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
	UserID             int    `json:"user_id" form:"user_id" schema:"user_id"`
	Username           string `json:"username" form:"username" schema:"username"`
	PhoneNumber        string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	ForwardPhoneNumber string `json:"forward_phone_number" form:"forward_phone_number" schema:"forward_phone_number"`
	Password           string `json:"password" form:"password" schema:"password"`
	UserRoleID         int    `json:"user_role_id" form:"user_role_id" schema:"user_role_id"`
	FirstName          string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName           string `json:"last_name" form:"last_name" schema:"last_name"`
}

type Lead struct {
	LeadID             int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FullName           string `json:"full_name" form:"full_name" schema:"full_name"`
	PhoneNumber        string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	OptInTextMessaging bool   `json:"opt_in_text_messaging" form:"opt_in_text_messaging" schema:"opt_in_text_messaging"`
	CreatedAt          int64  `json:"created_at" form:"created_at" schema:"created_at"`

	// Nullable
	Email   string `json:"email" form:"email" schema:"email"`
	Message string `json:"message" form:"message" schema:"message"`

	// Flow Management
	NextActionID   int `json:"next_action_id" form:"next_action_id" schema:"next_action_id"`
	LeadInterestID int `json:"lead_interest_id" form:"lead_interest_id" schema:"lead_interest_id"`
	LeadStatusID   int `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
}

type LeadMarketing struct {
	LeadMarketingID   int64  `json:"lead_marketing_id" form:"lead_marketing_id" schema:"lead_marketing_id"`
	LeadID            int64  `json:"lead_id" form:"lead_id" schema:"lead_id"`
	Source            string `json:"source" form:"source" schema:"source"`
	Medium            string `json:"medium" form:"medium" schema:"medium"`
	Channel           string `json:"channel" form:"channel" schema:"channel"`
	LandingPage       string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	Longitude         string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude          string `json:"latitude" form:"latitude" schema:"latitude"`
	Keyword           string `json:"keyword" form:"keyword" schema:"keyword"`
	Referrer          string `json:"referrer" form:"referrer" schema:"referrer"`
	ClickID           string `json:"click_id" form:"click_id" schema:"click_id"`
	CampaignID        int64  `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign        string `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID         int64  `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName       string `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID           int64  `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName         string `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID              int64  `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline        int64  `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
	Language          string `json:"language" form:"language" schema:"language"`
	OS                string `json:"os" form:"os" schema:"os"`
	UserAgent         string `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ButtonClicked     string `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	DeviceType        string `json:"device_type" form:"device_type" schema:"device_type"`
	IP                string `json:"ip" form:"ip" schema:"ip"`
	ExternalID        string `json:"external_id" form:"external_id" schema:"external_id"`
	GoogleClientID    string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID   string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID  string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	CSRFSecret        string `json:"csrf_secret" form:"csrf_secret" schema:"csrf_secret"`
	InstantFormLeadID int64  `json:"instant_form_lead_id" form:"instant_form_lead_id" schema:"instant_form_lead_id"`
	InstantFormID     int64  `json:"instant_form_id" form:"instant_form_id" schema:"instant_form_id"`
	InstantFormName   string `json:"instant_form_name" form:"instant_form_name" schema:"instant_form_name"`
	ReferralLeadID    int    `json:"referral_lead_id" form:"referral_lead_id" schema:"referral_lead_id"`
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
	Text        string `json:"text"`
	DateCreated int64  `json:"date_created"`
	TextFrom    string `json:"text_from"`
	TextTo      string `json:"text_to"`
	IsInbound   bool   `json:"is_inbound"`
	Status      string `json:"status" form:"status" schema:"status"`
	IsRead      bool   `json:"is_read"`
}

type PhoneCall struct {
	PhoneCallID  int    `json:"phone_call_id" form:"phone_call_id" schema:"phone_call_id"`
	ExternalID   string `json:"external_id" form:"external_id" schema:"external_id"`
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
	UserRoleID int    `json:"user_role_id" form:"user_role_id" schema:"user_role_id"`
	Role       string `json:"role" form:"role" schema:"role"`
}

type Event struct {
	EventID       int     `json:"event_id" form:"event_id" schema:"event_id"`
	BartenderID   int     `json:"bartender_id" form:"bartender_id" schema:"bartender_id"`
	LeadID        int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	StreetAddress string  `json:"street_address" form:"street_address" schema:"street_address"`
	City          string  `json:"city" form:"city" schema:"city"`
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
	CocktailID int    `json:"cocktail_id" form:"cocktail_id" schema:"cocktail_id"`
	Name       string `json:"name" form:"name" schema:"name"`
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

type Quote struct {
	QuoteID            int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	ExternalID         string  `json:"external_id" form:"external_id" schema:"external_id"`
	LeadID             int     `json:"lead_id" form:"lead_id" schema:"lead_id"`
	EventTypeID        int     `json:"event_type_id" form:"event_type_id" schema:"event_type_id"`
	VenueTypeID        int     `json:"venue_type_id" form:"venue_type_id" schema:"venue_type_id"`
	Guests             int     `json:"guests" form:"guests" schema:"guests"`
	Hours              float64 `json:"hours" form:"hours" schema:"hours"`
	NumberOfBartenders int     `json:"number_of_bartenders" form:"number_of_bartenders" schema:"number_of_bartenders"`
	EventDate          int64   `json:"event_date" form:"event_date" schema:"event_date"`
}

type InvoiceType struct {
	InvoiceTypeID    int     `json:"invoice_type_id" form:"invoice_type_id" schema:"invoice_type_id"`
	Type             string  `json:"type" form:"type" schema:"type"`
	AmountPercentage float64 `json:"amount_percentage" form:"amount_percentage" schema:"amount_percentage"`
}

type Invoice struct {
	InvoiceID       int    `json:"invoice_id" form:"invoice_id" schema:"invoice_id"`
	QuoteID         int    `json:"quote_id" form:"quote_id" schema:"quote_id"`
	DateCreated     int64  `json:"date_created" form:"date_created" schema:"date_created"`
	DatePaid        int64  `json:"date_paid" form:"date_paid" schema:"date_paid"`
	DueDate         int64  `json:"due_date" form:"due_date" schema:"due_date"`
	InvoiceTypeID   int    `json:"invoice_type_id" form:"invoice_type_id" schema:"invoice_type_id"`
	URL             string `json:"url" form:"url" schema:"url"`
	StripeInvoiceID string `json:"stripe_invoice_id" form:"stripe_invoice_id" schema:"stripe_invoice_id"`
}

type LeadStatus struct {
	LeadStatusID int    `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	Status       string `json:"status" form:"status" schema:"status"`
}

