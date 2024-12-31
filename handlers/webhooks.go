package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/conversions"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/webhooks/stripe-invoice":
			handleStripeInvoice(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleStripeInvoice(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Webhook error while parsing basic request. %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, constants.StrikeAPIKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handleStripePaymentSuccessful(invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling payment succeded: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "invoice.updated":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handleStripeInvoiceUpdated(invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling payment succeded: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func handleStripePaymentSuccessful(invoice stripe.Invoice) error {
	leadId, err := database.GetLeadByStripeCustomerID(invoice.Customer.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting lead details by stripe customer ID: %v\n", err)
		return err
	}

	lead, err := database.GetLeadDetails(fmt.Sprint(leadId))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting lead details: %v\n", err)
		return err
	}

	err = database.UpdateEstimateByWebhook(invoice.ID, string(invoice.Status), time.Now().Unix())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting updating estimate: %v\n", err)
		return err
	}

	fbEvent := types.FacebookEventData{
		EventName:      constants.InvoicePaidEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: lead.LandingPage,
		UserData: types.FacebookUserData{
			Email:           helpers.HashString(lead.Email),
			FirstName:       helpers.HashString(lead.FirstName),
			LastName:        helpers.HashString(lead.LastName),
			Phone:           helpers.HashString(lead.PhoneNumber),
			FBC:             lead.FacebookClickID,
			FBP:             lead.FacebookClientID,
			ExternalID:      helpers.HashString(lead.ExternalID),
			ClientIPAddress: lead.IP,
			ClientUserAgent: lead.UserAgent,
		},
		CustomData: types.FacebookCustomData{
			Currency: "USD",
			Value:    fmt.Sprint(float64(invoice.AmountPaid) / 100),
		},
		EventID: invoice.ID,
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	googlePayload := types.GooglePayload{
		ClientID: lead.GoogleClientID,
		UserId:   lead.ExternalID,
		Events: []types.GoogleEventLead{
			{
				Name: constants.InvoicePaidEventName,
				Params: types.GoogleEventParamsLead{
					GCLID:         lead.ClickID,
					TransactionID: invoice.ID,
					Value:         float64(invoice.AmountPaid) / 100,
					Currency:      constants.DefaultCurrency,
					CampaignID:    fmt.Sprint(lead.CampaignID),
					Campaign:      lead.CampaignName,
					Source:        lead.Source,
					Medium:        lead.Medium,
					Term:          lead.Keyword,
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
			Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},
			Address: []types.GoogleUserAddress{
				{
					Sha256FirstName: helpers.HashString(lead.FirstName),
					Sha256LastName:  helpers.HashString(lead.LastName),
				},
			},
		},
	}

	go conversions.SendGoogleConversion(googlePayload)
	go conversions.SendFacebookConversion(metaPayload)

	return nil
}

func handleStripeInvoiceUpdated(invoice stripe.Invoice) error {
	err := database.UpdateEstimatePriceByStripeInvoiceID(invoice.ID, float64(invoice.AmountDue)/100)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating estimate price by stripe invoice id: %v\n", err)
		return err
	}

	return nil
}
