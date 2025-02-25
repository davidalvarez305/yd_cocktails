package services

import (
	"fmt"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func UpdateInvoicesWorkflow(quoteId int, eventDate int64) error {
	// Get stripe invoices
	leadQuoteInvoices, err := database.GetLeadQuoteInvoices(quoteId)
	if err != nil {
		fmt.Printf("ERROR GETTING QUOTE INVOICES: %+v\n", err)
		return err
	}

	isDepositPaid, err := database.IsDepositPaid(quoteId)
	if err != nil {
		fmt.Printf("ERROR CHECKING IF DEPOSIT IS PAID: %+v\n", err)
		return err
	}

	var remainingInvoice types.LeadQuoteInvoice

	if isDepositPaid {
		remainingInvoice, err = database.GetRemainingInvoice(quoteId)
		if err != nil {
			fmt.Printf("ERROR GETTING REMAINING INVOICE: %+v\n", err)
			return err
		}

		depositStripeInvoiceId, err := database.GetDepositStripeInvoiceID(quoteId)
		if err != nil {
			fmt.Printf("ERROR GETTING DEPOSIT STRIPE INVOICE ID: %+v\n", err)
			return err
		}

		depositInvoice, err := GetStripeInvoice(depositStripeInvoiceId)
		if err != nil {
			fmt.Printf("ERROR GETTING STRIPE INVOICE: %+v\n", err)
			return err
		}

		// Subtract from the invoice amount what was already deducted from the deposit
		// This way the only thing that changes is the amount for which the invoice was updated by
		// The customer now only has to pay the difference between the new invoice amount and the deposit that was paid
		remainingInvoice.Amount = remainingInvoice.Amount - float64(depositInvoice.AmountPaid/100) // convert to cents before division
	}

	for _, leadQuoteInvoice := range leadQuoteInvoices {
		// If deposit is paid, the remaining invoice will be handled differently using pre-calculated values
		// In this way, we can continue this flow as usual, only changing the logic beforehand
		if isDepositPaid {
			leadQuoteInvoice = remainingInvoice
		}

		// Calculate new due date
		dueDate := time.Now().Add(24 * time.Hour).Unix()
		if leadQuoteInvoice.InvoiceTypeID == constants.RemainingInvoiceTypeID {
			t := time.Unix(eventDate, 0)
			dueDate = t.Add(-time.Duration(constants.InvoicePaymentDueInHours) * time.Hour).Unix()
		}
		leadQuoteInvoice.DueDate = dueDate

		// Void old invoice and copy over to new invoice on stripe
		invoice, err := UpdateStripeInvoice(leadQuoteInvoice)
		if err != nil {
			fmt.Printf("ERROR UPDATING STRIPE INVOICE: %+v\n", err)
			return err
		}

		// Set old invoice status to void
		err = database.UpdateInvoiceStatus(leadQuoteInvoice.StripeInvoiceID, constants.VoidInvoiceStatusID)
		if err != nil {
			fmt.Printf("ERROR UPDATING INVOICE STATUS: %+v\n", err)
			return err
		}

		// Create new invoice with status open
		err = database.CreateQuoteInvoice(invoice.ID, invoice.HostedInvoiceURL, quoteId, leadQuoteInvoice.InvoiceTypeID, invoice.DueDate)
		if err != nil {
			fmt.Printf("ERROR CREATING INVOICE: %+v\n", err)
			return err
		}
	}

	return nil
}
