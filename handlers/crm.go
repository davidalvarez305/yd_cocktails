package handlers

import (
	"encoding/json"
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

	unreadMessages, err := database.GetUnreadMessagesCount()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting unread messages from DB.", http.StatusInternalServerError)
		return
	}
	ctx["UnreadMessages"] = unreadMessages

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
	ctx["CRMUserPhoneNumber"] = phoneNumber

	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(path, "/")

		if strings.HasPrefix(path, "/crm/user/") {
			if len(path) > len("/crm/user/") && helpers.IsNumeric(path[len("/crm/user/"):]) {
				GetUserDetail(w, r, ctx)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/cocktail/") {
			if len(path) > len("/crm/cocktail/") && helpers.IsNumeric(path[len("/crm/cocktail/"):]) {
				GetCocktailDetail(w, r, ctx)
				return
			}
		}

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

		if strings.HasPrefix(path, "/crm/message/leads") {
			GetLeadsWithMessages(w, r, ctx)
			return
		}

		if strings.HasPrefix(path, "/crm/message") {
			if len(path) > len("/crm/message/") && helpers.IsNumeric(path[len("/crm/message/"):]) {
				GetMessagesByLeadID(w, r, ctx)
				return
			}
		}

		switch path {
		case "/crm/lead":
			GetLeads(w, r, ctx)
		case "/crm/user":
			GetUsers(w, r, ctx)
		case "/crm/cocktail":
			GetCocktails(w, r, ctx)
		case "/crm/service":
			GetServices(w, r, ctx)
		case "/crm/message":
			GetMessages(w, r, ctx)
		case "/crm/event":
			GetEvents(w, r, ctx)
		case "/crm/automated-follow-up":
			GetAutomatedFollowUpMessage(w, r)
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

		if strings.HasPrefix(path, "/crm/message/") {
			if len(parts) >= 5 && parts[4] == "read" && helpers.IsNumeric(parts[3]) {
				SetSMSToRead(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/user/") {
			if len(path) > len("/crm/user/") && helpers.IsNumeric(path[len("/crm/user/"):]) {
				PutUser(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/cocktail/") {
			if len(path) > len("/crm/cocktail/") && helpers.IsNumeric(path[len("/crm/cocktail/"):]) {
				PutCocktail(w, r)
				return
			}
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
		parts := strings.Split(path, "/")

		if strings.HasPrefix(path, "/crm/quote-service/") {
			DeleteQuoteService(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/user/") {
			if len(path) > len("/crm/user/") && helpers.IsNumeric(path[len("/crm/user/"):]) {
				DeleteUser(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/cocktail/") {
			if len(path) > len("/crm/cocktail/") && helpers.IsNumeric(path[len("/crm/cocktail/"):]) {
				DeleteCocktail(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/lead/") {
			if len(parts) >= 5 && parts[4] == "event" && helpers.IsNumeric(parts[3]) {
				DeleteEvent(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "next-action" && helpers.IsNumeric(parts[3]) {
				DeleteLeadNextAction(w, r)
				return
			}
			if len(parts) >= 5 && parts[4] == "quote" && helpers.IsNumeric(parts[3]) {
				DeleteLeadQuote(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/event/") {
			if len(parts) >= 5 && parts[4] == "staff" && helpers.IsNumeric(parts[3]) {
				DeleteEventStaff(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/service/") {
			DeleteService(w, r)
			return
		}
	case http.MethodPost:
		parts := strings.Split(path, "/")

		if strings.HasPrefix(path, "/crm/quote-service") {
			PostQuoteService(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/event/") {
			if len(parts) >= 5 && parts[4] == "staff" && helpers.IsNumeric(parts[3]) {
				PostEventStaff(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/lead/") {
			if strings.Contains(path, "quick-quote") {
				PostQuickQuote(w, r)
				return
			}
			if strings.Contains(path, "invoice-reminder") {
				PostSendInvoiceReminder(w, r)
				return
			}
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
			if len(parts) >= 5 && parts[4] == "next-action" && helpers.IsNumeric(parts[3]) {
				PostLeadNextAction(w, r)
				return
			}
		}
		switch path {
		case "/crm/service":
			PostService(w, r)
		case "/crm/user":
			PostUser(w, r)
		case "/crm/cocktail":
			PostCocktail(w, r)
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
	createLeadNextActionForm := constants.PARTIAL_TEMPLATES_DIR + "create_lead_next_action_form.html"
	leadNextActionsTable := constants.PARTIAL_TEMPLATES_DIR + "lead_next_actions_table.html"
	createQuickQuoteForm := constants.PARTIAL_TEMPLATES_DIR + "create_quick_quote_form.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, eventForm, eventTable, leadQuoteForm, leadQuoteTable, createLeadMessageForm, leadMessagesTemplate, createLeadNoteForm, leadNotesTemplate, createLeadNextActionForm, leadNextActionsTable, createQuickQuoteForm}
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

	leadNextActions, err := database.GetLeadNextActionsByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead next actions.", http.StatusInternalServerError)
		return
	}

	alcoholQuoteServices, err := database.GetServiceListByType(constants.AlcoholServiceTypeID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting alcohol quote services.", http.StatusInternalServerError)
		return
	}

	barRentalQuoteServices, err := database.GetServiceListByType(constants.BarRentalServiceTypeID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bar rental quote services.", http.StatusInternalServerError)
		return
	}

	coolerRentalQuoteServices, err := database.GetServiceListByType(constants.CoolerRentalServiceTypeID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bar rental quote services.", http.StatusInternalServerError)
		return
	}

	bartendingAddOnServices, err := database.GetQuickQuoteServiceListByTypeID(constants.BartendingAddOnServiceTypeID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bartending add on services.", http.StatusInternalServerError)
		return
	}

	bartendingHourlyServices, err := database.GetServiceListByType(constants.BartendingServiceTypeID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting bartending service.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Lead Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Lead"] = leadDetails
	data["Events"] = events
	data["Bartenders"] = bartenders
	data["Referrals"] = referrals
	data["LeadQuotes"] = leadQuotes
	data["LeadInterestList"] = leadInterestList
	data["LeadStatusList"] = leadStatusList
	data["NextActionList"] = nextActionList
	data["LeadNotes"] = leadNotes
	data["LeadMessages"] = leadMessages
	data["LeadNextActions"] = leadNextActions
	data["BarRentalQuoteServices"] = barRentalQuoteServices
	data["CoolerRentalQuoteServices"] = coolerRentalQuoteServices
	data["AlcoholQuoteServices"] = alcoholQuoteServices
	data["BartendingAddOnServices"] = bartendingAddOnServices
	data["BartendingHourlyServices"] = bartendingHourlyServices

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

	err = database.UpdateLeadStatus(leadId, constants.ArchivedLeadStatusID)
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
	createEventStaffForm := constants.PARTIAL_TEMPLATES_DIR + "create_event_staff_form.html"
	eventStaffTable := constants.PARTIAL_TEMPLATES_DIR + "event_staff_table.html"
	createEventCocktailsForm := constants.PARTIAL_TEMPLATES_DIR + "create_event_cocktails_form.html"
	eventCocktailsTable := constants.PARTIAL_TEMPLATES_DIR + "event_cocktails_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, createEventStaffForm, eventStaffTable, createEventCocktailsForm, eventCocktailsTable}
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

	eventStaff, err := database.GetEventStaff(eventId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting event staff.", http.StatusInternalServerError)
		return
	}

	users, err := database.GetUsers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users.", http.StatusInternalServerError)
		return
	}

	userRoles, err := database.GetUserRoles()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user roles.", http.StatusInternalServerError)
		return
	}

	eventCocktails, err := database.GetEventCocktails(eventId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting event cocktails.", http.StatusInternalServerError)
		return
	}

	cocktails, err := database.GetCocktails()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cocktails.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Event Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Event"] = eventDetails
	data["EventStaff"] = eventStaff
	data["Users"] = users
	data["UserRoles"] = userRoles
	data["EventCocktails"] = eventCocktails
	data["Cocktails"] = cocktails

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
		err = services.CreateInvoiceWorkflow(quote)

		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error creating invoice.",
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
	createServiceForm := constants.PARTIAL_TEMPLATES_DIR + "create_service_form.html"
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

	pageNum := helpers.ParsePageNum(r.URL.Query().Get("page_num"))
	services, totalRows, err := database.GetServicesList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting services from DB.", http.StatusInternalServerError)
		return
	}

	serviceTypes, err := database.GetServiceTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting service types from DB.", http.StatusInternalServerError)
		return
	}

	unitTypes, err := database.GetUnitTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting unit types from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Services — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Services"] = services
	data["ServiceTypes"] = serviceTypes
	data["UnitTypes"] = unitTypes
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

	pageNum := helpers.ParsePageNum(r.URL.Query().Get("page_num"))
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

	serviceTypes, err := database.GetServiceTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting service types from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "services_table.html",
		Data: map[string]any{
			"Services":     services,
			"ServiceTypes": serviceTypes,
			"CurrentPage":  pageNum,
			"MaxPages":     helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
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

	pageNum := helpers.ParsePageNum(r.URL.Query().Get("page_num"))
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

	serviceTypes, err := database.GetServiceTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting service types from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "services_table.html",
		Data: map[string]any{
			"Services":     services,
			"ServiceTypes": serviceTypes,
			"CurrentPage":  pageNum,
			"MaxPages":     helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
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

	pageNum := helpers.ParsePageNum(r.URL.Query().Get("page_num"))
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

	serviceTypes, err := database.GetServiceTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting service types from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "services_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "services_table.html",
		Data: map[string]any{
			"Services":     services,
			"ServiceTypes": serviceTypes,
			"CurrentPage":  pageNum,
			"MaxPages":     helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLeadQuoteDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lead_quote_detail.html"
	quoteServicesTable := constants.PARTIAL_TEMPLATES_DIR + "quote_services_table.html"
	createQuoteServiceForm := constants.PARTIAL_TEMPLATES_DIR + "create_quote_service_form.html"
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

	values, err := sessions.Get(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user ID from session.", http.StatusInternalServerError)
		return
	}

	leadNote := models.LeadNote{
		LeadID:        leadID,
		Note:          note,
		DateAdded:     time.Now().Unix(),
		AddedByUserID: values.UserID,
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

func GetMessages(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "messages.html"
	leadsWithMessagesTemplate := constants.PARTIAL_TEMPLATES_DIR + "leads_with_messages_list.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, leadsWithMessagesTemplate, baseFile}

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

	messages, err := database.GetLeadsWithMessages()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting messages from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Messages — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["LeadsWithMessages"] = messages

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetMessagesByLeadID(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/message/")
	if err != nil {
		http.Error(w, "Bad lead id.", http.StatusBadRequest)
		return
	}

	leadMessages, err := database.GetMessagesByLeadID(leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead messages from DB.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_messages.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_messages.html",
		Data: map[string]any{
			"LeadMessages": leadMessages,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func SetSMSToRead(w http.ResponseWriter, r *http.Request) {
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

	var form types.SetSMSToReadForm
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

	err = database.SetSMSToRead(helpers.SafeInt(form.MessageID))
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

	leadMessages, err := database.GetMessagesByLeadID(helpers.SafeInt(form.LeadID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting lead messages from DB.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_messages.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_messages.html",
		Data: map[string]any{
			"LeadMessages": leadMessages,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLeadsWithMessages(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	leadsWithMessages, err := database.GetLeadsWithMessages()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting leads with messages from DB.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "leads_with_messages_list.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "leads_with_messages_list.html",
		Data: map[string]any{
			"LeadsWithMessages": leadsWithMessages,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetAutomatedFollowUpMessage(w http.ResponseWriter, r *http.Request) {
	option := r.URL.Query().Get("option")
	leadId := r.URL.Query().Get("leadId")

	// Check for empty parameters
	if option == "" || leadId == "" {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Missing required parameters: option or leadId.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	lead, err := database.GetLeadDetails(leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting lead details from DB.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	previousConversations, err := database.GetPreviousConversations(lead.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting lead conversations from DB.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// Build conversation history string
	var conversationHistory string
	for _, convo := range previousConversations {
		conversationHistory += fmt.Sprintf("[%s] %s\n", convo.Type, convo.Content)
	}

	var prompt string

	// Define the prompts with previous conversations context
	firstFollowUpEnglishPrompt := fmt.Sprintf(`I have a lead whose full name is: %s. This person's original inquiry was: %s. 
	These are our past conversations with this person:
	%s
	This is my first follow-up message. Responding only with the message I'm going to send, write a follow-up message to re-engage the lead about our bartending services. 
	Address this person as cordially as possible while being professional. The message should read like two friends talking to each other through text.
	`, lead.FullName, lead.Message, conversationHistory)

	secondFollowUpEnglishPrompt := fmt.Sprintf(`I have a lead whose full name is: %s. This person's original inquiry was: %s. 
	These are our past conversations with this person:
	%s
	This is my second follow-up message. Responding only with the message I'm going to send, write a follow-up message to re-engage the lead about our bartending services. 
	Address this person as cordially as possible while being professional. The message should read like two friends talking to each other through text. 
	While being cordial, remind the person that we contacted them in the past and if they're not willing to move forward, we'll stop reaching out.
	`, lead.FullName, lead.Message, conversationHistory)

	firstFollowUpSpanishPrompt := fmt.Sprintf(`I have a lead whose full name is: %s. This person's original inquiry was: %s. 
	These are our past conversations with this person:
	%s
	This is my first follow-up message. Responding only with the message I'm going to send, write a follow-up message to re-engage the lead about our bartending services. 
	Address this person as cordially as possible while being professional. The message should read like two friends talking to each other through text, 
	address the person by their first name, maintain an attitude of friendliness & professionalism. Write this text in Spanish. 
	Address this person using "usted," which is a sign of respect.
	`, lead.FullName, lead.Message, conversationHistory)

	secondFollowUpSpanishPrompt := fmt.Sprintf(`I have a lead whose full name is: %s. This person's original inquiry was: %s. 
	These are our past conversations with this person:
	%s
	This is my second follow-up message. Responding only with the message I'm going to send, write a follow-up message to re-engage the lead about our bartending services. 
	Address this person as cordially as possible while being professional. The message should read like two friends talking to each other through text, 
	address the person by their first name, maintain an attitude of friendliness & professionalism. 
	While being cordial, remind the person that we contacted them in the past and if they're not willing to move forward, we'll stop reaching out. 
	Write this text in Spanish. Address this person using "usted," which is a sign of respect.
	`, lead.FullName, lead.Message, conversationHistory)

	switch option {
	case "First Follow Up (ENG)":
		prompt = firstFollowUpEnglishPrompt
	case "Second Follow Up (ENG)":
		prompt = secondFollowUpEnglishPrompt
	case "First Follow Up (ESP)":
		prompt = firstFollowUpSpanishPrompt
	case "Second Follow Up (ESP)":
		prompt = secondFollowUpSpanishPrompt
	default:
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid prompt.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	response, err := services.GetOpenAICompletionsResponse(prompt, 300)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting response from Open AI completions API.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(response))
}

func PostLeadNextAction(w http.ResponseWriter, r *http.Request) {
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

	var form types.LeadNextActionForm
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

	err = database.CreateLeadNextAction(form)
	if err != nil {
		fmt.Printf("Error creating lead next action: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead next action.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadNextActions, err := database.GetLeadNextActionsByLeadID(helpers.SafeInt(form.LeadID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting lead next actions from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_next_actions_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_next_actions_table.html",
		Data: map[string]any{
			"LeadNextActions": leadNextActions,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteLeadNextAction(w http.ResponseWriter, r *http.Request) {
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

	leadNextActionId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("Error getting lead next action id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead next action id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("Error getting lead id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.DeleteLeadNextAction(leadNextActionId)
	if err != nil {
		fmt.Printf("Error deleting lead's next action: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete lead's next action.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadNextActions, err := database.GetLeadNextActionsByLeadID(leadId)
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
			"LeadNextActions": leadNextActions,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostQuickQuote(w http.ResponseWriter, r *http.Request) {
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

	var form types.QuickQuoteForm
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

	var quoteServiceForm []types.QuoteServiceForm
	quoteServices := r.PostForm.Get("quote_services")
	if quoteServices != "" {
		err = json.Unmarshal([]byte(quoteServices), &quoteServiceForm)
		if err != nil {
			fmt.Printf("Error parsing quote_services JSON: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error parsing quote services data.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	quoteId, quoteExternalId, err := database.CreateQuickQuote(form, quoteServiceForm)
	if err != nil {
		fmt.Printf("Error creating quick quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating quick quote.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

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

	err = services.CreateInvoiceWorkflow(quote)
	if err != nil {
		fmt.Printf("Error creating quick quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating quote invoice.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	redirectURL := fmt.Sprintf("%s/external/%s", constants.RootDomain, quoteExternalId)

	var textMessageTemplateNotification = fmt.Sprintf(
		`BARTENDING QUOTE:
		Here's the link to your estimate: %s
	`, redirectURL)

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(quoteId)))
}

func DeleteLeadQuote(w http.ResponseWriter, r *http.Request) {
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

	leadQuoteId, err := helpers.GetSecondIDFromPath(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("Error getting lead quote id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead quote id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("Error getting lead id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.DeleteLeadQuote(leadQuoteId)
	if err != nil {
		fmt.Printf("Error deleting lead's quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete lead's quote.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadQuotes, err := database.GetLeadQuotes(leadId)
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

func GetEvents(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "events.html"
	createEventForm := constants.PARTIAL_TEMPLATES_DIR + "event_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "event_list_view.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createEventForm}

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

	events, totalRows, err := database.GetPaginatedEventList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting services from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Events — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Events"] = events
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostSendInvoiceReminder(w http.ResponseWriter, r *http.Request) {
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

	var externalQuoteView = fmt.Sprintf("%s/external/%s", constants.RootDomain, quote.ExternalID)
	var textMessageTemplateNotification = fmt.Sprintf(
		`BARTENDING FINAL PAYMENT REMINDER:
		Here's the link to the invoice: %s
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
			"AlertMessage": "Reminder has been sent.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostEventStaff(w http.ResponseWriter, r *http.Request) {
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

	var form types.EventStaffForm
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

	err = database.CreateEventStaff(form)
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

	eventStaff, err := database.GetEventStaff(helpers.SafeInt(form.EventID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting event staff from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "event_staff_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "event_staff_table.html",
		Data: map[string]any{
			"EventStaff": eventStaff,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteEventStaff(w http.ResponseWriter, r *http.Request) {
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

	eventStaffId, err := helpers.GetSecondIDFromPath(r, "/crm/event/")
	if err != nil {
		fmt.Printf("Error getting event staff id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get event staff id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	eventId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/event/")
	if err != nil {
		fmt.Printf("Error getting event id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get event id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.DeleteEventStaff(eventStaffId)
	if err != nil {
		fmt.Printf("Error deleting lead's quote: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete lead's quote.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	eventStaff, err := database.GetEventStaff(eventId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting event staff from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "event_staff_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "event_staff_table.html",
		Data: map[string]any{
			"EventStaff": eventStaff,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostEventCocktail(w http.ResponseWriter, r *http.Request) {
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

	var form types.EventCocktailForm
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

	err = database.CreateEventCocktail(form)
	if err != nil {
		fmt.Printf("Error creating event: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating event cocktail.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	eventCocktails, err := database.GetEventCocktails(helpers.SafeInt(form.EventID))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting event cocktails from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "event_cocktails_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "event_cocktails_table.html",
		Data: map[string]any{
			"EventCocktails": eventCocktails,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteEventCocktail(w http.ResponseWriter, r *http.Request) {
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

	eventCocktailId, err := helpers.GetSecondIDFromPath(r, "/crm/event/")
	if err != nil {
		fmt.Printf("Error getting event cocktail id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get event cocktail id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	eventId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/event/")
	if err != nil {
		fmt.Printf("Error getting event id from path: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get event id from path.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.DeleteEventStaff(eventCocktailId)
	if err != nil {
		fmt.Printf("Error deleting event cocktail: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete event cocktail.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	eventCocktails, err := database.GetEventCocktails(eventId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting event cocktails from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "event_cocktails_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "event_cocktails_table.html",
		Data: map[string]any{
			"EventCocktails": eventCocktails,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetUsers(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "users.html"
	createUserForm := constants.PARTIAL_TEMPLATES_DIR + "create_user_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "users_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createUserForm}

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

	users, totalRows, err := database.GetPaginatedUserList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users from DB.", http.StatusInternalServerError)
		return
	}

	userRoles, err := database.GetUserRoles()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user roles from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Users — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Users"] = users
	data["UserRoles"] = userRoles
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
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

	var form types.UserForm
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

	err = database.CreateUser(form)
	if err != nil {
		fmt.Printf("Error creating user: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating user.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	users, totalRows, err := database.GetPaginatedUserList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "users_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "users_table.html",
		Data: map[string]any{
			"Users":       users,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetUserDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "user_detail.html"
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

	userId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/user/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user id from path.", http.StatusInternalServerError)
		return
	}

	userDetails, err := database.GetUserDetails(fmt.Sprint(userId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user details from DB.", http.StatusInternalServerError)
		return
	}

	userRoles, err := database.GetUserRoles()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user roles.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "User Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["User"] = userDetails
	data["UserRoles"] = userRoles

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutUser(w http.ResponseWriter, r *http.Request) {
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

	var form types.UserForm
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

	err = database.UpdateUser(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating user.",
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
			"AlertMessage": "User has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	userId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/user/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteUser(userId)
	if err != nil {
		fmt.Printf("Error deleting event: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete user.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	users, totalRows, err := database.GetPaginatedUserList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "users_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "users_table.html",
		Data: map[string]any{
			"Users":       users,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetCocktails(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "cocktails.html"
	createCocktailForm := constants.PARTIAL_TEMPLATES_DIR + "create_cocktail_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "cocktails_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createCocktailForm}

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

	cocktails, totalRows, err := database.GetPaginatedCocktailList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cocktails from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Cocktails — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Cocktails"] = cocktails
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostCocktail(w http.ResponseWriter, r *http.Request) {
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

	var form types.CreateCocktailForm
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

	err = database.CreateCocktailMany(form)
	if err != nil {
		fmt.Printf("Error creating cocktails: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating cocktails.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	cocktails, totalRows, err := database.GetPaginatedCocktailList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting users from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "cocktails_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "cocktails_table.html",
		Data: map[string]any{
			"Cocktails":   cocktails,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetCocktailDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "cocktail_detail.html"
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

	cocktailId, err := helpers.GetSecondIDFromPath(r, "/crm/cocktail/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cocktail id from path.", http.StatusInternalServerError)
		return
	}

	cocktailDetails, err := database.GetCocktailDetails(fmt.Sprint(cocktailId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cocktail details from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "User Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Cocktail"] = cocktailDetails

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutCocktail(w http.ResponseWriter, r *http.Request) {
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

	var form types.CocktailForm
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

	err = database.UpdateCocktail(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating cocktail.",
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
			"AlertMessage": "Cocktail has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteCocktail(w http.ResponseWriter, r *http.Request) {
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

	cocktailId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/cocktail/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteCocktail(cocktailId)
	if err != nil {
		fmt.Printf("Error deleting cocktail: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete cocktail.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	cocktails, totalRows, err := database.GetPaginatedCocktailList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cocktails from DB.", http.StatusInternalServerError)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "cocktails_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "cocktails_table.html",
		Data: map[string]any{
			"Cocktails":   cocktails,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
