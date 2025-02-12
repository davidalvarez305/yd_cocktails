package services

import (
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
)

func UpdateInvoicesWorkflow(quoteId int, eventDate int64) error {
	// Get stripe invoices
	leadQuoteInvoices, err := database.GetLeadQuoteInvoices(quoteId)
	if err != nil {
		return err
	}

	for _, leadQuoteInvoice := range leadQuoteInvoices {

		// Calculate new due date
		dueDate := time.Now().Unix()
		if leadQuoteInvoice.InvoiceTypeID == constants.RemainingInvoiceTypeID {
			t := time.Unix(eventDate, 0)
			dueDate = t.Add(-time.Duration(constants.InvoicePaymentDueInHours) * time.Hour).Unix()
		}
		leadQuoteInvoice.DueDate = dueDate

		// Void old invoice and copy over to new invoice on stripe
		invoice, err := UpdateStripeInvoice(leadQuoteInvoice)
		if err != nil {
			return err
		}

		// Set old invoice status to void
		err = database.UpdateInvoiceStatus(leadQuoteInvoice.StripeInvoiceID, constants.VoidInvoiceStatusID)
		if err != nil {
			return err
		}

		// Create new invoice with status open
		err = database.CreateQuoteInvoice(invoice.ID, invoice.HostedInvoiceURL, quoteId, leadQuoteInvoice.InvoiceTypeID, invoice.DueDate)
		if err != nil {
			return err
		}
	}

	return nil
}
