package services

import (
	"fmt"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/invoice"
	"github.com/stripe/stripe-go/v81/invoiceitem"
)

func CreateStripeInvoiceForNewCustomer(lead models.Lead, packagePrice float64) (string, string, error) {
	stripe.Key = constants.StrikeAPIKey

	custParams := &stripe.CustomerParams{
		Email: stripe.String(fmt.Sprintf("%s", lead.Email)),
		Name:  stripe.String(fmt.Sprintf("%s %s", lead.FirstName, lead.LastName)),
	}
	cust, err := customer.New(custParams)
	if err != nil {
		return "", "", fmt.Errorf("failed to create customer: %v", err)
	}

	depositAmount := packagePrice * constants.DepositPercentageAmount

	_, err = invoiceitem.New(&stripe.InvoiceItemParams{
		Customer:    stripe.String(cust.ID),
		Amount:      stripe.Int64(int64(depositAmount)),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("50% Deposit"),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create invoice item: %v", err)
	}

	invoiceParams := &stripe.InvoiceParams{
		Customer: stripe.String(cust.ID),
	}
	inv, err := invoice.New(invoiceParams)
	if err != nil {
		return "", "", fmt.Errorf("failed to create invoice: %v", err)
	}

	inv, err = invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to finalize invoice: %v", err)
	}

	return inv.ID, cust.ID, nil
}
