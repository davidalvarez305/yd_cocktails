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
	OptInTextMessaging bool   `json:"opt_in_text_messaging" form:"opt_in_text_messaging" schema:"opt_in_text_messaging"`
	CreatedAt          int64  `json:"created_at" form:"created_at" schema:"created_at"`

	// Nullable
	Email   string `json:"email" form:"email" schema:"email"`
	Message string `json:"message" form:"message" schema:"message"`
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

type Event struct {
	EventID       int     `json:"event_id" form:"event_id" schema:"event_id"`
	BartenderID   int     `json:"bartender_id" form:"bartender_id" schema:"bartender_id"`
	LeadID        int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	StreetAddress string  `json:"street_address" form:"street_address" schema:"street_address"`
	City          string  `json:"city" form:"city" schema:"city"`
	State         string  `json:"state" form:"state" schema:"state"`
	ZipCode       string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	StartTime     int64   `json:"start_time" form:"start_time" schema:"start_time"`
	EndTime       int64   `json:"end_time" form:"end_time" schema:"end_time"`
	DateCreated   int64   `json:"date_created" form:"date_created" schema:"date_created"`
	DatePaid      int64   `json:"date_paid" form:"date_paid" schema:"date_paid"`
	Amount        float64 `json:"amount" form:"amount" schema:"amount"`
	Tip           float64 `json:"tip" form:"tip" schema:"tip"`
	EventTypeID   int     `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	VenueTypeID   int     `json:"venue_type_id" form:"venue_type_id" schema:"venue_type_id"`
	Guests        int     `json:"guests" form:"guests" schema:"guests"`
}

type EventCocktail struct {
	EventCocktailID int `json:"event_cocktail_id" form:"event_cocktail_id" schema:"event_cocktail_id"`
	CocktailID      int `json:"cocktail_id" form:"cocktail_id" schema:"cocktail_id"`
	EventID         int `json:"event_id" form:"event_id" schema:"event_id"`
}

type EventExpense struct {
	EventExpenseID     int     `json:"event_expense_id" form:"event_expense_id" schema:"event_expense_id"`
	EventExpenseTypeID int     `json:"event_expense_type_id" form:"event_expense_type_id" schema:"event_expense_type_id"`
	Name               string  `json:"name" form:"name" schema:"name"`
	Amount             float64 `json:"amount" form:"amount" schema:"amount"`
}

type EventExpenseType struct {
	EventExpenseTypeID int    `json:"event_expense_type_id" form:"event_expense_type_id" schema:"event_expense_type_id"`
	Type               string `json:"type" form:"type" schema:"type"`
}

type Cocktail struct {
	CocktailID   int    `json:"cocktail_id" form:"cocktail_id" schema:"cocktail_id"`
	Name         string `json:"name" form:"name" schema:"name"`
	Instructions string `json:"instructions" form:"instructions" schema:"instructions"`
	GlassType    string `json:"glass_type" form:"glass_type" schema:"glass_type"`
}

type Ingredient struct {
	IngredientID int    `json:"ingredient_id" form:"ingredient_id" schema:"ingredient_id"`
	Name         string `json:"name" form:"name" schema:"name"`
	Category     string `json:"category" form:"category" schema:"category"` // e.g., Liquor, Mixer, Garnish
}

type Unit struct {
	UnitID       int    `json:"unit_id" form:"unit_id" schema:"unit_id"`
	Name         string `json:"name" form:"name" schema:"name"`                         // e.g., Ounce, Teaspoon
	Abbreviation string `json:"abbreviation" form:"abbreviation" schema:"abbreviation"` // e.g., oz, tsp, part, dash, splash
}

type CocktailIngredient struct {
	CocktailID   int     `json:"cocktail_id" form:"cocktail_id" schema:"cocktail_id"`
	IngredientID int     `json:"ingredient_id" form:"ingredient_id" schema:"ingredient_id"`
	UnitID       int     `json:"unit_id" form:"unit_id" schema:"unit_id"`
	Amount       float64 `json:"amount" form:"amount" schema:"amount"` // Amount of the ingredient
}
