package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/davidalvarez305/yd_cocktails/services"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func PhoneServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/call/inbound":
			handleInboundCall(w, r)
		case "/call/inbound/end":
			handleInboundCallEnd(w, r)
		case "/sms/inbound":
			handleInboundSMS(w, r)
		case "/sms/outbound":
			handleOutboundSMS(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleInboundCall(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	incomingPhoneCall := types.TwilioIncomingCallBody{
		CallSid:       r.FormValue("CallSid"),
		AccountSid:    r.FormValue("AccountSid"),
		From:          r.FormValue("From"),
		To:            r.FormValue("To"),
		CallStatus:    r.FormValue("CallStatus"),
		ApiVersion:    r.FormValue("ApiVersion"),
		Direction:     r.FormValue("Direction"),
		ForwardedFrom: r.FormValue("ForwardedFrom"),
		CallerName:    r.FormValue("CallerName"),
		FromCity:      r.FormValue("FromCity"),
		FromState:     r.FormValue("FromState"),
		FromZip:       r.FormValue("FromZip"),
		FromCountry:   r.FormValue("FromCountry"),
		ToCity:        r.FormValue("ToCity"),
		ToState:       r.FormValue("ToState"),
		ToZip:         r.FormValue("ToZip"),
		ToCountry:     r.FormValue("ToCountry"),
		Caller:        r.FormValue("Caller"),
		Digits:        r.FormValue("Digits"),
		SpeechResult:  r.FormValue("SpeechResult"),
	}

	// Convert Confidence to float64
	if confidenceStr := r.FormValue("Confidence"); confidenceStr != "" {
		if confidence, err := strconv.ParseFloat(confidenceStr, 64); err == nil {
			incomingPhoneCall.Confidence = confidence
		}
	}

	forwardNumber, err := database.GetForwardPhoneNumber(helpers.RemoveCountryCode(incomingPhoneCall.To), helpers.RemoveCountryCode(incomingPhoneCall.From))
	if err != nil {
		fmt.Printf("Failed to get matching phone number: %+v\n", err)
		http.Error(w, "Failed to get matching phone number.", http.StatusInternalServerError)
		return
	}

	twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Dial action="%s">%s</Dial>
	</Response>`, constants.RootDomain+constants.TwilioCallbackWebhook, forwardNumber.ForwardPhoneNumber)

	phoneCall := models.PhoneCall{
		ExternalID:   incomingPhoneCall.CallSid,
		CallDuration: 0,
		DateCreated:  time.Now().Unix(),
		CallFrom:     incomingPhoneCall.From,
		CallTo:       incomingPhoneCall.To,
		IsInbound:    true,
		RecordingURL: "",
		Status:       incomingPhoneCall.CallStatus,
	}

	if err := database.SavePhoneCall(phoneCall); err != nil {
		fmt.Printf("Failed to save phone call: %+v\n", err)
		http.Error(w, "Failed to save phone call.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(twiML))
}

func handleInboundCallEnd(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	var dialStatus types.IncomingPhoneCallDialStatus

	dialStatus.DialCallStatus = r.FormValue("DialCallStatus")
	dialStatus.CallSid = r.FormValue("CallSid")

	if durationStr := r.FormValue("DialCallDuration"); durationStr != "" {
		if duration, err := strconv.Atoi(durationStr); err == nil {
			dialStatus.DialCallDuration = duration
		}
	}

	phoneCall, err := database.GetPhoneCallBySID(dialStatus.CallSid)
	if err != nil {
		fmt.Printf("FAILED TO GET PREVIOUS PHONE CALL: %+v\n", dialStatus)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	phoneCall.CallDuration = dialStatus.DialCallDuration
	phoneCall.RecordingURL = dialStatus.RecordingURL
	phoneCall.Status = dialStatus.DialCallStatus

	if err := database.UpdatePhoneCall(phoneCall); err != nil {
		fmt.Printf("FAILED TO UPDATE PREVIOUS PHONE CALL: %+v\n", dialStatus)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}

func handleInboundSMS(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	var twilioMessage types.TwilioMessage

	twilioMessage.MessageSid = r.FormValue("MessageSid")
	twilioMessage.AccountSid = r.FormValue("AccountSid")
	twilioMessage.MessagingServiceSid = r.FormValue("MessagingServiceSid")
	twilioMessage.From = r.FormValue("From")
	twilioMessage.To = r.FormValue("To")
	twilioMessage.Body = r.FormValue("Body")
	twilioMessage.NumMedia = r.FormValue("NumMedia")
	twilioMessage.NumSegments = r.FormValue("NumSegments")
	twilioMessage.SmsStatus = r.FormValue("SmsStatus")
	twilioMessage.ApiVersion = r.FormValue("ApiVersion")

	message := models.Message{
		ExternalID:  twilioMessage.MessageSid,
		Text:        twilioMessage.Body,
		TextFrom:    helpers.RemoveCountryCode(twilioMessage.From),
		TextTo:      helpers.RemoveCountryCode(twilioMessage.To),
		IsInbound:   true,
		DateCreated: time.Now().Unix(),
		Status:      twilioMessage.SmsStatus,
	}

	if err := database.SaveSMS(message); err != nil {
		log.Printf("Error saving SMS to database: %s", err)
		http.Error(w, "Failed to save message.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleOutboundSMS(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.OutboundMessageForm
	err = decoder.Decode(&form, r.PostForm)

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error decoding form data.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	messageResponse, err := services.SendTextMessage(form.To, form.From, form.Body)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to send text message.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var externalID = helpers.SafeString(messageResponse.Sid)
	var messageStatus = helpers.SafeString(messageResponse.Status)

	message := models.Message{
		ExternalID:  externalID,
		Text:        form.Body,
		TextFrom:    form.From,
		TextTo:      form.To,
		IsInbound:   false,
		DateCreated: time.Now().Unix(),
		Status:      messageStatus,
	}

	err = database.SaveSMS(message)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to save message.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	messages, err := database.GetMessagesByLeadID(form.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get new messages.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_messages.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_messages.html",
		Data: map[string]any{
			"Messages": messages,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
