package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
)

func createExternalViewContext() map[string]any {
	return map[string]any{
		"PageTitle":       constants.CompanyName,
		"MetaDescription": "YD Cocktails quote details.",
		"SiteName":        constants.SiteName,
		"StaticPath":      constants.StaticPath,
		"MediaPath":       constants.MediaPath,
		"PhoneNumber":     constants.CompanyPhoneNumber,
		"CurrentYear":     time.Now().Year(),
		"CompanyName":     constants.CompanyName,
	}
}

var externalBaseFilePath = constants.EXTERNAL_TEMPLATES_DIR + "base.html"

func ExternalHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx := createExternalViewContext()
	ctx["PagePath"] = constants.RootDomain + path

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/quote/":
			GetExternalQuoteDetails(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetExternalQuoteDetails(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	headerPath := "header_desktop.html"
	if helpers.IsMobileRequest(r) {
		headerPath = "header_mobile.html"
	}

	fileName := "external_quote_details.html"
	files := []string{externalBaseFilePath, websiteFooterFilePath, constants.EXTERNAL_TEMPLATES_DIR + fileName, constants.WEBSITE_TEMPLATES_DIR + headerPath}
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

	externalQuoteId := strings.TrimPrefix(r.URL.Path, "/quote/")

	quote, err := database.GetExternalQuoteDetails(externalQuoteId)
	if err != nil {
		http.Error(w, "Error retrieving quote details.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Quote View — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Quote"] = quote
	data["BartendingRate"] = constants.BartendingRate

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
