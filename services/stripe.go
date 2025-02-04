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

	_, err := invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(stripeCustomerId),
		Amount:      stripe.Int64(int64(quote * 100)),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Bartending service."),
	})
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice item: %v", err)
	}

	t := time.Unix(eventDate, 0)
	paymentDueDate := t.Add(-1 * time.Duration(constants.InvoicePaymentDueInHours) * time.Hour).Unix()

	memo := `***TERMS & CONDITIONS***`

	invoiceParams := &stripe.InvoiceParams{
		Customer:         stripe.String(stripeCustomerId),
		DueDate:          stripe.Int64(paymentDueDate),
		CollectionMethod: stripe.String("send_invoice"),
		Description:      stripe.String(memo),
		Currency:         stripe.String("USD"),
	}

	inv, err := invoice.New(invoiceParams)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to create invoice: %v", err)
	}

	inv, err = invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		return stripe.Invoice{}, fmt.Errorf("failed to finalize invoice: %v", err)
	}

	return *inv, nil
}
