package types

import (
	"time"

	"github.com/davidalvarez305/yd_cocktails/models"
)

type QuoteForm struct {
	FullName           *string `json:"full_name" form:"full_name" schema:"full_name"`
	PhoneNumber        *string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Message            *string `json:"message" form:"message" schema:"message"`
	Email              *string `json:"email" form:"email" schema:"email"`
	OptInTextMessaging *bool   `json:"opt_in_text_messaging" form:"opt_in_text_messaging" schema:"opt_in_text_messaging"`
	CreatedAt          *int64  `json:"created_at" form:"created_at" schema:"created_at"`

	Source        *string `json:"source" form:"source" schema:"source"`
	Medium        *string `json:"medium" form:"medium" schema:"medium"`
	Channel       *string `json:"channel" form:"channel" schema:"channel"`
	LandingPage   *string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	Keyword       *string `json:"keyword" form:"keyword" schema:"keyword"`
	Referrer      *string `json:"referrer" form:"referrer" schema:"referrer"`
	ClickID       *string `json:"click_id" form:"click_id" schema:"click_id"`
	CampaignID    *int64  `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign    *string `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID     *int64  `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName   *string `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID       *int64  `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName     *string `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID          *int64  `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline    *int64  `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
	Language      *string `json:"language" form:"language" schema:"language"`
	Longitude     *string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude      *string `json:"latitude" form:"latitude" schema:"latitude"`
	UserAgent     *string `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ButtonClicked *string `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	IP            *string `json:"ip" form:"ip" schema:"ip"`

	CSRFToken        *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	ExternalID       *string `json:"external_id" form:"external_id" schema:"external_id"`
	GoogleClientID   *string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID  *string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID *string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	CSRFSecret       *string `json:"csrf_secret" form:"csrf_secret"`

	InstantFormLeadID *int64  `json:"instant_form_lead_id" form:"instant_form_lead_id" schema:"instant_form_lead_id"`
	InstantFormID     *int64  `json:"instant_form_id" form:"instant_form_id" schema:"instant_form_id"`
	InstantFormName   *string `json:"instant_form_name" form:"instant_form_name" schema:"instant_form_name"`
	ReferralLeadID    *int    `json:"referral_lead_id" form:"referral_lead_id" schema:"referral_lead_id"`
}

type ContactForm struct {
	CSRFToken string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	FullName  string `json:"full_name" form:"full_name" schema:"full_name"`
	LastName  string `json:"last_name" form:"last_name" schema:"last_name"`
	Email     string `json:"email" form:"email" schema:"email"`
	Message   string `json:"message" form:"message" schema:"message"`
}

type OutboundMessageForm struct {
	To        string `json:"to" form:"to" schema:"to"`
	Body      string `json:"body" form:"body" schema:"body"`
	From      string `json:"from" form:"from" schema:"from"`
	UserID    int    `json:"user_id" form:"user_id" schema:"user_id"`
	LeadID    int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	CSRFToken string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
}

type LeadDetails struct {
	LeadID           int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FullName         string `json:"full_name" form:"full_name" schema:"full_name"`
	Email            string `json:"email" form:"email" schema:"email"`
	PhoneNumber      string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	StripeCustomerID string `json:"stripe_customer_id" form:"stripe_customer_id" schema:"stripe_customer_id"`

	CampaignName     string `json:"campaign_name" form:"campaign_name" schema:"campaign_name"`
	CampaignID       int64  `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	Medium           string `json:"medium" form:"medium" schema:"medium"`
	Source           string `json:"source" form:"source" schema:"source"`
	Referrer         string `json:"referrer" form:"referrer" schema:"referrer"`
	LandingPage      string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	IP               string `json:"ip" form:"ip" schema:"ip"`
	Keyword          string `json:"keyword" form:"keyword" schema:"keyword"`
	Channel          string `json:"channel" form:"channel" schema:"channel"`
	Language         string `json:"language" form:"language" schema:"language"`
	Message          string `json:"message" form:"message" schema:"message"`
	FacebookClickID  string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	UserAgent        string `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ExternalID       string `json:"external_id" form:"external_id" schema:"external_id"`
	ClickID          string `json:"click_id" form:"click_id" schema:"click_id"`
	GoogleClientID   string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	ButtonClicked    string `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	ReferralLeadID   int    `json:"referral_lead_id" form:"referral_lead_id" schema:"referral_lead_id"`

	InstantFormLeadID int64  `json:"instant_form_lead_id" form:"instant_form_lead_id" schema:"instant_form_lead_id"`
	InstantFormID     int64  `json:"instant_form_id" form:"instant_form_id" schema:"instant_form_id"`
	InstantFormName   string `json:"instant_form_name" form:"instant_form_name" schema:"instant_form_name"`

	NextActionID   int `json:"next_action_id" form:"next_action_id" schema:"next_action_id"`
	LeadInterestID int `json:"lead_interest_id" form:"lead_interest_id" schema:"lead_interest_id"`
	LeadStatusID   int `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
}

