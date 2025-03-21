package utils

import (
	"time"

	"github.com/davidalvarez305/yd_cocktails/types"
)

func FormatTimestampWithOptions(timestamp int64, opts *types.TimestampFormatOptions) string {
	if opts == nil {
		opts = &types.TimestampFormatOptions{
			Format:   "01/02/2006 03:04:05 PM", // Default format
			TimeZone: "",                       // Default to UTC
		}
	}

	var loc *time.Location
	var err error
	if opts.TimeZone != "" {
		loc, err = time.LoadLocation(opts.TimeZone)
		if err != nil {
			loc = time.UTC
		}
	} else {
		loc = time.UTC
	}

	t := time.Unix(timestamp, 0)
	t = t.In(loc)

	return t.Format(opts.Format)
}
