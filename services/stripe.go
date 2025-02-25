package services

import (
	"fmt"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/invoice"
	"github.com/stripe/stripe-go/v81/invoiceitem"
)

func CreateStripeInvoice(params types.CreateInvoiceParams) (stripe.Invoice, error) {
	stripe.Key = constants.StrikeAPIKey

	// Create a new customer if needed
	if len(params.StripeCustomerID) == 0 {
		custParams := &stripe.CustomerParams{
			Email: stripe.String(params.Email),
			Name:  stripe.String(params.FullName),
		}
		cust, err := customer.New(custParams)
		if err != nil {
			return stripe.Invoice{}, fmt.Errorf("failed to create customer: %v", err)
		}
		params.StripeCustomerID = cust.ID
	}

	// Create Invoice First
	invoiceParams := &stripe.InvoiceParams{
		Customer:         stripe.String(params.StripeCustomerID),
		CollectionMethod: stripe.String("send_invoice"),
		DueDate:          stripe.Int64(params.DueDate),
		Description:      stripe.String("This deposit will lock down your event date."),
		Footer:           stripe.String("**TERMS & CONDITIONS**\nBartending services invoice.\n\nThank you for your business!"),
		Currency:         stripe.String("USD"),
	}

	inv, err := invoice.New(invoiceParams)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice: %v", err)
	}

	// Add Invoice Item (Attaching to Invoice)
	_, err = invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(params.StripeCustomerID),
		Amount:      stripe.Int64(int64(params.Quote) * 100),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Bartending service."),
		Invoice:     stripe.String(inv.ID), // Attach to invoice
	})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice item: %v", err)
	}

	finalizedInvoice, err := invoice.FinalizeInvoice(inv.ID, &stripe.InvoiceFinalizeInvoiceParams{})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to send invoice: %v", err)
	}

	return *finalizedInvoice, nil
}

func UpdateStripeInvoice(leadQuoteInvoice types.LeadQuoteInvoice) (stripe.Invoice, error) {
	stripe.Key = constants.StrikeAPIKey
	var updatedInvoice stripe.Invoice

	// Retrieve the existing invoice
	originalInvoice, err := invoice.Get(leadQuoteInvoice.StripeInvoiceID, nil)
	if err != nil {
		return updatedInvoice, fmt.Errorf("failed to retrieve invoice: %v", err)
	}

	// Void the existing invoice
	_, err = invoice.VoidInvoice(originalInvoice.ID, nil)
	if err != nil {
		return updatedInvoice, fmt.Errorf("failed to void invoice: %v", err)
	}

	// Create New Invoice
	invoiceParams := &stripe.InvoiceParams{
		Customer:         stripe.String(leadQuoteInvoice.StripeCustomerID),
		CollectionMethod: stripe.String(string(originalInvoice.CollectionMethod)),
		DueDate:          stripe.Int64(leadQuoteInvoice.DueDate),
		Description:      stripe.String(originalInvoice.Description),
		Footer:           stripe.String(originalInvoice.Footer),
		Currency:         stripe.String(constants.DefaultCurrency),
	}

	newInvoice, err := invoice.New(invoiceParams)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice: %v", err)
	}

	// Add Invoice Item (Attaching to Invoice)
	_, err = invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(leadQuoteInvoice.StripeCustomerID),
		Amount:      stripe.Int64(int64(leadQuoteInvoice.Amount*leadQuoteInvoice.InvoiceTypeMultiplier) * 100),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Bartending service."),
		Invoice:     stripe.String(newInvoice.ID), // Attach to invoice
	})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice item: %v", err)
	}

	// Finalize the new invoice
	finalizedInvoice, err := invoice.FinalizeInvoice(newInvoice.ID, nil)
	if err != nil {
		return updatedInvoice, fmt.Errorf("failed to finalize new invoice: %v", err)
	}

	updatedInvoice = *finalizedInvoice

	return updatedInvoice, nil
}

func GetStripeInvoice(stripeInvoiceId string) (stripe.Invoice, error) {
	stripe.Key = constants.StrikeAPIKey
	var stripeInvoice stripe.Invoice

	originalInvoice, err := invoice.Get(stripeInvoiceId, nil)
	if err != nil {
		return stripeInvoice, fmt.Errorf("failed to retrieve invoice: %v", err)
	}

	if originalInvoice == nil {
		return stripeInvoice, fmt.Errorf("invoice with stripeInvoiceId %s not found", stripeInvoiceId)
	}

	stripeInvoice = *originalInvoice

	return stripeInvoice, nil
}
