package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/models"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendTextMessage(to, from, body string) (openapi.ApiV2010Message, error) {
	client := twilio.NewRestClient()

	var params openapi.CreateMessageParams
	var text openapi.ApiV2010Message

	params.SetTo("+1" + to)
	params.SetFrom("+1" + from)
	params.SetBody(body)

	sentMessage, err := client.Api.CreateMessage(&params)

	if err != nil || sentMessage == nil {
		return text, err
	}

	text = *sentMessage

	return text, nil
}

func InitiateOutboundCall(from, twiML string) (openapi.ApiV2010Call, error) {
	client := twilio.NewRestClient()

	var call openapi.ApiV2010Call
	var params openapi.CreateCallParams

	params.SetTo("+1" + from)
	params.SetFrom("+1" + from)
	params.SetTwiml(twiML)
	params.SetMethod("POST")
	params.SetRecord(true)

	outboundCall, err := client.Api.CreateCall(&params)

	if err != nil || outboundCall == nil {
		fmt.Printf("ERROR INITIATING OUTBOUND CALL: %+v\n", err)
		return call, err
	}

	call = *outboundCall

	return call, nil
}

func DownloadFileFromTwilio(fileURL, localFilePath string) error {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(constants.TwilioAccountSID, constants.TwilioAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status: %s", resp.Status)
	}

	outFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Println("File downloaded successfully:", localFilePath)
	return nil
}

func MissedCallFollowUpText(phoneCall models.PhoneCall) error {
	var textMessageTemplateNotification = fmt.Sprintf(
		`We missed you!
		%s
		`, `Hey! This is David with YD Cocktails.
		
		We're reaching out to you about your bartending service inquiry. We tried giving you a call but couldn't connect.
		
		Please give us a call back when you have a chance or let us know how we can help you.
		
		Todos hablamos español perfecto!`)

	sentMessage, err := SendTextMessage(phoneCall.CallTo, constants.CompanyPhoneNumber, textMessageTemplateNotification)

	if err != nil {
		fmt.Printf("ERROR SENDING MISSED CALL NOTIFICATION MSG: %+v\n", err)
		return err
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
		return err
	}

	return nil
}
