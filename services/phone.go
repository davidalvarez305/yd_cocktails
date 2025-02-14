package services

import (
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendTextMessage(to, from, body string) (*openapi.ApiV2010Message, error) {
	client := twilio.NewRestClient()

	var params openapi.CreateMessageParams

	params.SetTo("+1" + to)
	params.SetFrom("+1" + from)
	params.SetBody(body)

	return client.Api.CreateMessage(&params)
}

func InitiateOutboundCall(to, from, twiML string) (openapi.ApiV2010Call, error) {
	client := twilio.NewRestClient()

	var call openapi.ApiV2010Call
	var params openapi.CreateCallParams

	params.SetTo("+1" + to)
	params.SetFrom("+1" + from)
	params.SetTwiml(twiML)
	params.SetMethod("POST")
	params.SetRecord(true)

	outboundCall, err := client.Api.CreateCall(&params)

	if err != nil || outboundCall == nil {
		return call, err
	}

	call = *outboundCall

	return call, nil
}
