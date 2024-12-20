package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/conversions"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/models"
	"github.com/davidalvarez305/yd_cocktails/services"
	"github.com/davidalvarez305/yd_cocktails/sessions"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/davidalvarez305/yd_cocktails/utils"
	"github.com/gorilla/schema"
)

const (
	YovaHeroImage string = "https://ydcocktails.s3.us-east-1.amazonaws.com/media/yova_hero.jpeg"
	YovaMidCTA    string = "https://ydcocktails.s3.us-east-1.amazonaws.com/media/yova_mid_cta.png"
)

var decoder = schema.NewDecoder()

var websiteBaseFilePath = constants.WEBSITE_TEMPLATES_DIR + "base.html"
var websiteFooterFilePath = constants.WEBSITE_TEMPLATES_DIR + "footer.html"

func createWebsiteContext() types.WebsiteContext {
	return types.WebsiteContext{
		PageTitle:                    constants.CompanyName,
		MetaDescription:              "Get a quote for mobile bartending services in Miami, FL.",
		SiteName:                     constants.SiteName,
		StaticPath:                   constants.StaticPath,
		MediaPath:                    constants.MediaPath,
		PhoneNumber:                  helpers.FormatPhoneNumber(constants.CompanyPhoneNumber),
		CurrentYear:                  time.Now().Year(),
		GoogleAnalyticsID:            constants.GoogleAnalyticsID,
		GoogleAdsID:                  constants.GoogleAdsID,
		GoogleAdsCallConversionLabel: constants.GoogleAdsCallConversionLabel,
		FacebookDataSetID:            constants.FacebookDatasetID,
		CompanyName:                  constants.CompanyName,
		LeadEventName:                constants.LeadEventName,
		LeadGeneratedEventName:       constants.LeadGeneratedEventName,
		DefaultCurrency:              constants.DefaultCurrency,
		YovaHeroImage:                YovaHeroImage,
		YovaMidCTA:                   YovaMidCTA,
	}
}

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createWebsiteContext()
	ctx.PagePath = constants.RootDomain + r.URL.Path

	externalId, ok := r.Context().Value("external_id").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in context.", http.StatusInternalServerError)
		return
	}

	ctx.ExternalID = externalId

	switch r.Method {
	case http.MethodGet:

		path := r.URL.Path

		if strings.HasPrefix(path, "/quote/") {
			if len(path) > len("/quote/") && helpers.IsNumeric(path[len("/quote/"):]) {
				GetQuoteDetail(w, r, ctx)
				return
			}
			return
		}
		switch r.URL.Path {
		case "/contact":
			GetContactForm(w, r, ctx)
		case "/login":
			GetLogin(w, r, ctx)
		case "/privacy-policy":
			GetPrivacyPolicy(w, r, ctx)
		case "/terms-and-conditions":
			GetTermsAndConditions(w, r, ctx)
		case "/robots.txt":
			GetRobots(w, r, ctx)
		case "/":
			GetHome(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/quote":
			PostQuote(w, r)
		case "/estimate":
			PostEstimate(w, r)
		case "/contact":
			PostContactForm(w, r)
		case "/login":
			PostLogin(w, r)
		case "/logout":
			PostLogout(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetHome(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	isMobile := helpers.IsMobileRequest(r)
	heroImagePath := "hero_image_desktop.html"
	if isMobile {
		heroImagePath = "hero_image_mobile.html"
	}

	fileName := "home.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + heroImagePath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}
	eventTypes, err := database.GetEventTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting event types.", http.StatusInternalServerError)
		return
	}

	venueTypes, err := database.GetVenueTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting venue types.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Miami Mobile Bartending Services — " + constants.CompanyName
	data.Nonce = nonce
	data.Features = []string{
		"We'll work with you to create a custom menu that features our signature cocktails + your favorites.",
		"We'll always be early to setup & make sure everything that's necessary is ready for use.",
		"We have high standards of service to make sure your guests are able to enjoy their time with cold & delicious drinks.",
		"We will DEFINITELY clean up after ourselves and leave your area as clean as it was before we got there.",
		"Our team can dress to the occasion in the event that you require a specific outfit or a certain theme.",
		"We're very flexible in terms of capacity for number of attendees, and can serve small as well as larger events.",
		"We offer highly detailed & customized quotes so that you know 100% what you're paying for, and what we agree to.",
		"Your guests are first, and it's our priority to put forth an incredible service so that their experience at your event is awesome.",
		"Our bartenders are highly skilled with years of experience so your cocktails come out delish.",
	}
	data.CSRFToken = csrfToken
	data.VenueTypes = venueTypes
	data.EventTypes = eventTypes

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetRobots(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	robotsTxtContent := `
	# robots.txt for https://ydcocktails.com/

	# Allow all robots complete access
	User-agent: *
	Disallow:
	`

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err := w.Write([]byte(robotsTxtContent))
	if err != nil {
		http.Error(w, "Error writing robots.txt content.", http.StatusInternalServerError)
	}
}

func GetPrivacyPolicy(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "privacy.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName, quoteForm}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Privacy Policy — " + constants.CompanyName
	data.Nonce = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetTermsAndConditions(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "terms.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName, quoteForm}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Terms & Conditions — " + constants.CompanyName
	data.Nonce = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostQuote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.QuoteForm
	form.FirstName = helpers.GetStringPointerFromForm(r, "first_name")
	form.LastName = helpers.GetStringPointerFromForm(r, "last_name")
	form.PhoneNumber = helpers.GetStringPointerFromForm(r, "phone_number")
	form.EventType = helpers.GetIntPointerFromForm(r, "event_type")
	form.VenueType = helpers.GetIntPointerFromForm(r, "venue_type")
	form.Message = helpers.GetStringPointerFromForm(r, "message")
	form.Source = helpers.GetStringPointerFromForm(r, "source")
	form.Medium = helpers.GetStringPointerFromForm(r, "medium")
	form.Channel = helpers.GetStringPointerFromForm(r, "channel")
	form.LandingPage = helpers.GetStringPointerFromForm(r, "landing_page")
	form.Keyword = helpers.GetStringPointerFromForm(r, "keyword")
	form.Referrer = helpers.GetStringPointerFromForm(r, "referrer")
	form.ClickID = helpers.GetStringPointerFromForm(r, "click_id")
	form.CampaignID = helpers.GetInt64PointerFromForm(r, "campaign_id")
	form.AdCampaign = helpers.GetStringPointerFromForm(r, "ad_campaign")
	form.AdGroupID = helpers.GetInt64PointerFromForm(r, "ad_group_id")
	form.AdGroupName = helpers.GetStringPointerFromForm(r, "ad_group_name")
	form.AdSetID = helpers.GetInt64PointerFromForm(r, "ad_set_id")
	form.AdSetName = helpers.GetStringPointerFromForm(r, "ad_set_name")
	form.AdID = helpers.GetInt64PointerFromForm(r, "ad_id")
	form.AdHeadline = helpers.GetInt64PointerFromForm(r, "ad_headline")
	form.Language = helpers.GetStringPointerFromForm(r, "language")
	form.Longitude = helpers.GetStringPointerFromForm(r, "longitude")
	form.Latitude = helpers.GetStringPointerFromForm(r, "latitude")
	form.UserAgent = helpers.GetStringPointerFromForm(r, "user_agent")
	form.ButtonClicked = helpers.GetStringPointerFromForm(r, "button_clicked")
	form.IP = helpers.GetStringPointerFromForm(r, "ip")
	form.Email = helpers.GetStringPointerFromForm(r, "email")
	form.OptInTextMessaging = helpers.GetBoolPointerFromForm(r, "opt_in_text_messaging")

	form.Estimate = helpers.GetFloat64PointerFromForm(r, "estimate")
	form.PackageID = helpers.GetIntPointerFromForm(r, "package_id")

	form.FacebookClickID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbc", "facebook_click_id")
	form.FacebookClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbp", "facebook_client_id")
	form.GoogleClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_ga", "google_client_id")

	session, err := sessions.Get(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to retrieve session.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// User Marketing Variables
	var userIP = helpers.GetUserIPFromRequest(r)
	var userAgent = r.Header.Get("User-Agent")

	if userIP != "" {
		form.IP = &userIP
	}

	if userAgent != "" {
		form.UserAgent = &userAgent
	}

	if session.ExternalID != "" {
		form.ExternalID = &session.ExternalID
	}

	if session.CSRFSecret != "" {
		form.CSRFSecret = &session.CSRFSecret
	}

	leadID, err := database.CreateLeadAndMarketing(form)
	if err != nil {
		fmt.Printf("Error creating lead: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating quote request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	fbEvent := types.FacebookEventData{
		EventName:      constants.LeadEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: helpers.SafeString(form.LandingPage),
		UserData: types.FacebookUserData{
			FirstName:       helpers.HashString(helpers.SafeString(form.FirstName)),
			LastName:        helpers.HashString(helpers.SafeString(form.LastName)),
			Phone:           helpers.HashString(helpers.SafeString(form.PhoneNumber)),
			Email:           helpers.HashString(helpers.SafeString(form.Email)),
			FBC:             helpers.SafeString(form.FacebookClickID),
			FBP:             helpers.SafeString(form.FacebookClientID),
			State:           helpers.HashString("Florida"),
			ExternalID:      helpers.HashString(helpers.SafeString(form.ExternalID)),
			ClientIPAddress: helpers.SafeString(form.IP),
			ClientUserAgent: helpers.SafeString(form.UserAgent),
		},
		CustomData: types.FacebookCustomData{
			Value:    fmt.Sprint(helpers.SafeFloat64(form.Estimate)),
			Currency: constants.DefaultCurrency,
		},
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	payload := types.GooglePayload{
		ClientID: helpers.SafeString(form.GoogleClientID),
		UserId:   helpers.SafeString(form.ExternalID),
		Events: []types.GoogleEventLead{
			{
				Name: constants.LeadGeneratedEventName,
				Params: types.GoogleEventParamsLead{
					GCLID:      helpers.SafeString(form.ClickID),
					Value:      helpers.SafeFloat64(form.Estimate),
					Currency:   constants.DefaultCurrency,
					CampaignID: fmt.Sprint(helpers.SafeInt64(form.CampaignID)),
					Campaign:   helpers.SafeString(form.AdCampaign),
					Source:     helpers.SafeString(form.Source),
					Medium:     helpers.SafeString(form.Medium),
					Term:       helpers.SafeString(form.Keyword),
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256PhoneNumber:  []string{helpers.HashString(utils.AddPhonePrefixIfNeeded(helpers.SafeString(form.PhoneNumber)))},
			Sha256EmailAddress: []string{helpers.HashString(helpers.SafeString(form.Email))},
			Address: []types.GoogleUserAddress{
				{
					Sha256FirstName: helpers.HashString(helpers.SafeString(form.FirstName)),
					Sha256LastName:  helpers.HashString(helpers.SafeString(form.LastName)),
				},
			},
		},
	}

	// Send conversion events
	go conversions.SendGoogleConversion(payload)
	go conversions.SendFacebookConversion(metaPayload)

	lead, err := database.GetConversionLeadInfo(leadID)
	if err != nil {
		fmt.Printf("Error getting conversion: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Internal error reporting conversions to Google.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	go func() {
		subject := "YD Cocktails: New Lead"
		recipients := []string{constants.CompanyEmail}
		templateFile := constants.PARTIAL_TEMPLATES_DIR + "new_lead_notification_email.html"

		var notificationTemplateData = map[string]any{
			"Name":           helpers.SafeString(form.FirstName) + " " + helpers.SafeString(form.LastName),
			"PhoneNumber":    helpers.SafeString(form.PhoneNumber),
			"DateCreated":    utils.FormatTimestampEST(lead.CreatedAt),
			"EventType":      lead.EventType,
			"VenueType":      lead.VenueType,
			"Message":        helpers.SafeString(form.Message),
			"LeadDetailsURL": fmt.Sprintf("%s/crm/lead/%d", constants.RootDomain, leadID),
			"Location":       "",
		}

		if helpers.SafeString(form.Longitude) != "0.0" && len(helpers.SafeString(form.Longitude)) > 0 || helpers.SafeString(form.Latitude) != "0.0" && len(helpers.SafeString(form.Latitude)) > 0 {
			notificationTemplateData["Location"] = fmt.Sprintf("https://www.google.com/maps?q=%s,%s", helpers.SafeString(form.Latitude), helpers.SafeString(form.Longitude))
		}

		template, err := helpers.BuildStringFromTemplate(templateFile, "email", notificationTemplateData)

		if err != nil {
			fmt.Printf("ERROR BUILDING LEAD NOTIFICATION TEMPLATE: %+v\n", err)
			return
		}

		body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
		err = services.SendGmail(recipients, subject, constants.CompanyEmail, body)
		if err != nil {
			fmt.Printf("ERROR SENDING LEAD NOTIFICATION EMAIL: %+v\n", err)
			return
		}
	}()

	err = database.AssignLeadToPackage(helpers.SafeInt(form.PackageID), lead.LeadID)
	if err != nil {
		fmt.Printf("Error creating invoice: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Internal server error while adding package to your account.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	stripeInvoice, err := services.CreateStripeInvoiceForNewCustomer(helpers.SafeString(form.Email), helpers.SafeString(form.FirstName), helpers.SafeString(form.LastName), helpers.SafeFloat64(form.Estimate))
	if err != nil {
		fmt.Printf("Error creating invoice: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating invoice.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var redirectUrl = fmt.Sprintf("%s/quote/%s", constants.RootDomain, stripeInvoice.ID)

	go func() {
		to := "+1" + helpers.RemoveCountryCode(helpers.SafeString(form.PhoneNumber))
		from := "+1" + constants.CompanyPhoneNumber

		textMessage := fmt.Sprintf(`Hey! This is YD Cocktails, LLC.
		Your quote has been successfully generated. If you'd like to lock down a spot for your event, you can do so here: %s`, redirectUrl)

		twilioMessage, err := services.SendTextMessage(to, from, textMessage)
		if err != nil {
			fmt.Printf("Error creating invoice: %+v\n", err)
		}

		message := models.Message{
			ExternalID:  helpers.SafeString(twilioMessage.Sid),
			UserID:      constants.DavidUserID,
			LeadID:      lead.LeadID,
			Text:        helpers.SafeString(twilioMessage.Body),
			TextFrom:    constants.CompanyPhoneNumber,
			TextTo:      helpers.SafeString(form.PhoneNumber),
			IsInbound:   false,
			DateCreated: time.Now().Local().Unix(),
			Status:      helpers.SafeString(twilioMessage.Status),
		}

		if err := database.SaveSMS(message); err != nil {
			fmt.Printf("Error saving SMS to database: %s", err)
			http.Error(w, "Failed to save message.", http.StatusInternalServerError)
			return
		}
	}()

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func GetContactForm(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "contact_form.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName, quoteForm}

	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Contact Us — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostContactForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to parse form data.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.ContactForm
	err = decoder.Decode(&form, r.PostForm)

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error decoding form data.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	subject := "Contact Form: YD Cocktails"
	recipients := []string{constants.CompanyEmail}
	templateFile := constants.PARTIAL_TEMPLATES_DIR + "contact_form_email.html"

	template, err := helpers.BuildStringFromTemplate(templateFile, "email", form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error building e-mail template.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
	err = services.SendGmail(recipients, subject, form.Email, body)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to send message.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Sent!",
			"AlertMessage": "We've received your message and will be quick to respond.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLogin(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "login.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName, quoteForm}

	csrfSecret, ok := r.Context().Value("csrf_secret").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in login page.", http.StatusInternalServerError)
		return
	}

	session, err := database.GetSession(csrfSecret)
	if err != nil {
		http.Error(w, "Error trying to get session in login page.", http.StatusInternalServerError)
		return
	}

	if session.UserID > 0 {
		user, err := database.GetUserById(session.UserID)
		if err != nil {
			http.Error(w, "Error trying to get existing user from DB.", http.StatusInternalServerError)
			return
		}

		if user.UserRoleID == constants.UserAdminRoleID {
			http.Redirect(w, r, "/crm/dashboard", http.StatusSeeOther)
			return
		}
	}

	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Login — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "error",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
		Data:         map[string]any{},
	}

	user, err := database.GetUserByUsername(username)
	if err != nil {
		tmplCtx.Data["Message"] = "Invalid username."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	isValid := helpers.ValidatePassword(password, user.Password)
	if !isValid {
		tmplCtx.Data["Message"] = "Invalid password."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		tmplCtx.Data["Message"] = "Could not retrieve session."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	session.UserID = user.UserID
	err = sessions.Update(session)
	if err != nil {
		tmplCtx.Data["Message"] = "Could not update session."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	sessions.SetCookie(w, utils.GetSessionExpirationTime(), session.CSRFSecret)

	w.WriteHeader(http.StatusOK)
}

func PostLogout(w http.ResponseWriter, r *http.Request) {

	sessions.SetCookie(w, time.Now().Add(-1*time.Hour), "")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PostEstimate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.EstimateForm
	form.Guests = helpers.GetIntPointerFromForm(r, "guests")
	form.Hours = helpers.GetIntPointerFromForm(r, "hours")
	form.PackageTypeID = helpers.GetIntPointerFromForm(r, "package_type_id")
	form.AlcoholSegmentID = helpers.GetIntPointerFromForm(r, "alcohol_segment_id")
	form.WillProvideLiquor = helpers.GetBoolPointerFromForm(r, "will_provide_liquor")
	form.WillProvideBeerAndWine = helpers.GetBoolPointerFromForm(r, "will_provide_beer_and_wine")
	form.WillProvideMixers = helpers.GetBoolPointerFromForm(r, "will_provide_mixers")
	form.WillProvideJuices = helpers.GetBoolPointerFromForm(r, "will_provide_juices")
	form.WillProvideSoftDrinks = helpers.GetBoolPointerFromForm(r, "will_provide_soft_drinks")
	form.WillProvideCups = helpers.GetBoolPointerFromForm(r, "will_provide_cups")
	form.WillRequireGlassware = helpers.GetBoolPointerFromForm(r, "will_require_glassware")
	form.WillRequireMobileBar = helpers.GetBoolPointerFromForm(r, "will_require_mobile_bar")
	form.NumBars = helpers.GetIntPointerFromForm(r, "num_bars")

	packagePrice := helpers.CalculatePackagePrice(form)

	packageId, err := database.GeneratePackageEstimate(form, packagePrice)
	if err != nil {
		fmt.Printf("Error creating lead: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating quote request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	response := types.GenerateEstimateResponse{
		Data: types.EstimateData{
			PackageID: packageId,
			Estimate:  packagePrice,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while marshaling JSON.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	w.Write(jsonResponse)
}

func GetQuoteDetail(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "quote_detail.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	invoiceId := strings.TrimPrefix(r.URL.Path, "/crm/quote/")

	quoteDetails, err := database.GetQuoteDetailsByInvoiceID(invoiceId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting quote details from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = quoteDetails.FirstName + " " + quoteDetails.LastName + " Quote — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken
	data.Quote = quoteDetails

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
