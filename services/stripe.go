package services

import (
	"fmt"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/invoice"
	"github.com/stripe/stripe-go/v81/invoiceitem"
)

func CreateStripeInvoice(stripeCustomerId, email, fullName string, eventDate int64, quote float64) (stripe.Invoice, error) {
	stripe.Key = constants.StrikeAPIKey

	// Create a new customer if needed
	if len(stripeCustomerId) == 0 {
		custParams := &stripe.CustomerParams{
			Email: stripe.String(email),
			Name:  stripe.String(fullName),
		}
		cust, err := customer.New(custParams)
		if err != nil {
			return stripe.Invoice{}, fmt.Errorf("failed to create customer: %v", err)
		}
		stripeCustomerId = cust.ID
	}

	// Calculate payment due date
	t := time.Unix(eventDate, 0)
	paymentDueDate := t.Add(-time.Duration(constants.InvoicePaymentDueInHours) * time.Hour).Unix()

	// Create Invoice First
	invoiceParams := &stripe.InvoiceParams{
		Customer:         stripe.String(stripeCustomerId),
		CollectionMethod: stripe.String("send_invoice"),
		DueDate:          stripe.Int64(paymentDueDate),
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
		Customer:    stripe.String(stripeCustomerId),
		Amount:      stripe.Int64(int64(quote) * 100),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Bartending service."),
		Invoice:     stripe.String(inv.ID), // Attach to invoice
	})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice item: %v", err)
	}

	// Send the Invoice
	_, err = invoice.SendInvoice(inv.ID, &stripe.InvoiceSendInvoiceParams{})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to send invoice: %v", err)
	}

	return *inv, nil
}
