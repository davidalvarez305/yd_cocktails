package utils

import (
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
)

func FormatTimestampEST(timestamp int64) string {
	loc, _ := time.LoadLocation(constants.TimeZone)
	t := time.Unix(timestamp, 0).In(loc)
	formattedTime := t.Format("01/02/2006 03:04:05 PM")
	return formattedTime
}

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	formattedTime := t.Format("01/02/2006 03:04:05 PM")
	return formattedTime
}

func FormatDateMMDDYYYY(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	formattedTime := t.Format("01/02/2006")
	return formattedTime
}

func FormatDateJanDDYYYY(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	formattedTime := t.Format("Jan 2, 2006")
	return formattedTime
}
