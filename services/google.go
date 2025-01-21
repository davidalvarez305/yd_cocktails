package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func refreshAuthToken(config *oauth2.Config) (oauth2.Token, error) {
	var token oauth2.Token

	refreshToken := constants.GoogleRefreshToken
	client := &http.Client{}

	url := config.Endpoint.TokenURL
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
		return token, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("client_id", config.ClientID)
	q.Add("client_secret", config.ClientSecret)
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while getting auth token: ", err)
		return token, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("STATUS CODE: %+v\n", resp.Status)
		return token, errors.New("request failed")
	}

	json.NewDecoder(resp.Body).Decode(&token)

	return token, nil
}

func initializeGoogleClient(scope string) (*http.Client, error) {
	var client *http.Client

	b, err := os.ReadFile(constants.GoogleJSONPath)
	if err != nil {
		return client, err
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return client, err
	}

	token, err := refreshAuthToken(config)

	if err != nil {
		return nil, err
	}

	return config.Client(context.Background(), &token), nil
}

func SendGmail(recipients []string, subject, sender, body string) error {
	client, err := initializeGoogleClient(gmail.GmailSendScope)
	if err != nil {
		fmt.Printf("Unable to initialize Gmail client: %v", err)
		return err
	}

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	user := "me"

	var message gmail.Message

	emailContent := fmt.Sprintf("To: %s\r\nSubject: %s\r\nReply-To:%s\r\n%s", strings.Join(recipients, ", "), subject, sender, body)
	message.Raw = base64.URLEncoding.EncodeToString([]byte(emailContent))

	_, err = srv.Users.Messages.Send(user, &message).Do()
	if err != nil {
		fmt.Printf("Unable to send email: %v", err)
		return err
	}

	return nil
}

func GetDataFromSheets(spreadsheetId, spreadsheetRange string) (*sheets.ValueRange, error) {
	var values sheets.ValueRange

	client, err := initializeGoogleClient(sheets.SpreadsheetsScope)
	if err != nil {
		fmt.Printf("Unable to initialize Sheets client: %v", err)
		return &values, nil
	}

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Sheets client: %v", err)
		return &values, nil
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, spreadsheetRange).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve data from sheet: %v", err)
		return &values, nil
	}

	values = *resp

	return &values, nil
}

func OverwriteSheetsData(spreadsheetId, spreadsheetRange string, values *sheets.ValueRange) error {
	client, err := initializeGoogleClient(sheets.SpreadsheetsScope)
	if err != nil {
		fmt.Printf("Unable to initialize Sheets client: %v", err)
		return err
	}

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Sheets client: %v", err)
		return err
	}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, spreadsheetRange, values).Do()
	if err != nil {
		fmt.Printf("Unable to update data on sheet: %v", err)
		return err
	}

	return err
}
