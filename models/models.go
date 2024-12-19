package models

type EventType struct {
	EventTypeID int    `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	Name        string `json:"name" form:"name" schema:"name"`
}

type VenueType struct {
	VenueTypeID int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	Name        string `json:"name" form:"name" schema:"name"`
}

type User struct {
	UserID      int    `json:"user_id" form:"user_id" schema:"user_id"`
	Username    string `json:"username" form:"username" schema:"username"`
	PhoneNumber string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Password    string `json:"password" form:"password" schema:"password"`
	UserRoleID  int    `json:"user_role_id" form:"user_role_id" schema:"user_role_id"`
	FirstName   string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName    string `json:"last_name" form:"last_name" schema:"last_name"`
}

type Lead struct {
	LeadID             int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FirstName          string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName           string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber        string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email              string `json:"email" form:"email" schema:"email"`
	CreatedAt          int64  `json:"created_at" form:"created_at" schema:"created_at"`
	EventTypeID        int    `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	VenueTypeID        int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	Guests             int    `json:"guests" form:"guests" schema:"guests"`
	Message            string `json:"message" form:"message" schema:"message"`
	OptInTextMessaging bool   `json:"opt_in_text_messaging" form:"opt_in_text_messaging" schema:"opt_in_text_messaging"`
}

type LeadMarketing struct {
	LeadMarketingID  int64  `json:"lead_marketing_id"`
	LeadID           int64  `json:"lead_id"`
	Source           string `json:"source"`
	Medium           string `json:"medium"`
	Channel          string `json:"channel"`
	LandingPage      string `json:"landing_page"`
	Longitude        string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         string `json:"latitude" form:"latitude" schema:"latitude"`
	Keyword          string `json:"keyword"`
	Referrer         string `json:"referrer"`
	ClickID          string `json:"click_id"`
	CampaignID       int64  `json:"campaign_id"`
	AdCampaign       string `json:"ad_campaign"`
	AdGroupID        int64  `json:"ad_group_id"`
	AdGroupName      string `json:"ad_group_name"`
	AdSetID          int64  `json:"ad_set_id"`
	AdSetName        string `json:"ad_set_name"`
	AdID             int64  `json:"ad_id"`
	AdHeadline       int64  `json:"ad_headline"`
	Language         string `json:"language"`
	OS               string `json:"os"`
	UserAgent        string `json:"user_agent"`
	ButtonClicked    string `json:"button_clicked"`
	DeviceType       string `json:"device_type"`
	IP               string `json:"ip"`
	ExternalID       string `json:"external_id"`
	GoogleClientID   string `json:"google_client_id"`
	FacebookClickID  string `json:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id"`
	CSRFSecret       string `json:"csrf_secret"`
}

type CSRFToken struct {
	CSRFTokenID int    `json:"csrf_token_id"`
	ExpiryTime  int64  `json:"expiry_time"`
	Token       string `json:"token"`
	IsUsed      bool   `json:"is_used"`
}

type Message struct {
	MessageID   int    `json:"message_id"`
	ExternalID  string `json:"external_id"`
	UserID      int    `json:"user_id"`
	LeadID      int    `json:"lead_id"`
	Text        string `json:"text"`
	DateCreated int64  `json:"date_created"`
	TextFrom    string `json:"text_from"`
	TextTo      string `json:"text_to"`
	IsInbound   bool   `json:"is_inbound"`
	Status      string `json:"status" form:"status" schema:"status"`
}

type PhoneCall struct {
	PhoneCallID  int    `json:"phone_call_id" form:"phone_call_id" schema:"phone_call_id"`
	ExternalID   string `json:"external_id" form:"external_id" schema:"external_id"`
	UserID       int    `json:"user_id" form:"user_id" schema:"user_id"`
	LeadID       int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	CallDuration int    `json:"call_duration" form:"call_duration" schema:"call_duration"`
	DateCreated  int64  `json:"date_created" form:"date_created" schema:"date_created"`
	CallFrom     string `json:"call_from" form:"call_from" schema:"call_from"`
	CallTo       string `json:"call_to" form:"call_to" schema:"call_to"`
	IsInbound    bool   `json:"is_inbound" form:"is_inbound" schema:"is_inbound"`
	RecordingURL string `json:"recording_url" form:"recording_url" schema:"recording_url"`
	Status       string `json:"status" form:"status" schema:"status"`
}

