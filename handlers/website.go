package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/conversions"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/services"
	"github.com/davidalvarez305/yd_cocktails/sessions"
	"github.com/davidalvarez305/yd_cocktails/types"
	"github.com/davidalvarez305/yd_cocktails/utils"
	"github.com/gorilla/schema"
)

const (
	YovaHeroImage          string = "https://ydcocktails.s3.us-east-1.amazonaws.com/media/yova_hero.jpeg"
	YovaMostPopularPackage string = "https://ydcocktails.s3.us-east-1.amazonaws.com/media/yova_mid_cta.png"
	YovaBasicPackage       string = "https://ydcocktails.s3.us-east-1.amazonaws.com/media/yova_basic_package.jpeg"
	YovaOpenBarPackage     string = "https://ydcocktails.s3.us-east-1.amazonaws.com/media/yova_open_bar_package.jpeg"
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
		DefaultLeadValue:             constants.DefaultLeadValue,
		YovaHeroImage:                YovaHeroImage,
		YovaMostPopularPackage:       YovaMostPopularPackage,
		YovaBasicPackage:             YovaBasicPackage,
		YovaOpenBarPackage:           YovaOpenBarPackage,
	}
}

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createWebsiteContext()
	ctx.PagePath = constants.RootDomain + r.URL.Path
	ctx.IsMobile = helpers.IsMobileRequest(r)

	externalId, ok := r.Context().Value("external_id").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in context.", http.StatusInternalServerError)
		return
	}

	ctx.ExternalID = externalId

	switch r.Method {
	case http.MethodGet:
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
		case "/planning":
			GetPlanningLP(w, r, ctx)
		case "/staffing":
			GetStaffingLP(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/quote":
			PostQuote(w, r)
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
	heroImagePath := "hero_image_desktop.html"
	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		heroImagePath = "hero_image_mobile.html"
		headerPath = "header_mobile.html"
	}

	fileName := "home.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, constants.WEBSITE_TEMPLATES_DIR + heroImagePath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetPlanningLP(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	heroImagePath := "planning_hero_image_desktop.html"
	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		heroImagePath = "planning_hero_image_mobile.html"
		headerPath = "header_mobile.html"
	}

	fileName := "planning_lp.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, constants.WEBSITE_TEMPLATES_DIR + heroImagePath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}

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
	data.PageTitle = "Miami Party Planning Service — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken
	data.Features = []string{
		"We collaborate with you to design a custom event plan, including themed décor, catering, and entertainment.",
		"Our team arrives early to set up and ensure everything is perfectly in place before your guests arrive.",
		"We uphold the highest standards of service, ensuring your guests enjoy a seamless and memorable experience.",
		"We will put together a list of highly qualified & vetted vendors to make your dreams a reality.",
		"Our staff can dress to match your event’s theme, ensuring a cohesive and stylish atmosphere.",
		"We cater to events of all sizes, from intimate gatherings to large celebrations, with ease and efficiency.",
		"Our detailed and transparent quotes ensure you know exactly what’s included, with no hidden fees.",
		"Your guests’ experience is our top priority—we go above and beyond to create an unforgettable event.",
		"Our team consists of experienced professionals who bring creativity and expertise to every event we plan.",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetStaffingLP(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	heroImagePath := "staffing_hero_image_desktop.html"
	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		heroImagePath = "staffing_hero_image_mobile.html"
		headerPath = "header_mobile.html"
	}

	fileName := "staffing_lp.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, constants.WEBSITE_TEMPLATES_DIR + heroImagePath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}

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
	data.PageTitle = "Miami Event Staff for Hire — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken
	data.Features = []string{
		"We work with you to bring your event vision to life—whether it’s the vibe, the crew, or the full setup.",
		"Our team shows up early to set up, handle the details, and make sure everything runs smoothly before guests arrive.",
		"Top-tier service is our thing—we make sure your guests feel taken care of and have an amazing time.",
		"We’ve got a network of trusted vendors to hook you up with everything from food to décor to entertainment.",
		"Our staff can dress to match your event theme, so everything looks seamless and on point.",
		"Big or small, we’ve got the crew to handle it all—intimate gatherings, huge parties, and everything in between.",
		"No surprises here—our pricing is clear and upfront, so you know exactly what you’re getting.",
		"Your guests’ experience is our top priority—we go the extra mile to make your event unforgettable.",
		"Our team is packed with experienced pros who bring creativity, energy, and expertise to every event.",
	}

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

	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		headerPath = "header_mobile.html"
	}

	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}
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

	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		headerPath = "header_mobile.html"
	}

	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}
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

	// Duplicate phone number logic
	phoneNumber := helpers.GetStringPointerFromForm(r, "phone_number")
	cleanedPhoneNumber := helpers.ExtractPhoneNumber(helpers.SafeString(phoneNumber))

	if cleanedPhoneNumber == "" {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Phone number field cannot be empty.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	phoneNumberExists, err := database.IsPhoneNumberInDB(cleanedPhoneNumber)
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error checking if phone number exists in database.",
			},
		}

		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	if phoneNumberExists {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "We already have a request with this phone number.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.QuoteForm
	form.FullName = helpers.GetStringPointerFromForm(r, "full_name")
	form.PhoneNumber = &cleanedPhoneNumber
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

	form.FacebookClickID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbc", "facebook_click_id")
	form.FacebookClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbp", "facebook_client_id")
	form.GoogleClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_ga", "google_client_id")

	var createdAt = time.Now().Unix()
	form.CreatedAt = &createdAt

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

	if constants.Production {

		fbEvent := types.FacebookEventData{
			EventName:      constants.LeadEventName,
			EventTime:      time.Now().UTC().Unix(),
			ActionSource:   "website",
			EventSourceURL: helpers.SafeString(form.LandingPage),
			UserData: types.FacebookUserData{
				Phone:           helpers.HashString(helpers.SafeString(form.PhoneNumber)),
				FBC:             helpers.SafeString(form.FacebookClickID),
				FBP:             helpers.SafeString(form.FacebookClientID),
				ExternalID:      helpers.HashString(helpers.SafeString(form.ExternalID)),
				ClientIPAddress: helpers.SafeString(form.IP),
				ClientUserAgent: helpers.SafeString(form.UserAgent),
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
						CampaignID: fmt.Sprint(helpers.SafeInt64(form.CampaignID)),
						Campaign:   helpers.SafeString(form.AdCampaign),
						Source:     helpers.SafeString(form.Source),
						Medium:     helpers.SafeString(form.Medium),
						Term:       helpers.SafeString(form.Keyword),
						Value:      constants.DefaultLeadValue,
						Currency:   constants.DefaultCurrency,
					},
				},
			},
			UserData: types.GoogleUserData{
				Sha256PhoneNumber: []string{helpers.HashString(utils.AddPhonePrefixIfNeeded(helpers.SafeString(form.PhoneNumber)))},
			},
		}

		// Send conversion events
		go conversions.SendGoogleConversion(payload)
		go conversions.SendFacebookConversion(metaPayload)

		go func() {
			subject := "YD Cocktails: New Lead"
			recipients := []string{constants.CompanyEmail}
			templateFile := constants.PARTIAL_TEMPLATES_DIR + "new_lead_notification_email.html"

			var notificationTemplateData = map[string]any{
				"Name":           helpers.SafeString(form.FullName),
				"PhoneNumber":    helpers.SafeString(form.PhoneNumber),
				"DateCreated":    utils.FormatTimestampWithOptions(createdAt, &types.TimestampFormatOptions{TimeZone: constants.TimeZone}),
				"ButtonClicked":  helpers.SafeString(form.ButtonClicked),
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

			for _, phoneNumber := range constants.NotificationSubscribers {

				var textMessageTemplateNotification = fmt.Sprintf(
					`NEW LEAD:

				Phone: %s,
				Full Name: %s,
				Message: %s
			`, helpers.SafeString(form.PhoneNumber), helpers.SafeString(form.FullName), helpers.SafeString(form.Message))

				_, err := services.SendTextMessage(phoneNumber, constants.CompanyPhoneNumber, textMessageTemplateNotification)

				if err != nil {
					fmt.Printf("ERROR SENDING FB LEAD AD NOTIFICATION MSG: %+v\n", err)
				}
			}
		}()
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Awesome!",
			"AlertMessage": "We received your request and will be right with you.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetContactForm(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "contact_form.html"
	quoteForm := constants.WEBSITE_TEMPLATES_DIR + "quote_form.html"

	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		headerPath = "header_mobile.html"
	}

	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}

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

	headerPath := "header_desktop.html"
	if ctx.IsMobile {
		headerPath = "header_mobile.html"
	}

	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + headerPath, quoteForm, constants.WEBSITE_TEMPLATES_DIR + fileName}

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
