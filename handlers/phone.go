package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/davidalvarez305/yd_cocktails/services"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/google/uuid"
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
			handleInboundCallRecordingCallback(w, r)
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

		twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<Response>
			<Dial record="true" recordingStatusCallback="%s" action="%s">%s</Dial>
		</Response>`, constants.RootDomain+constants.TwilioCallbackWebhook, forwardPhoneNumber)

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
	phoneCall.RecordingURL = dialStatus.RecordingURL
	phoneCall.Status = dialStatus.DialCallStatus

	if err := database.UpdatePhoneCall(phoneCall); err != nil {
		fmt.Printf("FAILED TO UPDATE PREVIOUS PHONE CALL: %+v\n", dialStatus)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	callMissed := phoneCall.CallDuration < 45

	isFirstCall, err := database.CheckIsFirstLeadContact(phoneCall.CallTo)
	if err != nil {
		fmt.Printf("ERROR CHECKING IF PHONE CALL IS FIRST LEAD CONTACT: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send 1st missed call text
	if callMissed && isFirstCall && !phoneCall.IsInbound {
		var textMessageTemplateNotification = fmt.Sprintf(
			`We missed you!
			%s
			`, `Hey! This is David with YD Cocktails.
			
			We're reaching out to you about your bartending service inquiry. We tried giving you a call but couldn't connect.
			
			Please give us a call back when you have a chance or let us know how we can help you.
			
			Todos hablamos español perfecto!`)

		sentMessage, err := services.SendTextMessage(phoneCall.CallTo, constants.CompanyPhoneNumber, textMessageTemplateNotification)

		if err != nil {
			fmt.Printf("ERROR SENDING MISSED CALL NOTIFICATION MSG: %+v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var externalID = helpers.SafeString(sentMessage.Sid)
		var messageStatus = helpers.SafeString(sentMessage.Status)

		msg := models.Message{
			ExternalID:  externalID,
			Text:        textMessageTemplateNotification,
			TextFrom:    constants.CompanyPhoneNumber,
			TextTo:      phoneCall.CallTo,
			IsInbound:   false,
			DateCreated: time.Now().Unix(),
			Status:      messageStatus,
			IsRead:      true,
		}

		err = database.SaveSMS(msg)
		if err != nil {
			fmt.Printf("ERROR SAVING MISSED CALL NOTIFICATION MSG: %+v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Transcribe phone call
	go func() {
		// Download the file from Twilio
		fileName := uuid.New().String() + ".mp3"
		localFilePath := constants.LOCAL_FILES_DIR + fileName
		audioFileURL := phoneCall.RecordingURL
		err := helpers.DownloadFileFromURL(audioFileURL, localFilePath)
		if err != nil {
			fmt.Printf("ERROR DOWNLOADING AUDIO FILE: %+v\n", err)
			http.Error(w, "Failed to download audio file", http.StatusInternalServerError)
			return
		}

		// Get file info before opening
		fileInfo, err := os.Stat(localFilePath)
		if err != nil {
			fmt.Printf("ERROR GETTING FILE INFO: %+v\n", err)
			http.Error(w, "Failed to get file info", http.StatusInternalServerError)
			return
		}

		// Open the file before uploading to S3
		file, err := os.Open(localFilePath)
		if err != nil {
			fmt.Printf("ERROR OPENING FILE: %+v\n", err)
			http.Error(w, "Failed to open audio file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Correct S3 file path (bucket name should not be part of the key)
		s3FilePath := "/uploads/audio/" + fileName

		// Upload audio file to S3
		err = services.UploadFileToS3(file, fileInfo.Size(), s3FilePath)
		if err != nil {
			fmt.Printf("ERROR UPLOADING AUDIO TO S3: %+v\n", err)
			http.Error(w, "Failed to upload audio file to S3", http.StatusInternalServerError)
			return
		}

		// Transcribe audio file
		transcriptionID, transcriptionText, err := services.TranscribeAudio(audioFileS3URL)
		if err != nil {
			fmt.Printf("ERROR TRANSCRIBING AUDIO: %+v\n", err)
			http.Error(w, "Failed to transcribe audio", http.StatusInternalServerError)
			return
		}

		// Save transcription in DB
		transcription := models.PhoneCallTranscription{
			PhoneCallID:           phoneCall.PhoneCallID,
			TranscriptionID:       transcriptionID,
			Transcription:         transcriptionText,
			TranscriptionAudioURL: audioFileS3URL,
		}

		err = database.SavePhoneCallTranscription(transcription)
		if err != nil {
			fmt.Printf("ERROR SAVING TRANSCRIPTION: %+v\n", err)
			http.Error(w, "Failed to save transcription", http.StatusInternalServerError)
			return
		}

		// Upload transcription text to S3
		transcriptionFileUUID := uuid.New().String() + ".txt"

		// Define the local file path where the transcription text will be saved
		localTranscriptionTextPath := constants.LOCAL_FILES_DIR + transcriptionFileUUID

		// Create the local file and write the transcription text to it
		transcriptionFile, err := os.Create(localTranscriptionTextPath)
		if err != nil {
			fmt.Printf("ERROR CREATING LOCAL FILE: %+v\n", err)
			http.Error(w, "Failed to create transcription file", http.StatusInternalServerError)
			return
		}
		defer transcriptionFile.Close()

		_, err = transcriptionFile.Write([]byte(transcriptionText))
		if err != nil {
			fmt.Printf("ERROR WRITING TO LOCAL FILE: %+v\n", err)
			http.Error(w, "Failed to write transcription to file", http.StatusInternalServerError)
			return
		}

		// Open the file for uploading to S3
		transcriptionFileToUpload, err := os.Open(localTranscriptionTextPath)
		if err != nil {
			fmt.Printf("ERROR OPENING FILE: %+v\n", err)
			http.Error(w, "Failed to open transcription file", http.StatusInternalServerError)
			return
		}
		defer transcriptionFileToUpload.Close()

		// Determine file size
		fileInfo, err = transcriptionFileToUpload.Stat()
		if err != nil {
			fmt.Printf("ERROR GETTING FILE INFO: %+v\n", err)
			http.Error(w, "Failed to get file info", http.StatusInternalServerError)
			return
		}
		fileSize := fileInfo.Size()

		// S3 file path
		s3TranscriptionFilePath := "uploads/transcriptions/" + transcriptionFileUUID

		// Upload the file to S3
		err = services.UploadFileToS3(transcriptionFileToUpload, fileSize, s3TranscriptionFilePath)
		if err != nil {
			fmt.Printf("ERROR UPLOADING TRANSCRIPTION TO S3: %+v\n", err)
			http.Error(w, "Failed to upload transcription to S3", http.StatusInternalServerError)
			return
		}

		// Save the transcription file URL in the database
		transcription.TranscriptionTextURL = s3TranscriptionFilePath // The S3 path returned from the upload
		err = database.UpdatePhoneCallTranscription(transcription)
		if err != nil {
			fmt.Printf("ERROR UPDATING TRANSCRIPTION URL: %+v\n", err)
			http.Error(w, "Failed to update transcription URL", http.StatusInternalServerError)
			return
		}

	}()

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

	twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Dial action="%s">%s</Dial>
	</Response>`, constants.RootDomain+constants.TwilioCallbackWebhook, "+1"+to)

	outboundCall, err := services.InitiateOutboundCall(from, twiML)
	if err != nil {
		fmt.Println("Error initiating phone call:", err)
		http.Error(w, "Failed to initiate phone call", http.StatusInternalServerError)
		return
	}

	var recordingURL string
	subResources := outboundCall.SubresourceUris

	if subResources != nil {
		if recordings, ok := (*subResources)["recordings"].(string); ok {
			recordingURL = recordings
		}
	}

	phoneCall := models.PhoneCall{
		ExternalID:   helpers.SafeString(outboundCall.Sid),
		CallDuration: 0,
		DateCreated:  time.Now().Unix(),
		CallFrom:     from,
		CallTo:       to,
		IsInbound:    false,
		RecordingURL: recordingURL,
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

func handleInboundCallRecordingCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	callSid := r.FormValue("CallSid")
	recordingURL := r.FormValue("RecordingUrl")

	if callSid == "" || recordingURL == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	if err := database.UpdateCallRecordingURL(callSid, recordingURL); err != nil {
		fmt.Printf("Failed to update call recording: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
