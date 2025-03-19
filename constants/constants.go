package constants

import (
	"os"
)

const (
	UserAdminRoleID int = 1

	DavidUserID int = 1

	TimeZone string = "America/New_York"

	StructSpreadsheetHeaderTag string = "spreadsheet_header"

	LeadGeneratedEventName   string = "generate_lead"
	EventConversionEventName string = "event"
	LeadEventName            string = "Lead"
	EventSourceCRM           string = "crm"

	DefaultCurrency          string  = "USD"
	DefaultLeadValue         float64 = 150.00
	InvoicePaymentDueInHours int     = 48

	CallConversionDuration int = 15

	DepositInvoiceTypeID   int = 1
	RemainingInvoiceTypeID int = 2
	FullInvoiceTypeID      int = 3

	OpenInvoiceStatusID int = 1
	VoidInvoiceStatusID int = 2
	PaidInvoiceStatusID int = 3

	AlcoholServiceTypeID         int = 1
	BarRentalServiceTypeID       int = 2
	CoolerRentalServiceTypeID    int = 3
	BartendingAddOnServiceTypeID int = 4
	BartendingServiceTypeID      int = 5
	ExtraServiceTypeID           int = 8

	CupsStrawsNapkinsServiceID int = 14

	NoInterestLeadInterestID int = 4

	NewLeadStatusID      int = 1
	ArchivedLeadStatusID int = 7

	InitialContactActionID int = 1
	FirstFollowUpActionID  int = 3

	SocialMediaAdsMedium  = "paid"
	SocialMediaAdsChannel = "social"

	SessionName                    = "yd_vending_sessions"
	LeadsPerPage                   = 10
	TwilioCallbackWebhook          = "/call/inbound/end"
	TwilioRecordingCallbackWebhook = "/call/inbound/recording-callback"
	TwilioAmdCallbackWebhook       = "/call/inbound/amd"
)

var (
	Production                    bool
	StrikeAPIKey                  string
	StripeWebhookSecret           string
	FacebookAccessToken           string
	FacebookDatasetID             string
	GoogleAnalyticsID             string
	GoogleAdsID                   string
	GoogleAdsCallConversionLabel  string
	GoogleAnalyticsAPISecretKey   string
	GoogleRefreshToken            string
	GoogleJSONPath                string
	PostgresHost                  string
	PostgresPort                  string
	PostgresUser                  string
	PostgresPassword              string
	PostgresDBName                string
	DavidPhoneNumber              string
	DavidEmail                    string
	YovaEmail                     string
	YovaPhoneNumber               string
	ServerPort                    string
	RootDomain                    string
	AWSStorageBucket              string
	AWSS3BucketName               string
	AWSRegion                     string
	CookieName                    string
	DomainHost                    string
	SecretAESKey                  string
	AuthSecretKey                 string
	EncSecretKey                  string
	TwilioAccountSID              string
	TwilioAuthToken               string
	CompanyName                   string
	SiteName                      string
	CompanyPhoneNumber            string
	SessionLength                 int
	CSRFTokenLength               int
	StaticPath                    string
	MediaPath                     string
	MaxOpenConnections            string
	MaxIdleConnections            string
	MaxConnectionLifetime         string
	CompanyEmail                  string
	NotificationSubscribers       []string
	FacebookLeadsSpreadsheetID    string
	FacebookLeadsSpreadsheetRange string
	OpenAIApiKey                  string
)

func Init() {
	Production = os.Getenv("PRODUCTION") == "1"
	FacebookAccessToken = os.Getenv("FACEBOOK_ACCESS_TOKEN")
	FacebookDatasetID = os.Getenv("FACEBOOK_DATASET_ID")
	GoogleAnalyticsID = os.Getenv("GOOGLE_ANALYTICS_ID")
	GoogleAnalyticsAPISecretKey = os.Getenv("GOOGLE_ANALYTICS_API_KEY")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("PGUSER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDBName = os.Getenv("POSTGRES_DB")
	DavidPhoneNumber = os.Getenv("DAVID_PHONE_NUMBER")
	YovaPhoneNumber = os.Getenv("YOVA_PHONE_NUMBER")
	ServerPort = os.Getenv("SERVER_PORT")
	RootDomain = os.Getenv("ROOT_DOMAIN")
	AWSStorageBucket = os.Getenv("AWS_STORAGE_BUCKET")
	AWSS3BucketName = os.Getenv("AWS_S3_BUCKET_NAME")
	AWSRegion = os.Getenv("AWS_REGION")
	CookieName = os.Getenv("COOKIE_NAME")
	SecretAESKey = os.Getenv("SECRET_AES_KEY")
	AuthSecretKey = os.Getenv("AUTH_SECRET_KEY")
	EncSecretKey = os.Getenv("ENC_SECRET_KEY")
	TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	CompanyName = os.Getenv("COMPANY_NAME")
	SiteName = os.Getenv("SITE_NAME")
	GoogleRefreshToken = os.Getenv("GOOGLE_REFRESH_TOKEN")
	GoogleJSONPath = "./google.json"
	DavidEmail = os.Getenv("DAVID_EMAIL")
	YovaEmail = os.Getenv("YOVA_EMAIL")
	CompanyPhoneNumber = os.Getenv("COMPANY_PHONE_NUMBER")
	SessionLength = 7
	CSRFTokenLength = 1
	StaticPath = os.Getenv("STATIC_PATH")
	MediaPath = os.Getenv("MEDIA_PATH")
	MaxOpenConnections = os.Getenv("MAX_OPEN_CONNECTIONS")
	MaxIdleConnections = os.Getenv("MAX_IDLE_CONNECTIONS")
	MaxConnectionLifetime = os.Getenv("MAX_CONN_LIFETIME")
	DomainHost = os.Getenv("DOMAIN_HOST")
	CompanyEmail = os.Getenv("COMPANY_EMAIL")
	GoogleAdsID = os.Getenv("GOOGLE_ADS_ID")
	GoogleAdsCallConversionLabel = os.Getenv("GOOGLE_ADS_CALL_CONVERSION_LABEL")
	StrikeAPIKey = os.Getenv("STRIPE_API_KEY")
	StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
	FacebookLeadsSpreadsheetID = os.Getenv("FACEBOOK_LEADS_SPREADSHEET_ID")
	FacebookLeadsSpreadsheetRange = os.Getenv("FACEBOOK_LEADS_SPREADSHEET_RANGE")
	OpenAIApiKey = os.Getenv("OPEN_AI_API_KEY")

	NotificationSubscribers = []string{DavidPhoneNumber, YovaPhoneNumber}
}

var TEMPLATES_DIR = "./templates/"
var LOCAL_FILES_DIR = "./local_files/"
var WEBSITE_TEMPLATES_DIR = TEMPLATES_DIR + "website/"
var CRM_TEMPLATES_DIR = TEMPLATES_DIR + "crm/"
var PARTIAL_TEMPLATES_DIR = TEMPLATES_DIR + "partials/"
var EXTERNAL_TEMPLATES_DIR = TEMPLATES_DIR + "external/"