type LeadInterest struct {
	LeadInterestID int    `json:"lead_interest_id" form:"lead_interest_id" schema:"lead_interest_id"`
	Interest       string `json:"interest" form:"interest" schema:"interest"`
}

type NextAction struct {
	NextActionID int    `json:"next_action_id" form:"next_action_id" schema:"next_action_id"`
	Action       string `json:"action" form:"action" schema:"action"`
}

type LeadNextAction struct {
	LeadNextActionID int   `json:"lead_next_action_id" form:"lead_next_action_id" schema:"lead_next_action_id"`
	NextActionID     int   `json:"next_action_id" form:"next_action_id" schema:"next_action_id"`
	LeadID           int   `json:"lead_id" form:"lead_id" schema:"lead_id"`
	ActionDate       int64 `json:"action_date" form:"action_date" schema:"action_date"`
}

type Service struct {
	ServiceID      int     `json:"service_id" form:"service_id" schema:"service_id"`
	ServiceTypeID  int     `json:"service_type_id" form:"service_type_id" schema:"service_type_id"`
	Service        string  `json:"service" form:"service" schema:"service"`
	SuggestedPrice float64 `json:"suggested_price" form:"suggested_price" schema:"suggested_price"`
	GuestRatio     int     `json:"guest_ratio" form:"guest_ratio" schema:"guest_ratio"`
}

type ServiceType struct {
	ServiceTypeID int    `json:"service_type_id" form:"service_type_id" schema:"service_type_id"`
	Type          string `json:"type" form:"type" schema:"type"`
}

type QuoteService struct {
	QuoteServiceID int     `json:"quote_service_id" form:"quote_service_id" schema:"quote_service_id"`
	ServiceID      int     `json:"service_id" form:"service_id" schema:"service_id"`
	QuoteID        int     `json:"quote_id" form:"quote_id" schema:"quote_id"`
	Units          float64 `json:"units" form:"units" schema:"units"`
	PricePerUnit   float64 `json:"price_per_unit" form:"price_per_unit" schema:"price_per_unit"`
}

type LeadNote struct {
	LeadNoteID    int    `json:"lead_note_id" form:"lead_note_id" schema:"lead_note_id"`
	Note          string `json:"note" form:"note" schema:"note"`
	LeadID        int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	DateAdded     int64  `json:"date_added" form:"date_added" schema:"date_added"`
	AddedByUserID int    `json:"added_by_user_id" form:"added_by_user_id" schema:"added_by_user_id"`
}

type PhoneCallTranscription struct {
	PhoneCallTranscriptionID int    `json:"phone_call_transcription_id" form:"phone_call_transcription_id" schema:"phone_call_transcription_id"`
	PhoneCallID              int    `json:"phone_call_id" form:"phone_call_id" schema:"phone_call_id"`
	ExternalID               string `json:"external_id" form:"external_id" schema:"external_id"`
	Text                     string `json:"text" form:"text" schema:"text"`
	AudioURL                 string `json:"audio_url" form:"audio_url" schema:"audio_url"`
	TextURL                  string `json:"text_url" form:"text_url" schema:"text_url"`
}

type EventRole struct {
	EventRoleID int    `json:"event_role_id" form:"event_role_id" schema:"event_role_id"`
	Role        string `json:"role" form:"role" schema:"role"`
}

type EventStaff struct {
	EventStaffID int     `json:"event_staff_id" form:"event_staff_id" schema:"event_staff_id"`
	UserID       int     `json:"user_id" form:"user_id" schema:"user_id"`
	EventID      int     `json:"event_id" form:"event_id" schema:"event_id"`
	EventRoleID  int     `json:"event_role_id" form:"event_role_id" schema:"event_role_id"`
	HourlyRate   float64 `json:"hourly_rate" form:"hourly_rate" schema:"hourly_rate"`
}
