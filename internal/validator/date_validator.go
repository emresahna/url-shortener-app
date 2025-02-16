package validator

import (
	"time"

	"github.com/EmreSahna/url-shortener-app/internal/models"
)

func ParseDateWithTimeZone(date string) (time.Time, error) {
	dateParsed, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Time{}, models.ParseExpireTimeErr()
	}

	return dateParsed, nil
}

func ValidateFutureDate(parsedDate time.Time) (time.Duration, error) {
	now := time.Now().UTC()
	diff := parsedDate.Sub(now)

	if diff < 0 {
		return 0, models.ExpireTimeAlreadyPassedErr()
	}

	return diff, nil
}
