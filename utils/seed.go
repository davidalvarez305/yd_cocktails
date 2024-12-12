package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ConvertSeedTransactionTimestamp(day string, hourOfDay string) (int64, error) {
	dateParts := strings.Split(day, "/")
	if len(dateParts) != 3 {
		return 0, errors.New("invalid date format")
	}

	month, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid month format: %w", err)
	}

	dayOfMonth, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid day format: %w", err)
	}

	year, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid year format: %w", err)
	}

	hour, minute, err := parseSeedTransactionHour(hourOfDay)
	if err != nil {
		return 0, fmt.Errorf("error parsing hour: %w", err)
	}

	// EST must be used because that's how Seed is reporting transactions
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return 0, fmt.Errorf("invalid time zone: %w", err)
	}

	date := time.Date(year, time.Month(month), dayOfMonth, hour, minute, 0, 0, location)

	if date.IsZero() {
		return 0, errors.New("invalid date")
	}

	return date.Unix(), nil
}

func parseSeedTransactionHour(hourOfDay string) (int, int, error) {
	parts := strings.Split(hourOfDay, " ")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid hour format")
	}

	timeParts := strings.Split(parts[0], ":")
	hours, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hour value: %w", err)
	}

	minutes := 0
	if len(timeParts) > 1 {
		minutes, err = strconv.Atoi(timeParts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute value: %w", err)
		}
	}

	modifier := strings.ToUpper(parts[1])
	if modifier == "PM" && hours != 12 {
		hours += 12
	} else if modifier == "AM" && hours == 12 {
		hours = 0
	}

	return hours, minutes, nil
}
