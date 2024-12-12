package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashBytes)
	return hashedString
}

func FormatPhoneNumber(phoneNumber string) string {
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")
	if len(cleaned) == 10 {
		return fmt.Sprintf("(%s) %s - %s", cleaned[:3], cleaned[3:6], cleaned[6:])
	}
	return ""
}

func GetGoogleClientIDFromRequest(r *http.Request) (string, error) {
	gaCookie, err := r.Cookie("_ga")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _ga cookie found")
		}
		return "", err
	}

	parts := strings.Split(gaCookie.Value, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("unexpected _ga cookie format")
	}

	return parts[2] + "." + parts[3], nil
}

func GetFacebookClickIDFromRequest(r *http.Request) (string, error) {
	fbcCookie, err := r.Cookie("_fbc")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _fbc cookie found")
		}
		return "", err
	}

	return fbcCookie.Value, nil
}

func GetFacebookClientIDFromRequest(r *http.Request) (string, error) {
	fbpCookie, err := r.Cookie("_fbp")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _fbp cookie found")
		}
		return "", err
	}

	return fbpCookie.Value, nil
}
