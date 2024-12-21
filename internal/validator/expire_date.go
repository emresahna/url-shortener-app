package validator

import (
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"time"
)

func ValidateExpireDate(expireTime *string) (time.Duration, error) {
	if expireTime == nil {
		return 0, nil
	}

	expireParsed, err := time.Parse(time.RFC3339, *expireTime)
	if err != nil {
		return 0, models.ParseExpireTimeErr()
	}

	if time.Now().After(expireParsed) {
		return 0, models.ExpireTimeAlreadyPassedErr()
	}

	return expireParsed.Sub(time.Now()), nil
}
