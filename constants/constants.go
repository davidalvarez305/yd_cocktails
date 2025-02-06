package constants

import (
	"os"
)

const (
	UserAdminRoleID int = 1

	DavidUserID int = 1

	TimeZone string = "America/New_York"

	LeadGeneratedEventName   string = "generate_lead"
	EventConversionEventName string = "event"
	LeadEventName            string = "Lead"
	EventSourceCRM           string = "crm"

	DefaultCurrency          string  = "USD"
	DefaultLeadValue         float64 = 150.00
	InvoicePaymentDueInHours int     = 48

	CallConversionDuration int = 15

	PerGuestBartendingRatio       float64 = 70.00
	BartendingRate                float64 = 70.00
	BarRentalCost                 float64 = 200.00
	PerPersonAlcoholFee           float64 = 10.00
	PerPersonWineFee              float64 = 2.00
	PerPersonBeerFee              float64 = 3.00
	PerPersonMixersFee            float64 = 3.00
	PerPersonGarnishFee           float64 = 1.00
	PerPersonJuicesFee            float64 = 1.00
	PerPersonSoftDrinksFee        float64 = 2.00
	PerPersonCupsStrawsNapkinsFee float64 = 2.00
	PerPersonIceFee               float64 = 2.00
	PerPersonGlasswareFee         float64 = 3.00
	DepositPercentageAmount       float64 = 0.25

	DepositInvoiceTypeID   int = 1
	RemainingInvoiceTypeID int = 2
	FullInvoiceTypeID      int = 3

	SocialMediaAdsMedium  = "paid"
	SocialMediaAdsChannel = "social"
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
	TwilioCallbackWebhook         string
	CompanyName                   string
	SiteName                      string
	SessionName                   string
	LeadsPerPage                  int
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
	SessionName = "yd_vending_sessions"
	LeadsPerPage = 10
	TwilioCallbackWebhook = "/call/inbound/end"
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

	NotificationSubscribers = []string{DavidPhoneNumber, YovaPhoneNumber}
}

var TEMPLATES_DIR = "./templates/"
var LOCAL_FILES_DIR = "./local_files/"
var WEBSITE_TEMPLATES_DIR = TEMPLATES_DIR + "website/"
var CRM_TEMPLATES_DIR = TEMPLATES_DIR + "crm/"
var PARTIAL_TEMPLATES_DIR = TEMPLATES_DIR + "partials/"
var EXTERNAL_TEMPLATES_DIR = TEMPLATES_DIR + "external/"
