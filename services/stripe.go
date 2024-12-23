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

func CreateStripeInvoiceForNewCustomer(email, firstName, lastName string, packagePrice float64) (stripe.Invoice, error) {
	stripe.Key = constants.StrikeAPIKey

	custParams := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(fmt.Sprintf("%s %s", firstName, lastName)),
	}
	cust, err := customer.New(custParams)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create customer: %v", err)
	}

	_, err = invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(cust.ID),
		Amount:      stripe.Int64(int64(packagePrice * 100)), // Convert to cents
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Full open bar service."),
	})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice item: %v", err)
	}

	invoiceParams := &stripe.InvoiceParams{
		Customer: stripe.String(cust.ID),
	}

	inv, err := invoice.New(invoiceParams)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice: %v", err)
	}

	inv, err = invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to finalize invoice: %v", err)
	}

	inv.DueDate = time.Now().Local().Unix()

	return *inv, nil
}
