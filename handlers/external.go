package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/database"
	"github.com/davidalvarez305/yd_cocktails/helpers"
	"github.com/davidalvarez305/yd_cocktails/services"
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
		if strings.HasPrefix(path, "/external/") {
			GetExternalQuoteDetails(w, r, ctx)
			return
		}
		switch path {
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

	externalQuoteId := strings.TrimPrefix(r.URL.Path, "/external/")

	quote, err := database.GetExternalQuoteDetails(externalQuoteId)
	if err != nil {
		fmt.Printf("ERROR GETTING QUOTE DETAILS: %+v\n", err)
		http.Error(w, "Error retrieving quote details.", http.StatusInternalServerError)
		return
	}

	// I have to do this for now because I don't know if it's a good idea or not to save the quote as a column on invoices
	if quote.IsDepositPaid {
		inv, err := database.GetRemainingInvoice(quote.QuoteID)
		if err != nil {
			fmt.Printf("ERROR GETTING REMAINING INVOICE FROM DB: %+v\n", err)
			http.Error(w, "Error retrieving quote details.", http.StatusInternalServerError)
			return
		}

		remainingInvoice, err := services.GetStripeInvoice(inv.StripeInvoiceID)
		if err != nil {
			fmt.Printf("ERROR GETTING REMAINING INVOICE FROM STRIPE: %+v\n", err)
			http.Error(w, "Error retrieving quote details.", http.StatusInternalServerError)
			return
		}

		quote.RemainingAmount = float64(remainingInvoice.AmountDue / 100)
	}

	quoteServices, err := database.GetQuoteServices(quote.QuoteID)
	if err != nil {
		fmt.Printf("ERROR GETTING QUOTE SERVICES: %+v\n", err)
		http.Error(w, "Error retrieving quote services.", http.StatusInternalServerError)
		return
	}

	isWithin48Hours := false
	t := time.Unix(quote.EventDateTimestamp, 0)
	currentTime := time.Now()

	// Check if the event date is within 48 hours from now
	if t.Sub(currentTime) <= 48*time.Hour && t.After(currentTime) {
		isWithin48Hours = true
	}

	data := ctx
	data["PageTitle"] = "Quote View — " + constants.CompanyName
	data["Nonce"] = nonce
	data["Quote"] = quote
	data["QuoteServices"] = quoteServices
	data["IsWithin48Hours"] = isWithin48Hours

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