type LeadList struct {
	LeadID       int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FullName     string `json:"full_name" form:"full_name" schema:"full_name"`
	PhoneNumber  string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	CreatedAt    string `json:"created_at" form:"created_at" schema:"created_at"`
	Language     string `json:"language" form:"language" schema:"language"`
	NextAction   string `json:"next_action" form:"next_action" schema:"next_action"`
	LeadInterest string `json:"lead_interest" form:"lead_interest" schema:"lead_interest"`
	LeadStatus   string `json:"lead_status" form:"lead_status" schema:"lead_status"`
	TotalRows    int    `json:"total_rows" form:"total_rows" schema:"total_rows"`
}

type Referral struct {
	LeadID   int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FullName string `json:"full_name" form:"full_name" schema:"full_name"`
}

type GetLeadsParams struct {
	PageNum        *string `json:"page_num" form:"page_num" schema:"page_num"`
	Search         *string `json:"search" form:"search" schema:"search"`
	LeadInterestID *int    `json:"lead_interest_id" form:"lead_interest_id" schema:"lead_interest_id"`
	LeadStatusID   *int    `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	NextActionID   *int    `json:"next_action_id" form:"next_action_id" schema:"next_action_id"`
}

type DynamicPartialTemplate struct {
	TemplateName string
	TemplatePath string
	Data         map[string]any
}

type TwilioMessage struct {
	MessageSid          string    `json:"MessageSid"`
	AccountSid          string    `json:"AccountSid"`
	MessagingServiceSid string    `json:"MessagingServiceSid"`
	From                string    `json:"From"`
	To                  string    `json:"To"`
	Body                string    `json:"Body"`
	NumMedia            string    `json:"NumMedia"`
	NumSegments         string    `json:"NumSegments"`
	SmsStatus           string    `json:"SmsStatus"`
	ApiVersion          string    `json:"ApiVersion"`
	DateCreated         time.Time `json:"DateCreated"`
}

type UpdateLeadForm struct {
	Method           *string `json:"_method" form:"_method" schema:"_method"`
	CSRFToken        *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	LeadID           *string `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FullName         *string `json:"full_name" form:"full_name" schema:"full_name"`
	PhoneNumber      *string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email            *string `json:"email" form:"email" schema:"email"`
	StripeCustomerID *string `json:"stripe_customer_id" form:"stripe_customer_id" schema:"stripe_customer_id"`

	LeadInterestID *int `json:"lead_interest_id" form:"lead_interest_id" schema:"lead_interest_id"`
	LeadStatusID   *int `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	NextActionID   *int `json:"next_action_id" form:"next_action_id" schema:"next_action_id"`
}

type UpdateLeadMarketingForm struct {
	Method         *string `json:"_method" form:"_method" schema:"_method"`
	CSRFToken      *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	LeadID         *string `json:"lead_id" form:"lead_id" schema:"lead_id"`
	CampaignName   *string `json:"campaign_name" form:"campaign_name" schema:"campaign_name"`
	Medium         *string `json:"medium" form:"medium" schema:"medium"`
	Source         *string `json:"source" form:"source" schema:"source"`
	Referrer       *string `json:"referrer" form:"referrer" schema:"referrer"`
	LandingPage    *string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	IP             *string `json:"ip" form:"ip" schema:"ip"`
	Keyword        *string `json:"keyword" form:"keyword" schema:"keyword"`
	Channel        *string `json:"channel" form:"channel" schema:"channel"`
	Language       *string `json:"language" form:"language" schema:"language"`
	ButtonClicked  *string `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	ReferralLeadID *int    `json:"referral_lead_id" form:"referral_lead_id" schema:"referral_lead_id"`
}

