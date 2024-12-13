package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_cocktails/sessions"
	"github.com/davidalvarez305/yd_cocktails/utils"
)

var urlsToSkip = []string{"/static/", "/partials/", "/sms/", "/call/", "/webhooks/"}

func UserTracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.UrlsListHasCurrentPath(urlsToSkip, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		isNew, err := sessions.IsNew(r)
		if err != nil {
			fmt.Printf("FAILED TO CHECK SESSION AT USER TRACKING: %+v\n", err)
			http.Error(w, "Unable to check if session is new.", http.StatusInternalServerError)
			return
		}

		var externalId, csrfSecret string

		if isNew {
			session, err := sessions.Create(r, w)
			if err != nil {
				fmt.Printf("FAILED TO CREATE SESSION AT USER TRACKING: %+v\n", err)
				http.Error(w, "Failed to create session.", http.StatusInternalServerError)
				return
			}

			expirationTime := time.Unix(session.DateExpires, 0).UTC()

			sessions.SetCookie(w, expirationTime, session.CSRFSecret)

			externalId = session.ExternalID
			csrfSecret = session.CSRFSecret
		}

		if !isNew {
			session, err := sessions.Get(r)
			if err != nil {
				fmt.Printf("FAILED TO RETRIEVE SESSION AT USER TRACKING: %+v\n", err)
				http.Error(w, "Failed to retrieve session in user middleware.", http.StatusInternalServerError)
				return
			}

			externalId = session.ExternalID
			csrfSecret = session.CSRFSecret
		}

		r = r.WithContext(context.WithValue(r.Context(), "external_id", externalId))
		r = r.WithContext(context.WithValue(r.Context(), "csrf_secret", csrfSecret))
		next.ServeHTTP(w, r)
	})
}
