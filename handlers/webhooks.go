package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database
	"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/stripe/invoice":
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
			http.Error(w, "Webhook processing error", http.StatusInternalServerError)
			return
		}

		// Update invoice status to paid
		err = database.SetInvoiceStatusToPaid(invoice.ID, invoice.DatePaid)
		if err != nil {
			log.Printf("Failed to get invoice by stripe invoice id: %v", err)
			http.Error(w, "Error updating invoice status to paid", http.StatusInternalServerError)
			return
		}

		// If invoice type = deposit, schedule event + report to google
		if inv.InvoiceTypeID = constants.DepositInvoiceTypeID {
			err = database.CreateEvent()
			if err != nil {
				log.Printf("Failed to get invoice by stripe invoice id: %v", err)
				http.Error(w, "Error creating event", http.StatusInternalServerError)
				return
			}

			if constants.Production {
				lead, err := database.GetConversionReporting(int(helpers.SafeInt(form.LeadID)))
				if err != nil {
					fmt.Printf("Error getting conversion: %+v\n", err)
					tmplCtx := types.DynamicPartialTemplate{
						TemplateName: "error",
						TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
						Data: map[string]any{
							"Message": "Internal error reporting conversions to Google.",
						},
					}
		
					w.WriteHeader(http.StatusBadRequest)
					helpers.ServeDynamicPartialTemplate(w, tmplCtx)
					return
				}
		
				if lead.FacebookClickID != "" {
					fbEvent := types.FacebookEventData{
						EventName:      constants.EventConversionEventName,
						EventTime:      helpers.SafeInt64(form.DatePaid),
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
						EventTime:    helpers.SafeInt64(form.DatePaid),
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

		var notifyList = []string{quote.PhoneNumber, ...constants.NotificationSubscribers}

		// Text notification event details
		for _, phoneNumber := range notifyList {
			var textMessageTemplateNotification = fmt.Sprintf(
				`EVENT BOOKED:
	
				Date: %s,
				Full Name: %s
			`, utils.FormatTimestampEST(quote.EventDate), lead.FullName, helpers.SafeString(form.Message))
	
			_, err := services.SendTextMessage(constants.DavidPhoneNumber, constants.CompanyPhoneNumber, textMessageTemplateNotification)
	
			if err != nil {
				fmt.Printf("ERROR SENDING EVENT BOOKED NOTIFICATION MSG: %+v\n", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