type TwilioSMSResponse struct {
	Sid                 string            `json:"sid"`
	DateCreated         string            `json:"date_created"`
	DateUpdated         string            `json:"date_updated"`
	DateSent            string            `json:"date_sent"`
	AccountSid          string            `json:"account_sid"`
	To                  string            `json:"to"`
	From                string            `json:"from"`
	MessagingServiceSid string            `json:"messaging_service_sid"`
	Body                string            `json:"body"`
	Status              string            `json:"status"`
	NumSegments         string            `json:"num_segments"`
	NumMedia            string            `json:"num_media"`
	Direction           string            `json:"direction"`
	ApiVersion          string            `json:"api_version"`
	Price               string            `json:"price"`
	PriceUnit           string            `json:"price_unit"`
	ErrorCode           string            `json:"error_code"`
	ErrorMessage        string            `json:"error_message"`
	Uri                 string            `json:"uri"`
	SubresourceUris     map[string]string `json:"subresource_uris"`
}

type TwilioIncomingCallBody struct {
	CallSid       string  `json:"CallSid" form:"CallSid" schema:"CallSid"`
	AccountSid    string  `json:"AccountSid" form:"AccountSid" schema:"AccountSid"`
	From          string  `json:"From" form:"From" schema:"From"`
	To            string  `json:"To" form:"To" schema:"To"`
	CallStatus    string  `json:"CallStatus" form:"CallStatus" schema:"CallStatus"`
	ApiVersion    string  `json:"ApiVersion" form:"ApiVersion" schema:"ApiVersion"`
	Direction     string  `json:"Direction" form:"Direction" schema:"Direction"`
	ForwardedFrom string  `json:"ForwardedFrom" form:"ForwardedFrom" schema:"ForwardedFrom"`
	CallerName    string  `json:"CallerName" form:"CallerName" schema:"CallerName"`
	FromCity      string  `json:"FromCity" form:"FromCity" schema:"FromCity"`
	FromState     string  `json:"FromState" form:"FromState" schema:"FromState"`
	FromZip       string  `json:"FromZip" form:"FromZip" schema:"FromZip"`
	FromCountry   string  `json:"FromCountry" form:"FromCountry" schema:"FromCountry"`
	ToCity        string  `json:"ToCity" form:"ToCity" schema:"ToCity"`
	ToState       string  `json:"ToState" form:"ToState" schema:"ToState"`
	ToZip         string  `json:"ToZip" form:"ToZip" schema:"ToZip"`
	ToCountry     string  `json:"ToCountry" form:"ToCountry" schema:"ToCountry"`
	Caller        string  `json:"Caller" form:"Caller" schema:"Caller"`
	Digits        string  `json:"Digits" form:"Digits" schema:"Digits"`
	SpeechResult  string  `json:"SpeechResult" form:"SpeechResult" schema:"SpeechResult"`
	Confidence    float64 `json:"Confidence" form:"Confidence" schema:"Confidence"`
}

type IncomingPhoneCallForwarding struct {
	FirstName          string `json:"first_name"`
	UserID             int    `json:"user_id"`
	LeadID             int    `json:"lead_id"`
	ForwardPhoneNumber string `json:"forward_phone_number"`
}

