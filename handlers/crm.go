package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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
			if len(parts) >= 5 && parts[4] == "event" && helpers.IsNumeric(parts[3]) {
				GetEventDetail(w, r, ctx)
				return
			}
			if len(parts) >= 5 && parts[4] == "quote" && helpers.IsNumeric(parts[3]) {
				GetLeadQuoteDetail(w, r, ctx)
				return
			}
			return
		}

		switch path {
		case "/crm/lead":
			GetLeads(w, r, ctx)
		case "/crm/service":
			GetServices(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPut:
		parts := strings.Split(path, "/")
		if strings.HasPrefix(path, "/crm/quote-service/") {
			PutQuoteService(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/service/") {
			PutService(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/lead/") {
			if len(path) > len("/crm/lead/") && strings.Contains(path, "archive") {
				ArchiveLead(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "marketing" && helpers.IsNumeric(parts[3]) {
				PutLeadMarketing(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "event" && helpers.IsNumeric(parts[3]) {
				PutEvent(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "quote" && helpers.IsNumeric(parts[3]) {
				PutLeadQuote(w, r)
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
		if strings.HasPrefix(path, "/crm/quote-service/") {
			DeleteQuoteService(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/lead/") {
			parts := strings.Split(path, "/")
			if len(parts) >= 5 && parts[4] == "event" && helpers.IsNumeric(parts[3]) {
				DeleteEvent(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/service/") {
			DeleteService(w, r)
			return
		}
	case http.MethodPost:
		if strings.HasPrefix(path, "/crm/quote-service") {
			PostQuoteService(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/lead/") {
			parts := strings.Split(path, "/")
			if strings.Contains(path, "invoice") {
				PostSendInvoice(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "event" && helpers.IsNumeric(parts[3]) {
				PostEvent(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "quote" && helpers.IsNumeric(parts[3]) {
				PostLeadQuote(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "note" && helpers.IsNumeric(parts[3]) {
				PostLeadNote(w, r)
				return
			}
		}
		switch path {
		case "/crm/service":
			PostService(w, r)
		case "/crm/quote-service":
			PostSendInvoice(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
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
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))
	params.Search = helpers.SafeStringToPointer(r.URL.Query().Get("search"))
	params.LeadInterestID = helpers.SafeStringToIntPointer(r.URL.Query().Get("lead_interest_id"))
	params.LeadStatusID = helpers.SafeStringToIntPointer(r.URL.Query().Get("lead_status_id"))
	params.NextActionID = helpers.SafeStringToIntPointer(r.URL.Query().Get("next_action_id"))

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

	interests, err := database.GetLeadInterestList()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	statuses, err := database.GetLeadStatusList()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	nextActions, err := database.GetNextActionList()
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
	data["Interests"] = interests
	data["Statuses"] = statuses
	data["NextActions"] = nextActions

	data["CurrentPage"] = 1
	if params.PageNum != nil {
		data["CurrentPage"] = *params.PageNum
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetLeadDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lead_detail.html"
	eventForm := constants.PARTIAL_TEMPLATES_DIR + "event_form.html"
	eventTable := constants.PARTIAL_TEMPLATES_DIR + "events_table.html"
	leadQuoteForm := constants.PARTIAL_TEMPLATES_DIR + "lead_quote_form.html"
	leadQuoteTable := constants.PARTIAL_TEMPLATES_DIR + "lead_quotes_table.html"
	createLeadNoteForm := constants.PARTIAL_TEMPLATES_DIR + "create_lead_note_form.html"
	leadNotesTemplate := constants.PARTIAL_TEMPLATES_DIR + "lead_notes.html"
	createLeadMessageForm := constants.PARTIAL_TEMPLATES_DIR + "create_lead_message_form.html"
	leadMessagesTemplate := constants.PARTIAL_TEMPLATES_DIR + "lead_messages.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, eventForm, eventTable, leadQuoteForm, leadQuoteTable, createLeadMessageForm, leadMessagesTemplate, createLeadNoteForm, leadNotesTemplate}
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
		http.Error(w, "Error getting venue types.", http.StatusInternalServerError)
		return
	}

	events, err := database.GetEventList(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting events.", http.StatusInternalServerError)
		return
	}

	leadQuotes, err := database.GetLeadQuotes(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead quotes.", http.StatusInternalServerError)
		return
	}

	bartenders, err := database.GetUsers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bartenders.", http.StatusInternalServerError)
		return
	}

	leadInterestList, err := database.GetLeadInterestList()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead interest list.", http.StatusInternalServerError)
		return
	}

	leadStatusList, err := database.GetLeadStatusList()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead status list.", http.StatusInternalServerError)
		return
	}

	nextActionList, err := database.GetNextActionList()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting next action list.", http.StatusInternalServerError)
		return
	}

	var params types.GetLeadsParams
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))

	referrals, err := database.GetReferrals()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting referrals.", http.StatusInternalServerError)
		return
	}

	leadMessages, err := database.GetMessagesByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead messages.", http.StatusInternalServerError)
		return
	}

	leadNotes, err := database.GetLeadNotesByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead notes.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Lead Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Lead"] = leadDetails
	data["CRMUserPhoneNumber"] = phoneNumber
	data["UserID"] = values.UserID
	data["EventTypes"] = eventTypes
	data["VenueTypes"] = venueTypes
	data["Events"] = events
	data["Bartenders"] = bartenders
	data["Referrals"] = referrals
	data["LeadQuotes"] = leadQuotes
	data["LeadInterestList"] = leadInterestList
	data["LeadStatusList"] = leadStatusList
	data["NextActionList"] = nextActionList
	data["LeadNotes"] = leadNotes
	data["LeadMessages"] = leadMessages

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

func ArchiveLead(w http.ResponseWriter, r *http.Request) {
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

	err = database.ArchiveLead(leadId)
	if err != nil {
		fmt.Printf("Error archiving lead: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to archive lead.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var params types.GetLeadsParams
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

func PostEvent(w http.ResponseWriter, r *http.Request) {
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

	var form types.EventForm
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

	err = database.CreateEvent(form)
	if err != nil {
		fmt.Printf("Error creating event: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating event.",
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

	if constants.Production {
		lead, err := database.GetConversionReporting(int(helpers.SafeInt(form.LeadID)))
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

		if lead.FacebookClickID != "" {
			fbEvent := types.FacebookEventData{
				EventName:      constants.EventConversionEventName,
				EventTime:      helpers.SafeInt64(form.DatePaid),
				ActionSource:   "phone_call",
				EventSourceURL: lead.LandingPage,
				UserData: types.FacebookUserData{
					Email:           helpers.HashString(lead.Email),
					Phone:           helpers.HashString(lead.PhoneNumber),
					FBC:             lead.FacebookClickID,
					FBP:             lead.FacebookClientID,
					ExternalID:      helpers.HashString(lead.ExternalID),
					ClientIPAddress: lead.IP,
					ClientUserAgent: lead.UserAgent,
				},
				CustomData: types.FacebookCustomData{
					Currency: constants.DefaultCurrency,
					Value:    fmt.Sprint(lead.Revenue),
				},
				EventID: fmt.Sprint(lead.EventID),
			}

			metaPayload := types.FacebookPayload{
				Data: []types.FacebookEventData{fbEvent},
			}

			go conversions.SendFacebookConversion(metaPayload)
		} else {
			fbLeadAdEvent := types.FacebookEventData{
				EventName:    constants.EventConversionEventName,
				EventTime:    helpers.SafeInt64(form.DatePaid),
				ActionSource: "phone_call",
				UserData: types.FacebookUserData{
					LeadID: lead.InstantFormLeadID,
				},
				CustomData: types.FacebookCustomData{
					Currency:        constants.DefaultCurrency,
					Value:           fmt.Sprint(lead.Revenue),
					EventSource:     constants.EventSourceCRM,
					LeadEventSource: constants.CompanyName,
				},
				EventID: fmt.Sprint(lead.EventID),
			}

			metaLeadAdPayload := types.FacebookPayload{
				Data: []types.FacebookEventData{fbLeadAdEvent},
			}

			go conversions.SendFacebookConversion(metaLeadAdPayload)
		}

		googlePayload := types.GooglePayload{
			ClientID: lead.GoogleClientID,
			UserId:   lead.ExternalID,
			Events: []types.GoogleEventLead{
				{
					Name: constants.EventConversionEventName,
					Params: types.GoogleEventParamsLead{
						GCLID:         lead.ClickID,
						TransactionID: fmt.Sprint(lead.EventID),
						Value:         lead.Revenue,
						Currency:      constants.DefaultCurrency,
						CampaignID:    fmt.Sprint(lead.CampaignID),
						Campaign:      lead.CampaignName,
					},
				},
			},
			UserData: types.GoogleUserData{
				Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
				Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},
			},
		}

		go conversions.SendGoogleConversion(googlePayload)
	}

	events, err := database.GetEventList(helpers.SafeInt(form.LeadID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting product slot assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "events_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "events_table.html",
		Data: map[string]any{
			"Events": events,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetEventDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "event_detail.html"
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

	eventId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting event id from path.", http.StatusInternalServerError)
		return
	}

	eventDetails, err := database.GetEventDetails(fmt.Sprint(eventId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting event details from DB.", http.StatusInternalServerError)
		return
	}

	bartenders, err := database.GetUsers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bartenders.", http.StatusInternalServerError)
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

	data := ctx
	data["PageTitle"] = "Event Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Event"] = eventDetails
	data["Bartenders"] = bartenders
	data["EventTypes"] = eventTypes
	data["VenueTypes"] = venueTypes

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutEvent(w http.ResponseWriter, r *http.Request) {
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

	var form types.EventForm
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

	err = database.UpdateEvent(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating event.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	if constants.Production {
		lead, err := database.GetConversionReporting(int(helpers.SafeInt(form.LeadID)))
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

		if lead.FacebookClickID != "" {
			fbEvent := types.FacebookEventData{
				EventName:      constants.EventConversionEventName,
				EventTime:      helpers.SafeInt64(form.DatePaid),
				ActionSource:   "phone_call",
				EventSourceURL: lead.LandingPage,
				UserData: types.FacebookUserData{
					Email:           helpers.HashString(lead.Email),
					Phone:           helpers.HashString(lead.PhoneNumber),
					FBC:             lead.FacebookClickID,
					FBP:             lead.FacebookClientID,
					ExternalID:      helpers.HashString(lead.ExternalID),
					ClientIPAddress: lead.IP,
					ClientUserAgent: lead.UserAgent,
				},
				CustomData: types.FacebookCustomData{
					Currency: constants.DefaultCurrency,
					Value:    fmt.Sprint(lead.Revenue),
				},
				EventID: fmt.Sprint(lead.EventID),
			}

			metaPayload := types.FacebookPayload{
				Data: []types.FacebookEventData{fbEvent},
			}

			go conversions.SendFacebookConversion(metaPayload)
		} else {
			fbLeadAdEvent := types.FacebookEventData{
				EventName:    constants.EventConversionEventName,
				EventTime:    helpers.SafeInt64(form.DatePaid),
				ActionSource: "phone_call",
				UserData: types.FacebookUserData{
					LeadID: lead.InstantFormLeadID,
				},
				CustomData: types.FacebookCustomData{
					Currency:        constants.DefaultCurrency,
					Value:           fmt.Sprint(lead.Revenue),
					EventSource:     constants.EventSourceCRM,
					LeadEventSource: constants.CompanyName,
				},
				EventID: fmt.Sprint(lead.EventID),
			}

			metaLeadAdPayload := types.FacebookPayload{
				Data: []types.FacebookEventData{fbLeadAdEvent},
			}

			go conversions.SendFacebookConversion(metaLeadAdPayload)
		}

		googlePayload := types.GooglePayload{
			ClientID: lead.GoogleClientID,
			UserId:   lead.ExternalID,
			Events: []types.GoogleEventLead{
				{
					Name: constants.EventConversionEventName,
					Params: types.GoogleEventParamsLead{
						GCLID:         lead.ClickID,
						TransactionID: fmt.Sprint(lead.EventID),
						Value:         lead.Revenue,
						Currency:      constants.DefaultCurrency,
						CampaignID:    fmt.Sprint(lead.CampaignID),
						Campaign:      lead.CampaignName,
					},
				},
			},
			UserData: types.GoogleUserData{
				Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
				Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},
			},
		}

		go conversions.SendGoogleConversion(googlePayload)
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Event has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	eventId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteEvent(eventId)
	if err != nil {
		fmt.Printf("Error deleting event: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete event.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	events, err := database.GetEventList(leadId)
	if err != nil {
		fmt.Printf("Error querying events after deletion: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to query events after deletion.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "events_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "events_table.html",
		Data: map[string]any{
			"Events": events,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLeadQuote(w http.ResponseWriter, r *http.Request) {
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

	var form types.LeadQuoteForm
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

	err = database.CreateLeadQuote(form)
	if err != nil {
		fmt.Printf("Error creating lead quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead quote.",
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

	leadQuotes, err := database.GetLeadQuotes(helpers.SafeInt(form.LeadID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting lead quotes from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_quotes_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_quotes_table.html",
		Data: map[string]any{
			"LeadQuotes": leadQuotes,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutLeadQuote(w http.ResponseWriter, r *http.Request) {
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

	var form types.LeadQuoteForm
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

	err = database.UpdateLeadQuote(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating quote.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// Check for event date validation
	formEventDate := helpers.SafeInt64(form.EventDate)
	if formEventDate == 0 {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Event date cannot be nil.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = services.UpdateInvoicesWorkflow(helpers.SafeInt(form.QuoteID), formEventDate)
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error during invoices workflow.",
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
			"AlertMessage": "Quote has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostSendInvoice(w http.ResponseWriter, r *http.Request) {
	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to parse lead id.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	quoteId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to parse quote id.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	quote, err := database.GetLeadQuoteInvoiceDetails(fmt.Sprint(leadId), fmt.Sprint(quoteId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting quote details.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// Do not create invoice if quote has invoice already
	if quote.InvoiceID == 0 {

		invoiceTypes, err := database.GetInvoiceTypes()
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error getting invoice types.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		stripeCustomerId := quote.StripeCustomerID

		for _, invoiceType := range invoiceTypes {
			invoiceDueDate := time.Now().Add(24 * time.Hour).Unix()

			// Remaining Invoice (0.75% due 48 hours prior to event)
			if invoiceType.InvoiceTypeID == constants.RemainingInvoiceTypeID {
				t := time.Unix(quote.EventDate, 0)
				invoiceDueDate = t.Add(-time.Duration(constants.InvoicePaymentDueInHours) * time.Hour).Unix()
			}

			createInvoiceParams := types.CreateInvoiceParams{
				Email:            quote.Email,
				StripeCustomerID: stripeCustomerId,
				FullName:         quote.FullName,
				DueDate:          invoiceDueDate,
				Quote:            quote.Amount * invoiceType.AmountPercentage,
			}

			createdInvoice, err := services.CreateStripeInvoice(createInvoiceParams)
			if err != nil {
				fmt.Printf("%+v\n", err)
				tmplCtx := types.DynamicPartialTemplate{
					TemplateName: "error",
					TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
					Data: map[string]any{
						"Message": "Error sending stripe invoice.",
					},
				}
				w.WriteHeader(http.StatusBadRequest)
				helpers.ServeDynamicPartialTemplate(w, tmplCtx)
				return
			}

			err = database.CreateQuoteInvoice(createdInvoice.ID, createdInvoice.HostedInvoiceURL, quoteId, invoiceType.InvoiceTypeID, invoiceDueDate)
			if err != nil {
				fmt.Printf("%+v\n", err)
				tmplCtx := types.DynamicPartialTemplate{
					TemplateName: "error",
					TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
					Data: map[string]any{
						"Message": "Error assigning stripe invoice ID to quote.",
					},
				}
				w.WriteHeader(http.StatusBadRequest)
				helpers.ServeDynamicPartialTemplate(w, tmplCtx)
				return
			}

			// This is only true when a new customer is created
			if stripeCustomerId == "" {
				stripeCustomerId = createdInvoice.Customer.ID
			}
		}

		err = database.AssignStripeCustomerIDToLead(stripeCustomerId, quote.LeadID)
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error assigning stripe customer ID to lead.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	var externalQuoteView = fmt.Sprintf("%s/external/%s", constants.RootDomain, quote.ExternalID)
	var textMessageTemplateNotification = fmt.Sprintf(
		`BARTENDING QUOTE:
		Here's the link to your estimate: %s
	`, externalQuoteView)

	_, err = services.SendTextMessage(quote.PhoneNumber, constants.CompanyPhoneNumber, textMessageTemplateNotification)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error sending invoice via text.",
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
			"AlertMessage": "Invoice has been sent.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetServices(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "services.html"
	createServiceForm := constants.CRM_TEMPLATES_DIR + "create_service_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "services_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createServiceForm}

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

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	services, totalRows, err := database.GetServicesList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting services from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Services — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Services"] = services
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostService(w http.ResponseWriter, r *http.Request) {
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

	var form types.ServiceForm
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

	err = database.CreateService(form)
	if err != nil {
		fmt.Printf("Error creating service: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create service.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	services, totalRows, err := database.GetServicesList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting services from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "services_table.html",
		Data: map[string]any{
			"Services":    services,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutService(w http.ResponseWriter, r *http.Request) {
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

	var form types.ServiceForm
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

	err = database.UpdateService(form)
	if err != nil {
		fmt.Printf("Error updating service: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update service.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	services, totalRows, err := database.GetServicesList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting services from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "services_table.html",
		Data: map[string]any{
			"Services":    services,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
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

	serviceId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/service/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteService(serviceId)
	if err != nil {
		fmt.Printf("Error deleting service: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete service.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	services, totalRows, err := database.GetServicesList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting services from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "services_table.html",
		Data: map[string]any{
			"Services":    services,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLeadQuoteDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lead_quote_detail.html"
	quoteServicesTable := constants.PARTIAL_TEMPLATES_DIR + "quote_services_table.html"
	createQuoteServiceForm := constants.CRM_TEMPLATES_DIR + "create_quote_service_form.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, quoteServicesTable, createQuoteServiceForm}
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

	quoteId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting quote id from URL.", http.StatusInternalServerError)
		return
	}

	quoteDetails, err := database.GetLeadQuoteDetails(fmt.Sprint(quoteId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting quote details from DB.", http.StatusInternalServerError)
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

	services, err := database.GetServices()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bar types.", http.StatusInternalServerError)
		return
	}

	quoteServices, err := database.GetQuoteServices(quoteId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting quote services.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Quote Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Quote"] = quoteDetails
	data["EventTypes"] = eventTypes
	data["VenueTypes"] = venueTypes
	data["QuoteServices"] = quoteServices
	data["Services"] = services

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func DeleteQuoteService(w http.ResponseWriter, r *http.Request) {
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

	quoteServiceId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/quote-service/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cannot get quote id after quote service has been deleted...
	quoteId, err := database.GetQuoteIDByQuoteServiceID(quoteServiceId)
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting quote id from quote service id.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.DeleteQuoteService(quoteServiceId)
	if err != nil {
		fmt.Printf("Error deleting quote service: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete quote service.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	hasInvoice, err := database.CheckQuoteHasInvoiceID(quoteId)
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error checking if quote has invoices.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	if hasInvoice {
		quote, err := database.GetLeadQuoteDetails(fmt.Sprint(quoteId))
		if err != nil {
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error during invoices workflow.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		err = services.UpdateInvoicesWorkflow(quote.QuoteID, quote.EventDate)
		if err != nil {
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error during invoices workflow.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	quoteServices, err := database.GetQuoteServices(quoteId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting services from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "quote_services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "quote_services_table.html",
		Data: map[string]any{
			"QuoteServices": quoteServices,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostQuoteService(w http.ResponseWriter, r *http.Request) {
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

	var form types.QuoteServiceForm
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

	err = database.CreateQuoteService(form)
	if err != nil {
		fmt.Printf("Error creating lead quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead quote.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	hasInvoice, err := database.CheckQuoteHasInvoiceID(helpers.SafeInt(form.QuoteID))
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error checking if quote has invoices.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	if hasInvoice {
		quote, err := database.GetLeadQuoteDetails(fmt.Sprint(helpers.SafeInt(form.QuoteID)))
		if err != nil {
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error during invoices workflow.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		err = services.UpdateInvoicesWorkflow(quote.QuoteID, quote.EventDate)
		if err != nil {
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error during invoices workflow.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	quoteServices, err := database.GetQuoteServices(helpers.SafeInt(form.QuoteID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting quote services from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "quote_services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "quote_services_table.html",
		Data: map[string]any{
			"QuoteServices": quoteServices,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutQuoteService(w http.ResponseWriter, r *http.Request) {
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

	var form types.QuoteServiceForm
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

	err = database.UpdateQuoteService(form)
	if err != nil {
		fmt.Printf("Error updating quote service: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while updating quote service.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	hasInvoice, err := database.CheckQuoteHasInvoiceID(helpers.SafeInt(form.QuoteID))
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error checking if quote has invoices.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	if hasInvoice {
		quote, err := database.GetLeadQuoteDetails(fmt.Sprint(helpers.SafeInt(form.QuoteID)))
		if err != nil {
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error during invoices workflow.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		err = services.UpdateInvoicesWorkflow(quote.QuoteID, quote.EventDate)
		if err != nil {
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error during invoices workflow.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	quoteServices, err := database.GetQuoteServices(helpers.SafeInt(form.QuoteID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting quote services from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "quote_services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "quote_services_table.html",
		Data: map[string]any{
			"QuoteServices": quoteServices,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLeadNote(w http.ResponseWriter, r *http.Request) {
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

	leadIdForm := r.FormValue("lead_id")
	userIdForm := r.FormValue("user_id")
	note := r.FormValue("note")

	leadID, err := strconv.Atoi(leadIdForm)
	if err != nil {
		fmt.Printf("Error converting lead_id to int: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid lead ID.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	userID, err := strconv.Atoi(userIdForm)
	if err != nil {
		fmt.Printf("Error converting user_id to int: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid user ID.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadNote := models.LeadNote{
		LeadID:        leadID,
		Note:          note,
		DateAdded:     time.Now().Unix(),
		AddedByUserID: userID,
	}

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

	err = database.CreateLeadNote(leadNote)
	if err != nil {
		fmt.Printf("Error creating lead quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead quote.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadNotes, err := database.GetLeadNotesByLeadID(leadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting lead quotes from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_notes.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_notes.html",
		Data: map[string]any{
			"LeadNotes": leadNotes,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
