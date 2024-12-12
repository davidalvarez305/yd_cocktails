package handlers

import (
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/types"
)

func PartialsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/partials/pop-up-modal":
			GetPopUpModal(w, r)
		case "/partials/error-modal":
			GetErrorModal(w, r)
		case "/partials/opt-out-confirmation-modal":
			GetOptOutConfirmationModal(w, r)
		default:
			http.Error(w, "No partials found.", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetPopUpModal(w http.ResponseWriter, r *http.Request) {
	fileName := "pop_up.html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+fileName)
}

func GetErrorModal(w http.ResponseWriter, r *http.Request) {
	errMessage := r.URL.Query().Get("err")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "error",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
		Data: map[string]any{
			"Message": errMessage,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetOptOutConfirmationModal(w http.ResponseWriter, r *http.Request) {
	fileName := "opt_out_confirmation_modal.html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+fileName)
}