type IncomingPhoneCallDialStatus struct {
	Called           string `json:"called" form:"called" schema:"called"`
	ToState          string `json:"to_state" form:"to_state" schema:"to_state"`
	DialCallStatus   string `json:"dial_call_status" form:"dial_call_status" schema:"dial_call_status"`
	CallerCountry    string `json:"caller_country" form:"caller_country" schema:"caller_country"`
	Direction        string `json:"direction" form:"direction" schema:"direction"`
	CallerState      string `json:"caller_state" form:"caller_state" schema:"caller_state"`
	ToZip            string `json:"to_zip" form:"to_zip" schema:"to_zip"`
	DialCallSid      string `json:"dial_call_sid" form:"dial_call_sid" schema:"dial_call_sid"`
	CallSid          string `json:"call_sid" form:"call_sid" schema:"call_sid"`
	To               string `json:"to" form:"to" schema:"to"`
	CallerZip        string `json:"caller_zip" form:"caller_zip" schema:"caller_zip"`
	ToCountry        string `json:"to_country" form:"to_country" schema:"to_country"`
	CalledZip        string `json:"called_zip" form:"called_zip" schema:"called_zip"`
	ApiVersion       string `json:"api_version" form:"api_version" schema:"api_version"`
	CalledCity       string `json:"called_city" form:"called_city" schema:"called_city"`
	CallStatus       string `json:"call_status" form:"call_status" schema:"call_status"`
	From             string `json:"from" form:"from" schema:"from"`
	DialBridged      bool   `json:"dial_bridged" form:"dial_bridged" schema:"dial_bridged"`
	AccountSid       string `json:"account_sid" form:"account_sid" schema:"account_sid"`
	DialCallDuration int    `json:"dial_call_duration" form:"dial_call_duration" schema:"dial_call_duration"`
	CalledCountry    string `json:"called_country" form:"called_country" schema:"called_country"`
	CallerCity       string `json:"caller_city" form:"caller_city" schema:"caller_city"`
	ToCity           string `json:"to_city" form:"to_city" schema:"to_city"`
	FromCountry      string `json:"from_country" form:"from_country" schema:"from_country"`
	Caller           string `json:"caller" form:"caller" schema:"caller"`
	FromCity         string `json:"from_city" form:"from_city" schema:"from_city"`
	CalledState      string `json:"called_state" form:"called_state" schema:"called_state"`
	FromZip          string `json:"from_zip" form:"from_zip" schema:"from_zip"`
	FromState        string `json:"from_state" form:"from_state" schema:"from_state"`
	RecordingURL     string `json:"recording_url" form:"recording_url" schema:"recording_url"`
}

type WebsiteContext struct {
	PageTitle                    string             `json:"page_title" form:"page_title"`
	MetaDescription              string             `json:"meta_description" form:"meta_description"`
	SiteName                     string             `json:"site_name" form:"site_name"`
	StaticPath                   string             `json:"static_path" form:"static_path"`
	MediaPath                    string             `json:"media_path" form:"media_path"`
	PhoneNumber                  string             `json:"phone_number" form:"phone_number"`
	CurrentYear                  int                `json:"current_year" form:"current_year"`
	GoogleAnalyticsID            string             `json:"google_analytics_id" form:"google_analytics_id"`
	FacebookDataSetID            string             `json:"facebook_data_set_id" form:"facebook_data_set_id"`
	CompanyName                  string             `json:"company_name" form:"company_name"`
	PagePath                     string             `json:"page_path" form:"page_path"`
	Nonce                        string             `json:"nonce" form:"nonce"`
	Features                     []string           `json:"features" form:"features"`
	CSRFToken                    string             `json:"csrf_token" form:"csrf_token"`
	EventTypes                   []models.EventType `json:"event_types" form:"event_types"`
	VenueTypes                   []models.VenueType `json:"venue_types" form:"venue_types"`
	ExternalID                   string             `json:"external_id" form:"external_id"`
	GoogleAdsID                  string             `json:"google_ads_id"`
	GoogleAdsCallConversionLabel string             `json:"google_ads_call_conversion_label"`
	LeadEventName                string             `json:"lead_event_name"`
	LeadGeneratedEventName       string             `json:"lead_generated_event_name"`
	DefaultCurrency              string             `json:"default_currency"`
	YovaHeroImage                string             `json:"yova_hero_image"`
	YovaMostPopularPackage       string             `json:"yova_most_popular_package"`
	YovaBasicPackage             string             `json:"yova_basic_package"`
	YovaOpenBarPackage           string             `json:"yova_open_bar_package"`
	IsMobile                     bool               `json:"is_bool"`
	DefaultLeadValue             float64            `json:"default_lead_value"`
}

type FacebookUserData struct {
	Phone           string `json:"ph,omitempty" form:"ph,omitempty" schema:"ph,omitempty"`
	Email           string `json:"em,omitempty" form:"em,omitempty" schema:"em,omitempty"`
	ClientIPAddress string `json:"client_ip_address,omitempty" form:"client_ip_address,omitempty" schema:"client_ip_address,omitempty"`
	ClientUserAgent string `json:"client_user_agent,omitempty" form:"client_user_agent,omitempty" schema:"client_user_agent,omitempty"`
	FBC             string `json:"fbc,omitempty" form:"fbc,omitempty" schema:"fbc,omitempty"`
	FBP             string `json:"fbp,omitempty" form:"fbp,omitempty" schema:"fbp,omitempty"`
	ExternalID      string `json:"external_id,omitempty" form:"external_id,omitempty" schema:"external_id,omitempty"`
	LeadID          int64  `json:"lead_id,omitempty" form:"lead_id,omitempty" schema:"lead_id,omitempty"`
}

