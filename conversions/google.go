package conversions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func SendGoogleConversion(payload types.GooglePayload) error {
	endpoint := "https://www.google-analytics.com/mp/collect"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling Google payload: %+v\n", err)
		return err
	}

	url := endpoint + "?measurement_id=" + constants.GoogleAnalyticsID + "&api_secret=" + constants.GoogleAnalyticsAPISecretKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending Google request: %+v\n", err)
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("GOOGLE REPORTING ERROR: %s\n", bodyString)

		return fmt.Errorf("google API returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
