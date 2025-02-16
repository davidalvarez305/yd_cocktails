package services

import (
	"fmt"

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

func InitiateOutboundCall(reversedNumberToDial, reversedCallerId, twiML string) (openapi.ApiV2010Call, error) {
	client := twilio.NewRestClient()

	var call openapi.ApiV2010Call
	var params openapi.CreateCallParams

	params.SetTo("+1" + reversedCallerId)       // caller id must be set "to" because it is the first call that happens (where the call will be coming from)
	params.SetFrom("+1" + reversedNumberToDial) // this number, the client's number, will be the number that is called to from the caller id's number
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