type FacebookCustomData struct {
	Value           string `json:"value,omitempty" form:"value,omitempty" schema:"value,omitempty"`
	Currency        string `json:"currency,omitempty" form:"currency,omitempty" schema:"currency,omitempty"`
	EventSource     string `json:"event_source,omitempty" form:"event_source,omitempty" schema:"event_source,omitempty"`
	LeadEventSource string `json:"lead_event_source,omitempty" form:"lead_event_source,omitempty" schema:"lead_event_source,omitempty"`
}

type FacebookEventData struct {
	EventName      string             `json:"event_name,omitempty" form:"event_name,omitempty" schema:"event_name,omitempty"`
	EventTime      int64              `json:"event_time,omitempty" form:"event_time,omitempty" schema:"event_time,omitempty"`
	EventID        string             `json:"event_id,omitempty" form:"event_id,omitempty" schema:"event_id,omitempty"`
	ActionSource   string             `json:"action_source,omitempty" form:"action_source,omitempty" schema:"action_source,omitempty"`
	EventSourceURL string             `json:"event_source_url,omitempty" form:"event_source_url,omitempty" schema:"event_source_url,omitempty"`
	UserData       FacebookUserData   `json:"user_data,omitempty" form:"user_data,omitempty" schema:"user_data,omitempty"`
	CustomData     FacebookCustomData `json:"custom_data,omitempty" form:"custom_data,omitempty" schema:"custom_data,omitempty"`
}

type FacebookPayload struct {
	Data []FacebookEventData `json:"data,omitempty" form:"data,omitempty" schema:"data,omitempty"`
}

type GoogleEventParamsLead struct {
	GCLID         string  `json:"gclid,omitempty" form:"gclid" schema:"gclid"`
	TransactionID string  `json:"transaction_id,omitempty" form:"transaction_id" schema:"transaction_id"`
	Value         float64 `json:"value,omitempty" form:"value" schema:"value"`
	Currency      string  `json:"currency,omitempty" form:"currency" schema:"currency"`
	CampaignID    string  `json:"campaign_id,omitempty" form:"campaign_id" schema:"campaign_id"`
	Campaign      string  `json:"campaign,omitempty" form:"campaign" schema:"campaign"`
	Source        string  `json:"source,omitempty" form:"source" schema:"source"`
	Medium        string  `json:"medium,omitempty" form:"medium" schema:"medium"`
	Term          string  `json:"term,omitempty" form:"term" schema:"term"`
}

type GoogleEventLead struct {
	Name   string                `json:"name,omitempty" form:"name" schema:"name"`
	Params GoogleEventParamsLead `json:"params,omitempty" form:"params" schema:"params"`
}

type GooglePayload struct {
	ClientID string            `json:"client_id,omitempty" form:"client_id" schema:"client_id"`
	UserId   string            `json:"userId,omitempty" form:"userId" schema:"userId"`
	Events   []GoogleEventLead `json:"events,omitempty" form:"events" schema:"events"`
	UserData GoogleUserData    `json:"user_data,omitempty" form:"user_data" schema:"user_data"`
}

type GoogleUserData struct {
	Sha256EmailAddress []string            `json:"sha256_email_address,omitempty" form:"sha256_email_address,omitempty" schema:"sha256_email_address,omitempty"`
	Sha256PhoneNumber  []string            `json:"sha256_phone_number,omitempty" form:"sha256_phone_number,omitempty" schema:"sha256_phone_number,omitempty"`
	Address            []GoogleUserAddress `json:"address,omitempty" form:"address,omitempty" schema:"address,omitempty"`
}

type GoogleUserAddress struct {
	Sha256FullName string `json:"sha256_full_name,omitempty" form:"sha256_full_name,omitempty" schema:"sha256_full_name,omitempty"`
	Sha256Street   string `json:"sha256_street,omitempty" form:"sha256_street,omitempty" schema:"sha256_street,omitempty"`
	City           string `json:"city,omitempty" form:"city,omitempty" schema:"city,omitempty"`
	Region         string `json:"region,omitempty" form:"region,omitempty" schema:"region,omitempty"`
	PostalCode     string `json:"postal_code,omitempty" form:"postal_code,omitempty" schema:"postal_code,omitempty"`
	Country        string `json:"country,omitempty" form:"country,omitempty" schema:"country,omitempty"`
}

