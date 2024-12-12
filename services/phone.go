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
