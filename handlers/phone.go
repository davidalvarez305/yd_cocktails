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
		case "/call/outbound":
			handleOutboundCall(w, r)
		case "/call/inbound/end":
			handleInboundCallEnd(w, r)
		case "/call/inbound/recording-callback":
			handleCallRecordingCallback(w, r)
		case "/sms/inbound":
			handleInboundSMS(w, r)
		case "/sms/outbound":
			handleOutboundSMS(w, r)
		case "/call/inbound/amd":
			handleAmdStatusCallback(w, r)
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
		CallSid:    r.FormValue("CallSid"),
		From:       r.FormValue("From"),
		To:         r.FormValue("To"),
		CallStatus: r.FormValue("CallStatus"),
	}

	if incomingPhoneCall.To != incomingPhoneCall.From {
		forwardPhoneNumber, err := database.GetForwardPhoneNumber(helpers.RemoveCountryCode(incomingPhoneCall.To), helpers.RemoveCountryCode(incomingPhoneCall.From))
		if err != nil {
			fmt.Printf("Failed to get matching phone number: %+v\n", err)
			http.Error(w, "Failed to get matching phone number.", http.StatusInternalServerError)
			return
		}

		recordingCallbackURL := fmt.Sprintf("%s%s", constants.RootDomain, constants.TwilioRecordingCallbackWebhook)

		twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<Response>
			<Dial record="true" recordingStatusCallback="%s" recordingStatusCallbackEvent="completed" action="%s">%s</Dial>
		</Response>`, recordingCallbackURL, constants.RootDomain+constants.TwilioCallbackWebhook, forwardPhoneNumber)

		phoneCall := models.PhoneCall{
			ExternalID:   incomingPhoneCall.CallSid,
			CallDuration: 0,
			DateCreated:  time.Now().Unix(),
			CallFrom:     helpers.RemoveCountryCode(incomingPhoneCall.From),
			CallTo:       helpers.RemoveCountryCode(incomingPhoneCall.To),
			IsInbound:    true,
			RecordingURL: "",
			Status:       incomingPhoneCall.CallStatus,
		}

		if err := database.SavePhoneCall(phoneCall); err != nil {
			fmt.Printf("Failed to save phone call: %+v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(twiML))
	}
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
	phoneCall.Status = dialStatus.DialCallStatus

	if err := database.UpdatePhoneCall(phoneCall); err != nil {
		fmt.Printf("FAILED TO UPDATE PREVIOUS PHONE CALL: %+v\n", dialStatus)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}

func handleOutboundCall(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	from := r.FormValue("from")
	to := r.FormValue("to")

	if to == "" || from == "" {
		http.Error(w, "Missing required parameters (From, To)", http.StatusBadRequest)
		return
	}

	recordingCallbackURL := fmt.Sprintf("%s%s", constants.RootDomain, constants.TwilioRecordingCallbackWebhook)

	twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Dial record="true"
			  recordingStatusCallback="%s"
			  recordingStatusCallbackEvent="completed"
			  action="%s">%s</Dial>
	</Response>`,
		recordingCallbackURL,
		constants.RootDomain+constants.TwilioCallbackWebhook,
		"+1"+to)

	outboundCall, err := services.InitiateOutboundCall(from, twiML)
	if err != nil {
		fmt.Println("Error initiating phone call:", err)
		http.Error(w, "Failed to initiate phone call", http.StatusInternalServerError)
		return
	}

	phoneCall := models.PhoneCall{
		ExternalID:   helpers.SafeString(outboundCall.Sid),
		CallDuration: 0,
		DateCreated:  time.Now().Unix(),
		CallFrom:     from,
		CallTo:       to,
		IsInbound:    false,
		RecordingURL: "",
		Status:       helpers.SafeString(outboundCall.Status),
	}

	if err := database.SavePhoneCall(phoneCall); err != nil {
		fmt.Printf("Failed to save phone call: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleInboundSMS(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	var twilioMessage types.TwilioMessage

	twilioMessage.MessageSid = r.FormValue("MessageSid")
	twilioMessage.From = r.FormValue("From")
	twilioMessage.To = r.FormValue("To")
	twilioMessage.Body = r.FormValue("Body")
	twilioMessage.NumMedia = r.FormValue("NumMedia")
	twilioMessage.NumSegments = r.FormValue("NumSegments")
	twilioMessage.SmsStatus = r.FormValue("SmsStatus")

	message := models.Message{
		ExternalID:  twilioMessage.MessageSid,
		Text:        twilioMessage.Body,
		TextFrom:    helpers.RemoveCountryCode(twilioMessage.From),
		TextTo:      helpers.RemoveCountryCode(twilioMessage.To),
		IsInbound:   true,
		DateCreated: time.Now().Unix(),
		Status:      twilioMessage.SmsStatus,
		IsRead:      false,
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
		IsRead:      true,
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

	leadMessages, err := database.GetMessagesByLeadID(form.LeadID)
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_messages.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_messages.html",
		Data: map[string]any{
			"LeadMessages": leadMessages,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func handleCallRecordingCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	callSID := r.FormValue("CallSid")
	recordingSID := r.FormValue("RecordingSid")

	if callSID == "" || recordingSID == "" {
		http.Error(w, "Missing CallSid or RecordingSid", http.StatusBadRequest)
		return
	}

	recordingURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Recordings/%s.mp3?RequestedChannels=2", constants.TwilioAccountSID, recordingSID)

	if err := database.SetRecordingURLToPhoneCall(callSID, recordingURL); err != nil {
		fmt.Printf("FAILED TO UPDATE PHONE CALL WITH RECORDING URL: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}

func handleAmdStatusCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for AMD status callback")

	// Parse request form data
	if err := r.ParseForm(); err != nil {
		fmt.Println("ERROR: Failed to parse form data:", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract Twilio parameters
	callSid := r.FormValue("CallSid")
	amdStatus := r.FormValue("AnsweredBy") // Possible values: 'human', 'machine_start', etc.

	fmt.Printf("Parsed AMD Callback Data - CallSid: %s, AnsweredBy: %s\n", callSid, amdStatus)

	// Validate required parameters
	if callSid == "" || amdStatus == "" {
		fmt.Println("ERROR: Missing required parameters (CallSid or AnsweredBy)")
		http.Error(w, "Missing required parameters (CallSid, AnsweredBy)", http.StatusBadRequest)
		return
	}

	// If a human answered, no further action is needed
	if amdStatus == "human" {
		fmt.Println("Call was answered by a human, no further action required.")
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("CALL NOT HUMAN.")

	// Fetch the existing call record from the database
	fmt.Println("Fetching phone call record from database...")
	phoneCall, err := database.GetPhoneCallBySID(callSid)
	if err != nil {
		fmt.Printf("ERROR: Failed to get phone call record (CallSid: %s) - %+v\n", callSid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Successfully fetched phone call record: %+v\n", phoneCall)

	// Get user details based on the caller's phone number
	fmt.Println("Fetching user details from database...")
	user, err := database.GetUserByPhoneNumber(phoneCall.CallFrom)
	if err != nil {
		fmt.Printf("ERROR: Failed to get user from phone number (%s) - %+v\n", phoneCall.CallFrom, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Successfully fetched user: %+v\n", user)

	// Mark call as "missed"
	phoneCall.Status = "missed"
	fmt.Println("Updating phone call status to 'missed'...")
	if err := database.UpdatePhoneCall(phoneCall); err != nil {
		fmt.Printf("ERROR: Failed to update phone call record - %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully updated phone call status to 'missed'.")

	// Check if this is the first time this lead has been contacted
	fmt.Println("Checking if this is the first lead contact...")
	isFirstCall, err := database.CheckIsFirstLeadContact(phoneCall.CallTo)
	if err != nil {
		fmt.Printf("ERROR: Failed to check if first lead contact - %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Is first lead contact: %v\n", isFirstCall)

	// If not the first call, we can stop here
	if !isFirstCall {
		fmt.Println("Not the first lead contact. No further action required.")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Send follow-up text for a missed call
	fmt.Println("Sending missed call follow-up text...")
	err = services.MissedCallFollowUpText(phoneCall, user)
	if err != nil {
		fmt.Printf("ERROR: Failed to send missed call follow-up text - %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully sent missed call follow-up text.")

	// Fetch the lead ID from the phone number
	fmt.Println("Fetching lead ID from phone number...")
	leadId, err := database.GetLeadIDFromPhoneNumber(phoneCall.CallTo)
	if err != nil {
		fmt.Printf("ERROR: Failed to get lead ID from phone number - %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Successfully fetched lead ID: %d\n", leadId)

	// Schedule the next follow-up action
	fmt.Println("Scheduling first follow-up action...")
	firstFollowUpActionID := constants.FirstFollowUpActionID
	twentyFourHours := time.Now().Add(24 * time.Hour).Unix()

	err = database.CreateLeadNextAction(types.LeadNextActionForm{
		NextActionID:   &firstFollowUpActionID,
		LeadID:         &leadId,
		NextActionDate: &twentyFourHours,
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to save next lead action - %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully scheduled first follow-up action.")

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Println("AMD callback processing completed successfully.")
}
