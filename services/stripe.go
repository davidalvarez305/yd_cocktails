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

	if !params.ShouldSendInvoice {
		return *inv, nil
	}

	// Send the Invoice
	sentInvoice, err := invoice.SendInvoice(inv.ID, &stripe.InvoiceSendInvoiceParams{})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to send invoice: %v", err)
	}

	return *sentInvoice, nil
}

func UpdateStripeInvoice(leadQuoteInvoice types.LeadQuoteInvoice) error {
	stripe.Key = constants.StrikeAPIKey

	// Retrieve the existing invoice
	inv, err := invoice.Get(leadQuoteInvoice.StripeInvoiceID, &stripe.InvoiceParams{})
	if err != nil {
		return fmt.Errorf("failed to retrieve invoice: %v", err)
	}

	// Retrieve all invoice items associated with this invoice
	invoiceItems := invoiceitem.List(&stripe.InvoiceItemListParams{
		Invoice: stripe.String(inv.ID),
	})
	if err != nil {
		return fmt.Errorf("failed to list invoice items: %v", err)
	}

	// Delete the old invoice items that are associated with this invoice
	for invoiceItems.Next() {
		var item = invoiceItems.InvoiceItem()
		// Only delete invoice items that are associated with the current invoice and have not been paid
		if item.Invoice != nil && item.Invoice.ID == inv.ID && !item.Invoice.Paid {
			// Delete the invoice item
			_, err := invoiceitem.Del(item.ID, nil)
			if err != nil {
				return fmt.Errorf("failed to delete invoice item %v: %v", item.ID, err)
			}
		}
	}

	// Add the new invoice item with the updated amount
	_, err = invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(leadQuoteInvoice.StripeCustomerID),
		Amount:      stripe.Int64(int64(leadQuoteInvoice.Amount) * 100), // Amount in cents
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Bartending service."),
		Invoice:     stripe.String(inv.ID), // Attach to the existing invoice
	})

	if err != nil {
		return fmt.Errorf("failed to create invoice item: %v", err)
	}

	_, err = invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to finalize invoice: %v", err)
	}

	return nil
}