type EventList struct {
	LeadID    int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	EventID   int     `json:"event_id" form:"event_id" schema:"event_id"`
	Amount    float64 `json:"amount" form:"amount" schema:"amount"`
	EventTime string  `json:"event_time" form:"event_time" schema:"event_time"`
	LeadName  string  `json:"lead_name" form:"lead_name" schema:"lead_name"`
	Bartender string  `json:"bartender" form:"bartender" schema:"bartender"`
	EventType string  `json:"event_type" form:"event_type" schema:"event_type"`
	VenueType string  `json:"venue_type" form:"venue_type" schema:"venue_type"`
	Guests    int     `json:"guests" form:"guests" schema:"guests"`
}

type EventDetails struct {
	EventID     int `json:"event_id" form:"event_id" schema:"event_id"`
	LeadID      int `json:"lead_id" form:"lead_id" schema:"lead_id"`
	BartenderID int `json:"bartender_id" form:"bartender_id" schema:"bartender_id"`
	EventTypeID int `json:"event_type" form:"event_type" schema:"event_type"`
	VenueTypeID int `json:"venue_type" form:"venue_type" schema:"venue_type"`

	StreetAddress string  `json:"street_address" form:"street_address" schema:"street_address"`
	City          string  `json:"city" form:"city" schema:"city"`
	ZipCode       string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	StartTime     int64   `json:"start_time" form:"start_time" schema:"start_time"`
	EndTime       int64   `json:"end_time" form:"end_time" schema:"end_time"`
	DateCreated   int64   `json:"date_created" form:"date_created" schema:"date_created"`
	DatePaid      int64   `json:"date_paid" form:"date_paid" schema:"date_paid"`
	Amount        float64 `json:"amount" form:"amount" schema:"amount"`
	Guests        int     `json:"guests" form:"guests" schema:"guests"`
}

type EventForm struct {
	CSRFToken *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	EventID   *int    `json:"event_id" form:"event_id" schema:"event_id"`
	LeadID    *int    `json:"lead_id" form:"lead_id" schema:"lead_id"`

	BartenderID *int `json:"bartender_id" form:"bartender_id" schema:"bartender_id"`
	EventTypeID *int `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	VenueTypeID *int `json:"venue_type_id" form:"venue_type_id" schema:"venue_type_id"`

	StreetAddress *string  `json:"street_address" form:"street_address" schema:"street_address"`
	City          *string  `json:"city" form:"city" schema:"city"`
	ZipCode       *string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	StartTime     *int64   `json:"start_time" form:"start_time" schema:"start_time"`
	EndTime       *int64   `json:"end_time" form:"end_time" schema:"end_time"`
	DateCreated   *int64   `json:"date_created" form:"date_created" schema:"date_created"`
	DatePaid      *int64   `json:"date_paid" form:"date_paid" schema:"date_paid"`
	Amount        *float64 `json:"amount" form:"amount" schema:"amount"`
	Tip           *float64 `json:"tip" form:"tip" schema:"tip"`
	Guests        *int     `json:"guests" form:"guests" schema:"guests"`
}

type FacebookInstantFormLead struct {
	ID               string `json:"id"`
	CreatedTime      string `json:"created_time"`
	AdID             string `json:"ad_id"`
	AdName           string `json:"ad_name"`
	AdsetID          string `json:"adset_id"`
	AdsetName        string `json:"adset_name"`
	CampaignID       string `json:"campaign_id"`
	CampaignName     string `json:"campaign_name"`
	FormID           string `json:"form_id"`
	FormName         string `json:"form_name"`
	IsOrganic        string `json:"is_organic"`
	Platform         string `json:"platform"`
	FullName         string `json:"whats_your_full_name"`
	PhoneNumber      string `json:"whats_the_best_phone_number_to_reach_you_at"`
	EventDescription string `json:"give_us_a_brief_description_of_your_event"`
	Email            string `json:"email"`
	IsQualified      string `json:"is_qualified"`
	IsQuality        string `json:"is_quality"`
	IsConverted      string `json:"is_converted"`
}

type ConversionReporting struct {
	LeadID            int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	Email             string  `json:"email" form:"email" schema:"email"`
	PhoneNumber       string  `json:"phone_number" form:"phone_number" schema:"phone_number"`
	CampaignName      string  `json:"campaign_name" form:"campaign_name" schema:"campaign_name"`
	CampaignID        int64   `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	LandingPage       string  `json:"landing_page" form:"landing_page" schema:"landing_page"`
	IP                string  `json:"ip" form:"ip" schema:"ip"`
	FacebookClickID   string  `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID  string  `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	UserAgent         string  `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ExternalID        string  `json:"external_id" form:"external_id" schema:"external_id"`
	ClickID           string  `json:"click_id" form:"click_id" schema:"click_id"`
	GoogleClientID    string  `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	ReferralLeadID    int     `json:"referral_lead_id" form:"referral_lead_id" schema:"referral_lead_id"`
	InstantFormLeadID int64   `json:"instant_form_lead_id" form:"instant_form_lead_id" schema:"instant_form_lead_id"`
	EventID           int     `json:"event_id" form:"event_id" schema:"event_id"`
	Revenue           float64 `json:"revenue" form:"revenue" schema:"revenue"`
}

