package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/conversions"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/models"
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
		UserID:       forwardNumber.UserID,
		LeadID:       forwardNumber.LeadID,
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

	if phoneCall.CallDuration > 60 {
		fbEvent := types.FacebookEventData{
			EventName:      constants.LeadEventName,
			EventTime:      phoneCall.DateCreated,
			ActionSource:   "website",
			EventSourceURL: constants.RootDomain,
			UserData: types.FacebookUserData{
				Phone: helpers.HashString(phoneCall.CallFrom),
				State: helpers.HashString("Florida"),
			},
		}

		metaPayload := types.FacebookPayload{
			Data: []types.FacebookEventData{fbEvent},
		}

		err = conversions.SendFacebookConversion(metaPayload)

		if err != nil {
			fmt.Printf("Error sending Facebook conversion: %+v\n", err)
		}
	}

	if err := database.UpdatePhoneCall(phoneCall); err != nil {
		fmt.Printf("FAILED TO UPDATE PREVIOUS PHONE CALL: %+v\n", dialStatus)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}
