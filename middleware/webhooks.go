package middleware

import (
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/twilio/twilio-go/client"
)

func validateTwilioWebhook(r *http.Request) error {
	requestValidator := client.NewRequestValidator(constants.TwilioAuthToken)
	twilioSignature := r.Header.Get("X-Twilio-Signature")

	url := "https://" + r.Host + r.URL.Path

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}
	}

	params := make(map[string]string)
	for key, values := range r.Form {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	isValid := requestValidator.Validate(url, params, twilioSignature)

	if !isValid {
		return fmt.Errorf("invalid Twilio signature")
	}

	return nil
}