type LeadQuoteList struct {
	LeadID    int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	QuoteID   int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	EventType string  `json:"event_type" form:"event_type" schema:"event_type"`
	VenueType string  `json:"venue_type" form:"venue_type" schema:"venue_type"`
	Guests    int     `json:"guests" form:"guests" schema:"guests"`
	EventDate string  `json:"event_date" form:"event_date" schema:"event_date"`
	Amount    float64 `json:"amount" form:"amount" schema:"amount"`
}

type LeadQuoteForm struct {
	CSRFToken          *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	QuoteID            *int    `json:"quote_id" form:"quote_id" schema:"quote_id"`
	LeadID             *int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	ExternalID         *string `json:"external_id" form:"external_id" schema:"external_id"`
	EventTypeID        *int    `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	VenueTypeID        *int    `json:"venue_type_id" form:"venue_type_id" schema:"venue_type_id"`
	NumberOfBartenders *int    `json:"number_of_bartenders" form:"number_of_bartenders" schema:"number_of_bartenders"`
	Guests             *int    `json:"guests" form:"guests" schema:"guests"`
	Hours              *int    `json:"hours" form:"hours" schema:"hours"`
	EventDate          *int64  `json:"event_date" form:"event_date" schema:"event_date"`
}

type QuoteDetails struct {
	LeadID           int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	ExternalID       string  `json:"external_id" form:"external_id" schema:"external_id"`
	FullName         string  `json:"full_name" form:"full_name" schema:"full_name"`
	Email            string  `json:"email" form:"email" schema:"email"`
	PhoneNumber      string  `json:"phone_number" form:"phone_number" schema:"phone_number"`
	StripeCustomerID string  `json:"stripe_customer_id" form:"stripe_customer_id" schema:"stripe_customer_id"`
	EventDate        int64   `json:"event_date" form:"event_date" schema:"event_date"`
	InvoiceID        int     `json:"invoice_id" form:"invoice_id" schema:"invoice_id"`
	Amount           float64 `json:"amount" form:"amount" schema:"amount"`
}

type ExternalQuoteDetails struct {
	QuoteID             int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	ExternalID          string  `json:"external_id" form:"external_id" schema:"external_id"`
	Amount              float64 `json:"amount" form:"amount" schema:"amount"`
	Deposit             float64 `json:"deposit" form:"deposit" schema:"deposit"`
	RemainingAmount     float64 `json:"remaining_amount" form:"remaining_amount" schema:"remaining_amount"`
	EventType           string  `json:"event_type" form:"event_type" schema:"event_type"`
	VenueType           string  `json:"venue_type" form:"venue_type" schema:"venue_type"`
	Guests              int     `json:"guests" form:"guests" schema:"guests"`
	Hours               int     `json:"hours" form:"hours" schema:"hours"`
	EventDate           string  `json:"event_date" form:"event_date" schema:"event_date"`
	NumberOfBartenders  int     `json:"number_of_bartenders" form:"number_of_bartenders" schema:"number_of_bartenders"`
	FullName            string  `json:"full_name" form:"full_name" schema:"full_name"`
	PhoneNumber         string  `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email               string  `json:"email" form:"email" schema:"email"`
	DepositInvoiceURL   string  `json:"deposit_invoice_url" form:"deposit_invoice_url" schema:"deposit_invoice_url"`
	FullInvoiceURL      string  `json:"full_invoice_url" form:"full_invoice_url" schema:"full_invoice_url"`
	RemainingInvoiceURL string  `json:"remaining_invoice_url" form:"remaining_invoice_url" schema:"remaining_invoice_url"`
	IsDepositPaid       bool    `json:"is_deposit_paid" form:"is_deposit_paid" schema:"is_deposit_paid"`
}

