package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/conversions"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/services"
	"github.com/davidalvarez305/yd_cocktails/sessions"
	"github.com/davidalvarez305/yd_cocktails/types"
)

var crmBaseFilePath = constants.CRM_TEMPLATES_DIR + "base.html"
var crmFooterFilePath = constants.CRM_TEMPLATES_DIR + "footer.html"

func createCrmContext() map[string]any {
	return map[string]any{
		"PageTitle":       constants.CompanyName,
		"MetaDescription": "Get a quote for mobile bartending services in Miami, FL.",
		"SiteName":        constants.SiteName,
		"StaticPath":      constants.StaticPath,
		"MediaPath":       constants.MediaPath,
		"PhoneNumber":     constants.DavidPhoneNumber,
		"CurrentYear":     time.Now().Year(),
		"CompanyName":     constants.CompanyName,
	}
}

func CRMHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createCrmContext()
	ctx["PagePath"] = constants.RootDomain + r.URL.Path
	path := r.URL.Path

	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(path, "/")
		if strings.HasPrefix(path, "/crm/lead/") {
			if len(path) > len("/crm/lead/") && helpers.IsNumeric(path[len("/crm/lead/"):]) {
				GetLeadDetail(w, r, ctx)
				return
			}
			if len(parts) >= 5 && parts[4] == "booking" && helpers.IsNumeric(parts[3]) {
				GetBookingDetail(w, r, ctx)
				return
			}
			if len(parts) >= 5 && parts[4] == "estimate" && helpers.IsNumeric(parts[3]) {
				GetEstimateDetail(w, r, ctx)
				return
			}
			return
		}

		switch path {
		case "/crm/lead":
			GetLeads(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPut:
		parts := strings.Split(path, "/")
		if strings.HasPrefix(path, "/crm/lead/") {
			if len(parts) >= 5 && parts[4] == "marketing" && helpers.IsNumeric(parts[3]) {
				PutLeadMarketing(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "booking" && helpers.IsNumeric(parts[3]) {
				PutBooking(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "estimate" && helpers.IsNumeric(parts[3]) {
				PutEstimate(w, r)
				return
			}
			if len(path) > len("/crm/lead/") && helpers.IsNumeric(path[len("/crm/lead/"):]) {
				PutLead(w, r)
				return
			}
		}
		switch path {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if strings.HasPrefix(path, "/crm/lead/") {
			parts := strings.Split(path, "/")
			if len(path) > len("/crm/lead/") && helpers.IsNumeric(path[len("/crm/lead/"):]) {
				DeleteLead(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "booking" && helpers.IsNumeric(parts[3]) {
				DeleteBooking(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "estimate" && helpers.IsNumeric(parts[3]) {
				DeleteEstimate(w, r)
				return
			}
		}
	case http.MethodPost:
		if strings.HasPrefix(path, "/crm/lead/") {
			parts := strings.Split(path, "/")
			if len(parts) >= 5 && parts[4] == "booking" && helpers.IsNumeric(parts[3]) {
				PostBooking(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "estimate" && helpers.IsNumeric(parts[3]) {
				PostEstimate(w, r)
				return
			}
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetLeads(w http.ResponseWriter, r *http.Request, ctx map[string]interface{}) {
	baseFile := constants.CRM_TEMPLATES_DIR + "leads.html"
	leadsTable := constants.PARTIAL_TEMPLATES_DIR + "leads_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, leadsTable, baseFile}

	// Retrieve nonce from request context
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	// Retrieve CSRF token from request context
	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	var params types.GetLeadsParams
	params.EventType = helpers.SafeStringToPointer(r.URL.Query().Get("event_type"))
	params.VenueType = helpers.SafeStringToPointer(r.URL.Query().Get("venue_type"))
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))

	leads, totalRows, err := database.GetLeadList(params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting leads from DB.", http.StatusInternalServerError)
		return
	}

	eventTypes, err := database.GetEventTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending types.", http.StatusInternalServerError)
		return
	}

	venueTypes, err := database.GetVenueTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Leads — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Leads"] = leads
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["EventTypes"] = eventTypes
	data["VenueTypes"] = venueTypes

	data["CurrentPage"] = 1
	if params.PageNum != nil {
		data["CurrentPage"] = *params.PageNum
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetLeadDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lead_detail.html"
	bookingForm := constants.PARTIAL_TEMPLATES_DIR + "booking_form.html"
	bookingTable := constants.PARTIAL_TEMPLATES_DIR + "bookings_table.html"
	estimateForm := constants.PARTIAL_TEMPLATES_DIR + "estimate_form.html"
	estimateTable := constants.PARTIAL_TEMPLATES_DIR + "estimates_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, bookingForm, bookingTable, estimateForm, estimateTable}
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

	leadId := strings.TrimPrefix(r.URL.Path, "/crm/lead/")

	leadDetails, err := database.GetLeadDetails(leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead details from DB.", http.StatusInternalServerError)
		return
	}

	values, err := sessions.Get(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user ID from session.", http.StatusInternalServerError)
		return
	}

	phoneNumber, err := database.GetPhoneNumberFromUserID(values.UserID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting phone number from user ID.", http.StatusInternalServerError)
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
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	bookings, err := database.GetBookingList(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	estimates, err := database.GetEstimateList(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Lead Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Lead"] = leadDetails
	data["CRMUserPhoneNumber"] = phoneNumber
	data["EventTypes"] = eventTypes
	data["VenueTypes"] = venueTypes
	data["Bookings"] = bookings
	data["Estimates"] = estimates

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutLead(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GenerateTokenInHeader(w, r)
	if err != nil {
		fmt.Printf("Error generating token: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error generating new token. Reload page.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	w.Header().Set("X-Csrf-Token", token)

	err = r.ParseForm()
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

	var form types.UpdateLeadForm
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

	err = database.UpdateLead(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating lead.",
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
			"AlertHeader":  "Success!",
			"AlertMessage": "Lead has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutLeadMarketing(w http.ResponseWriter, r *http.Request) {
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

	var form types.UpdateLeadMarketingForm
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

	err = database.UpdateLeadMarketing(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating lead marketing.",
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
			"AlertHeader":  "Success!",
			"AlertMessage": "Lead marketing has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteLead(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteLead(leadId)
	if err != nil {
		fmt.Printf("Error deleting lead: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete lead.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var params types.GetLeadsParams
	params.EventType = helpers.SafeStringToPointer(r.URL.Query().Get("event_type"))
	params.VenueType = helpers.SafeStringToPointer(r.URL.Query().Get("venue_type"))
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))

	leads, totalRows, err := database.GetLeadList(params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting leads from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := "1"
	safePageNum := helpers.SafeString(params.PageNum)
	if safePageNum != "" {
		pageNum = safePageNum
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "leads_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "leads_table.html",
		Data: map[string]any{
			"Leads":       leads,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
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
	form.WillProvideIce = helpers.GetBoolPointerFromForm(r, "will_provide_ice")
	form.WillRequireGlassware = helpers.GetBoolPointerFromForm(r, "will_require_glassware")
	form.WillRequireBar = helpers.GetBoolPointerFromForm(r, "will_require_bar")
	form.NumBars = helpers.GetIntPointerFromForm(r, "num_bars")

	packagePrice := helpers.CalculatePackagePrice(form)

	err = database.CreateEstimate(form, packagePrice)
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

	if form.LeadID == nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Lead ID cannot be null.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	lead, err := database.GetLeadDetails(fmt.Sprint(helpers.SafeInt(form.LeadID)))
	if err != nil {
		fmt.Printf("Error querying lead details: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error querying lead details.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	stripeInvoice, err := services.CreateStripeInvoiceForNewCustomer(lead.Email, lead.FirstName, lead.LastName, packagePrice)
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

	err = database.AssignStripeCustomerToLead(stripeInvoice.Customer.ID, lead.LeadID)
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

	err = database.AssignStripeInvoiceToEstimate(stripeInvoice.ID, helpers.SafeInt(form.EstimateID))
	if err != nil {
		fmt.Printf("Error assigning invoice to pkg: %+v\n", err)
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

	fbEvent := types.FacebookEventData{
		EventName:      constants.EstimateEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: lead.LandingPage,
		UserData: types.FacebookUserData{
			Email:           helpers.HashString(lead.Email),
			FirstName:       helpers.HashString(lead.FirstName),
			LastName:        helpers.HashString(lead.LastName),
			Phone:           helpers.HashString(lead.PhoneNumber),
			FBC:             lead.FacebookClickID,
			FBP:             lead.FacebookClientID,
			ExternalID:      helpers.HashString(lead.ExternalID),
			ClientIPAddress: lead.IP,
			ClientUserAgent: lead.UserAgent,
		},
		CustomData: types.FacebookCustomData{
			Currency: "USD",
			Value:    fmt.Sprint(float64(stripeInvoice.AmountDue) / 100),
		},
		EventID: stripeInvoice.ID,
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	googlePayload := types.GooglePayload{
		ClientID: lead.GoogleClientID,
		UserId:   lead.ExternalID,
		Events: []types.GoogleEventLead{
			{
				Name: constants.EstimateEventName,
				Params: types.GoogleEventParamsLead{
					GCLID:         lead.ClickID,
					TransactionID: stripeInvoice.ID,
					Value:         float64(stripeInvoice.AmountDue) / 100,
					Currency:      constants.DefaultCurrency,
					CampaignID:    fmt.Sprint(lead.CampaignID),
					Campaign:      lead.CampaignName,
					Source:        lead.Source,
					Medium:        lead.Medium,
					Term:          lead.Keyword,
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
			Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},
			Address: []types.GoogleUserAddress{
				{
					Sha256FirstName: helpers.HashString(lead.FirstName),
					Sha256LastName:  helpers.HashString(lead.LastName),
				},
			},
		},
	}

	go conversions.SendGoogleConversion(googlePayload)
	go conversions.SendFacebookConversion(metaPayload)

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

func PostBooking(w http.ResponseWriter, r *http.Request) {
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

	var form types.BookingForm
	form.BookingID = helpers.GetIntPointerFromForm(r, "booking_id")
	form.EstimateID = helpers.GetIntPointerFromForm(r, "estimate_id")
	form.StreetAddress = helpers.GetStringPointerFromForm(r, "street_address")
	form.City = helpers.GetStringPointerFromForm(r, "city")
	form.State = helpers.GetStringPointerFromForm(r, "state")
	form.PostalCode = helpers.GetStringPointerFromForm(r, "postal_code")
	form.Country = helpers.GetStringPointerFromForm(r, "country")
	form.StartTime = helpers.GetInt64PointerFromForm(r, "start_time")
	form.EndTime = helpers.GetInt64PointerFromForm(r, "end_time")
	form.BartenderID = helpers.GetIntPointerFromForm(r, "bartender_id")
	form.LeadID = helpers.GetIntPointerFromForm(r, "lead_id")

	err = database.CreateBooking(form)
	if err != nil {
		fmt.Printf("Error creating booking: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating booking.",
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
			"AlertHeader":  "Success!",
			"AlertMessage": "Booking has been created.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetBookingDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "booking_detail.html"
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

	bookingId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting booking id from path.", http.StatusInternalServerError)
		return
	}

	bookingDetails, err := database.GetBookingDetails(fmt.Sprint(bookingId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting booking details from DB.", http.StatusInternalServerError)
		return
	}

	bartenders, err := database.GetUsers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users from DB.", http.StatusInternalServerError)
		return
	}

	estimates, err := database.GetEstimateList(bookingDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Booking Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Booking"] = bookingDetails
	data["Bartenders"] = bartenders
	data["Estimates"] = estimates

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetEstimateDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "estimate_detail.html"
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

	estimateId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting estimate id from path.", http.StatusInternalServerError)
		return
	}

	estimateDetails, err := database.GetEstimateDetails(fmt.Sprint(estimateId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting estimate details from DB.", http.StatusInternalServerError)
		return
	}

	packageTypes, err := database.GetPackageTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting package types from DB.", http.StatusInternalServerError)
		return
	}

	alcoholSegments, err := database.GetAlcoholSegments()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting alcohol segments from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Estimate Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Estimate"] = estimateDetails
	data["PackageTypes"] = packageTypes
	data["AlcoholSegments"] = alcoholSegments

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutBooking(w http.ResponseWriter, r *http.Request) {
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

	var form types.BookingForm
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

	err = database.UpdateBooking(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating booking.",
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
			"AlertHeader":  "Success!",
			"AlertMessage": "Booking has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutEstimate(w http.ResponseWriter, r *http.Request) {
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

	var form types.EstimateForm
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

	err = database.UpdateEstimate(form, helpers.CalculatePackagePrice(form))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating booking.",
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
			"AlertHeader":  "Success!",
			"AlertMessage": "Estimate has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookingId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteBooking(bookingId)
	if err != nil {
		fmt.Printf("Error deleting booking: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete booking.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	bookings, err := database.GetBookingList(leadId)
	if err != nil {
		fmt.Printf("Error querying bookings after deletion: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to query bookings after deletion.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "bookings_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "bookings_table.html",
		Data: map[string]any{
			"Bookings": bookings,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteEstimate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	estimateId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteEstimate(estimateId)
	if err != nil {
		fmt.Printf("Error deleting estimate: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete estimate.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	estimates, err := database.GetEstimateList(leadId)
	if err != nil {
		fmt.Printf("Error querying estimates after deletion: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to query estimates after deletion.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "estimates_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "estimates_table.html",
		Data: map[string]any{
			"Estimates": estimates,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
