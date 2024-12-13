package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetSessionExpirationTime() time.Time {
	return time.Now().Add(time.Duration(constants.SessionLength) * 24 * time.Hour)
}

func GenerateTokenExpiryTime() time.Time {
	return time.Now().Add(time.Duration(constants.CSRFTokenLength) * 24 * time.Hour)
}

func GetCurrentTimeInEST() (int64, error) {
	var est int64

	loc, err := time.LoadLocation(constants.TimeZone)
	if err != nil {
		return est, err
	}

	est = time.Now().In(loc).Unix()

	return est, nil
}

func ConvertTimestampToESTDateTime(timestamp int64) (time.Time, error) {
	t := time.Unix(timestamp, 0)

	loc, err := time.LoadLocation(constants.TimeZone)
	if err != nil {
		return t, err
	}

	estTime := t.In(loc)

	return estTime, nil
}

func ParseDateInLocation(timestamp int64) (int64, error) {
	loc, err := time.LoadLocation(constants.TimeZone)
	if err != nil {
		return 0, fmt.Errorf("error loading location: %w", err)
	}

	localTime := time.Unix(timestamp, 0).In(loc)
	return localTime.Unix(), nil
}

func UrlsListHasCurrentPath(urls []string, url string) bool {
	for _, protectedUrl := range urls {
		if strings.Contains(url, protectedUrl) {
			return true
		}
	}
	return false
}

func GetStartAndEndDatesFromMonthYear(monthYear string) (time.Time, time.Time, error) {
	parts := strings.Split(monthYear, ", ")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid format")
	}

	selectedMonth := parts[0]
	selectedYear := parts[1]

	year, err := strconv.Atoi(selectedYear)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year: %v", err)
	}

	monthTime, err := time.Parse("January", selectedMonth)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month: %v", err)
	}

	start := time.Date(year, monthTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return start, end, nil
}

func GetBusinessNameFromURL(location string) (string, error) {
	var businessName string
	parts := strings.Split(location, "/")

	if len(parts) > 3 {
		locationPart := parts[3]

		decodedLocation, err := url.PathUnescape(locationPart)
		if err != nil {
			return "", err
		}

		businessName = decodedLocation
	} else {
		return "", fmt.Errorf("incorrect URL structure")
	}

	return businessName, nil
}

func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")

	c := cases.Title(language.Und)
	for i := range parts {
		parts[i] = c.String(parts[i])
	}

	return strings.Join(parts, "")
}

func AddPhonePrefixIfNeeded(phoneNumber string) string {
	if !strings.HasPrefix(phoneNumber, "+") {
		return "+1" + phoneNumber
	}
	return phoneNumber
}
