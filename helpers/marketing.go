package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/davidalvarez305/yd_cocktails/utils"
)

func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashBytes)
	return hashedString
}

func FormatPhoneNumber(phoneNumber string) string {
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")
	if len(cleaned) == 10 {
		return fmt.Sprintf("(%s) %s - %s", cleaned[:3], cleaned[3:6], cleaned[6:])
	}
	return ""
}

func GetGoogleClientIDFromRequest(r *http.Request) (string, error) {
	gaCookie, err := r.Cookie("_ga")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _ga cookie found")
		}
		return "", err
	}

	parts := strings.Split(gaCookie.Value, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("unexpected _ga cookie format")
	}

	return parts[2] + "." + parts[3], nil
}

func GetFacebookClickIDFromRequest(r *http.Request) (string, error) {
	fbcCookie, err := r.Cookie("_fbc")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _fbc cookie found")
		}
		return "", err
	}

	return fbcCookie.Value, nil
}

func GetFacebookClientIDFromRequest(r *http.Request) (string, error) {
	fbpCookie, err := r.Cookie("_fbp")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _fbp cookie found")
		}
		return "", err
	}

	return fbpCookie.Value, nil
}

func MapInstantFormToQuoteForm(lead types.FacebookInstantFormLead) (types.QuoteForm, error) {
	var quoteForm types.QuoteForm

	fullName := lead.FullName
	phoneNumber := ExtractPhoneNumber(lead.PhoneNumber)
	message := lead.EventDescription
	email := lead.Email
	optInTextMessaging := true

	source := lead.Platform
	medium := constants.SocialMediaAdsMedium
	channel := constants.SocialMediaAdsChannel

	campaignID := ExtractMarketingID(lead.CampaignID)
	adCampaign := lead.AdName
	adSetID := ExtractMarketingID(lead.AdsetID)
	adSetName := lead.AdsetName
	adID := ExtractMarketingID(lead.AdID)
	adHeadline := ExtractMarketingID(lead.AdName)

	instantFormLeadId := ExtractMarketingID(lead.ID)
	formId := ExtractMarketingID(lead.FormID)
	formName := lead.FormName

	createdAt, err := utils.GetDateFromInstantForm(lead.CreatedTime)
	if err != nil {
		return quoteForm, fmt.Errorf("invalid date")
	}

	quoteForm = types.QuoteForm{
		FullName:           &fullName,
		PhoneNumber:        &phoneNumber,
		Message:            &message,
		Email:              &email,
		OptInTextMessaging: &optInTextMessaging,
		CreatedAt:          &createdAt,

		Source:        &source,
		Medium:        &medium,
		Channel:       &channel,
		LandingPage:   nil,
		Keyword:       nil,
		Referrer:      nil,
		ClickID:       nil,
		CampaignID:    campaignID,
		AdCampaign:    &adCampaign,
		AdSetID:       adSetID,
		AdSetName:     &adSetName,
		AdID:          adID,
		AdHeadline:    adHeadline,
		Language:      nil,
		Longitude:     nil,
		Latitude:      nil,
		UserAgent:     nil,
		ButtonClicked: nil,
		IP:            nil,

		CSRFToken:        nil,
		ExternalID:       nil,
		GoogleClientID:   nil,
		FacebookClickID:  nil,
		FacebookClientID: nil,
		CSRFSecret:       nil,

		InstantFormLeadID: instantFormLeadId,
		InstantFormID:     formId,
		InstantFormName:   &formName,
	}

	return quoteForm, nil
}

func ExtractPhoneNumber(input string) string {
	re := regexp.MustCompile(`\d+`)

	matches := re.FindAllString(input, -1)

	phoneNumber := ""
	for _, match := range matches {
		phoneNumber += match
	}

	return phoneNumber
}

func ExtractMarketingID(input string) *int64 {
	re := regexp.MustCompile(`\d+`)

	matches := re.FindAllString(input, -1)

	result := ""
	for _, match := range matches {
		result += match
	}

	marketingID, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return nil
	}

	return &marketingID
}