type Session struct {
	SessionID   int    `json:"session_id" form:"session_id" schema:"session_id"`
	UserID      int    `json:"user_id" form:"user_id" schema:"user_id"`
	CSRFSecret  string `json:"csrf_secret" form:"csrf_secret" schema:"csrf_secret"`
	ExternalID  string `json:"external_id" form:"external_id" schema:"external_id"`
	DateCreated int64  `json:"date_created" form:"date_created" schema:"date_created"`
	DateExpires int64  `json:"date_expires" form:"date_expires" schema:"date_expires"`
}

type UserRole struct {
	RoleID int    `json:"role_id" form:"role_id" schema:"role_id"`
	Role   string `json:"role" form:"role" schema:"role"`
}

type Package struct {
	PackageID              int     `json:"package_id" form:"package_id" schema:"package_id"`
	PackageTypeID          int     `json:"package_type_id" form:"package_type_id" schema:"package_type_id"`
	AlcoholSegmentID       int     `json:"alcohol_segment_id" form:"alcohol_segment_id" schema:"alcohol_segment_id"`
	LeadID                 int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	Price                  float64 `json:"price" form:"price" schema:"price"`
	Guests                 int     `json:"guests" form:"guests" schema:"guests"`
	Hours                  int     `json:"hours" form:"hours" schema:"hours"`
	WillProvideLiquor      bool    `json:"will_provide_liquor" form:"will_provide_liquor" schema:"will_provide_liquor"`
	WillProvideBeerAndWine bool    `json:"will_provide_beer_and_wine" form:"will_provide_beer_and_wine" schema:"will_provide_beer_and_wine"`
	WillProvideMixers      bool    `json:"will_provide_mixers" form:"will_provide_mixers" schema:"will_provide_mixers"`
	WillProvideJuices      bool    `json:"will_provide_juices" form:"will_provide_juices" schema:"will_provide_juices"`
	WillProvideSoftDrinks  bool    `json:"will_provide_soft_drinks" form:"will_provide_soft_drinks" schema:"will_provide_soft_drinks"`
	WillProvideCups        bool    `json:"will_provide_cups" form:"will_provide_cups" schema:"will_provide_cups"`
	WillProvideIce         bool    `json:"will_provide_ice" form:"will_provide_ice" schema:"will_provide_ice"`
	WillRequireGlassware   bool    `json:"will_require_glassware" form:"will_require_glassware" schema:"will_require_glassware"`
	WillRequireMobileBar   bool    `json:"will_require_mobile_bar" form:"will_require_mobile_bar" schema:"will_require_mobile_bar"`
	NumBars                int     `json:"num_bars" form:"num_bars" schema:"num_bars"`
	DateCreated            int64   `json:"date_created" form:"date_created" schema:"date_created"`
	DateUpdated            int64   `json:"date_updated" form:"date_updated" schema:"date_updated"`
}

// Full Open Bar, Partial Open Bar
type PackageType struct {
	PackageTypeID     int     `json:"package_type_id" form:"package_type_id" schema:"package_type_id"`
	Name              string  `json:"name" form:"name" schema:"name"`
	PriceModification float64 `json:"price_modification" form:"price_modification" schema:"price_modification"`
}

// Top Shelf, Premium, Standard
type AlcoholSegment struct {
	AlcoholSegmentID  int     `json:"package_type_id" form:"package_type_id" schema:"package_type_id"`
	PriceModification float64 `json:"price_modification" form:"price_modification" schema:"price_modification"`
}
