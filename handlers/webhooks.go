package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/conversions"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/services"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/davidalvarez305/yd_cocktails/utils"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/webhooks/stripe/invoice":
			handleStripeInvoicePayment(w, r)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleStripeInvoicePayment(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request body read error", http.StatusBadRequest)
		return
	}

	// Retrieve the Stripe signature header
	sigHeader := r.Header.Get("Stripe-Signature")
	signingSecret := constants.StripeWebhookSecret

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(body, sigHeader, signingSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	// Handle the event
	switch event.Type {
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			log.Printf("Failed to parse invoice payment succeeded event: %v", err)
			http.Error(w, "Webhook processing error", http.StatusInternalServerError)
			return
		}

		inv, err := database.GetInvoiceByStripeInvoiceID(invoice.ID)
		if err != nil {
			log.Printf("Failed to get invoice by stripe invoice id: %v", err)
			http.Error(w, "Failed to find invoice by stripe invoice id.", http.StatusInternalServerError)
			return
		}

		// Update invoice status to paid
		datePaid := time.Now().Unix()
		dateEventCreated := time.Now().Unix()
		err = database.SetInvoiceStatusToPaid(invoice.ID, datePaid)
		if err != nil {
			log.Printf("Failed to update invoice status to paid: %v", err)
			http.Error(w, "Error updating invoice status to paid", http.StatusInternalServerError)
			return
		}

		quote, err := database.GetQuoteDetailsByStripeInvoiceID(inv.StripeInvoiceID)
		if err != nil {
			log.Printf("Failed to get quote details by stripe invoice id: %v", err)
			http.Error(w, "Error creating event", http.StatusInternalServerError)
			return
		}

		// If invoice type = deposit || full, schedule event + report to google
		if inv.InvoiceTypeID == constants.DepositInvoiceTypeID || inv.InvoiceTypeID == constants.FullInvoiceTypeID {
			eventForm := types.EventForm{
				LeadID:      &quote.LeadID,
				EventTypeID: &quote.EventTypeID,
				VenueTypeID: &quote.VenueTypeID,
				DateCreated: &dateEventCreated,
				DatePaid:    &datePaid,
				Amount:      &quote.Amount,
				Guests:      &quote.Guests,
			}
			err = database.CreateEvent(eventForm)
			if err != nil {
				log.Printf("Failed to create event after successful payment: %v", err)
				http.Error(w, "Error creating event", http.StatusInternalServerError)
				return
			}

			if constants.Production {
				lead, err := database.GetConversionReporting(int(helpers.SafeInt(eventForm.LeadID)))
				if err != nil {
					log.Printf("Error reporting getting conversion details: %v", err)
					http.Error(w, "Error reporting getting conversion details.", http.StatusInternalServerError)
					return
				}

				if lead.FacebookClickID != "" {
					fbEvent := types.FacebookEventData{
						EventName:      constants.EventConversionEventName,
						EventTime:      helpers.SafeInt64(eventForm.DatePaid),
						ActionSource:   "phone_call",
						EventSourceURL: lead.LandingPage,
						UserData: types.FacebookUserData{
							Email:           helpers.HashString(lead.Email),
							Phone:           helpers.HashString(lead.PhoneNumber),
							FBC:             lead.FacebookClickID,
							FBP:             lead.FacebookClientID,
							ExternalID:      helpers.HashString(lead.ExternalID),
							ClientIPAddress: lead.IP,
							ClientUserAgent: lead.UserAgent,
						},
						CustomData: types.FacebookCustomData{
							Currency: constants.DefaultCurrency,
							Value:    fmt.Sprint(lead.Revenue),
						},
						EventID: fmt.Sprint(lead.EventID),
					}

					metaPayload := types.FacebookPayload{
						Data: []types.FacebookEventData{fbEvent},
					}

					go conversions.SendFacebookConversion(metaPayload)
				} else {
					fbLeadAdEvent := types.FacebookEventData{
						EventName:    constants.EventConversionEventName,
						EventTime:    helpers.SafeInt64(eventForm.DatePaid),
						ActionSource: "phone_call",
						UserData: types.FacebookUserData{
							LeadID: lead.InstantFormLeadID,
						},
						CustomData: types.FacebookCustomData{
							Currency:        constants.DefaultCurrency,
							Value:           fmt.Sprint(lead.Revenue),
							EventSource:     constants.EventSourceCRM,
							LeadEventSource: constants.CompanyName,
						},
						EventID: fmt.Sprint(lead.EventID),
					}

					metaLeadAdPayload := types.FacebookPayload{
						Data: []types.FacebookEventData{fbLeadAdEvent},
					}

					go conversions.SendFacebookConversion(metaLeadAdPayload)
				}

				googlePayload := types.GooglePayload{
					ClientID: lead.GoogleClientID,
					UserId:   lead.ExternalID,
					Events: []types.GoogleEventLead{
						{
							Name: constants.EventConversionEventName,
							Params: types.GoogleEventParamsLead{
								GCLID:         lead.ClickID,
								TransactionID: fmt.Sprint(lead.EventID),
								Value:         lead.Revenue,
								Currency:      constants.DefaultCurrency,
								CampaignID:    fmt.Sprint(lead.CampaignID),
								Campaign:      lead.CampaignName,
							},
						},
					},
					UserData: types.GoogleUserData{
						Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
						Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},
					},
				}

				go conversions.SendGoogleConversion(googlePayload)
			}

			var notifyList = append([]string{quote.PhoneNumber}, constants.NotificationSubscribers...)

			// Text notification event details
			for _, phoneNumber := range notifyList {
				var textMessageTemplateNotification = fmt.Sprintf(
					`EVENT BOOKED:
	
				Date: %s,
				Full Name: %s
			`, utils.FormatTimestamp(quote.EventDate), quote.FullName)

				_, err := services.SendTextMessage(phoneNumber, constants.CompanyPhoneNumber, textMessageTemplateNotification)

				if err != nil {
					fmt.Printf("ERROR SENDING EVENT BOOKED NOTIFICATION MSG: %+v\n", err)
				}
			}
		}

		// Void all invoices if full or remaining has been paid
		if inv.InvoiceTypeID == constants.FullInvoiceTypeID || inv.InvoiceTypeID == constants.RemainingInvoiceTypeID {
			err = database.SetOpenInvoicesToVoid(quote.QuoteID)
			if err != nil {
				fmt.Printf("ERROR SETTING INVOICES TO VOID: %+v\n", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