type CreateInvoiceParams struct {
	StripeCustomerID string
	Email            string
	FullName         string
	DueDate          int64
	Quote            float64
}

type InvoiceQuoteDetails struct {
	FullName         string  `json:"full_name" form:"full_name" schema:"full_name"`
	StripeCustomerID string  `json:"stripe_customer_id" form:"stripe_customer_id" schema:"stripe_customer_id"`
	Amount           float64 `json:"amount" form:"amount" schema:"amount"`
	EventTypeID      int     `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	VenueTypeID      int     `json:"venue_type_id" form:"venue_type_id" schema:"venue_type_id"`
	LeadID           int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	QuoteID          int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	Guests           int     `json:"guests" form:"guests" schema:"guests"`
	PhoneNumber      string  `json:"phone_number" form:"phone_number" schema:"phone_number"`
	EventDate        int64   `json:"event_date" form:"event_date" schema:"event_date"`
}

type LeadQuoteInvoice struct {
	StripeCustomerID      string  `json:"stripe_customer_id" form:"stripe_customer_id" schema:"stripe_customer_id"`
	StripeInvoiceID       string  `json:"stripe_invoice_id" form:"stripe_invoice_id" schema:"stripe_invoice_id"`
	Amount                float64 `json:"amount" form:"amount" schema:"amount"`
	DueDate               int64   `json:"due_date" form:"due_date" schema:"due_date"`
	InvoiceTypeMultiplier float64 `json:"invoice_type_multiplier" form:"invoice_type_multiplier" schema:"invoice_type_multiplier"`
	InvoiceTypeID         int     `json:"invoice_type_id" form:"invoice_type_id" schema:"invoice_type_id"`
}

type QuoteServiceList struct {
	QuoteServiceID int     `json:"quote_service_id" form:"quote_service_id" schema:"quote_service_id"`
	ServiceID      int     `json:"service_id" form:"service_id" schema:"service_id"`
	QuoteID        int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	Service        string  `json:"service" form:"service" schema:"service"`
	Units          int     `json:"units" form:"units" schema:"units"`
	PricePerUnit   float64 `json:"price_per_unit" form:"price_per_unit" schema:"price_per_unit"`
	Total          float64 `json:"total" form:"total" schema:"total"`
}

type QuoteServiceForm struct {
	CSRFToken      *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	QuoteServiceID *int     `json:"quote_service_id" form:"quote_service_id" schema:"quote_service_id"`
	ServiceID      *int     `json:"service_id" form:"service_id" schema:"service_id"`
	QuoteID        *int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	Units          *int     `json:"units" form:"units" schema:"units"`
	PricePerUnit   *float64 `json:"price_per_unit" form:"price_per_unit" schema:"price_per_unit"`
}

type ServiceForm struct {
	CSRFToken      *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	ServiceID      *int     `json:"service_id" form:"service_id" schema:"service_id"`
	Service        *string  `json:"service" form:"service" schema:"service"`
	SuggestedPrice *float64 `json:"suggested_price" form:"suggested_price" schema:"suggested_price"`
}

type FrontendMessage struct {
	ClientName  string `json:"client_name"`
	UserName    string `json:"user_name"`
	DateCreated string `json:"date_created"`
	Message     string `json:"message"`
	IsInbound   bool   `json:"is_inbound"`
	MessageID   int    `json:"message_id"`
	LeadID      int    `json:"lead_id"`
	IsRead      bool   `json:"is_read"`
}

type FrontendNote struct {
	UserName  string `json:"user_name"`
	DateAdded string `json:"date_added"`
	Note      string `json:"note"`
}

type SetSMSToReadForm struct {
	MessageID *int    `json:"message_id" form:"message_id" schema:"message_id"`
	LeadID    *int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	IsRead    *bool   `json:"is_read" form:"is_read" schema:"is_read"`
	CSRFToken *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
}
