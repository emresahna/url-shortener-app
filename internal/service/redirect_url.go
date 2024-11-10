package service

import (
	"context"
	"errors"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

func (s *service) RedirectUrl(ctx context.Context, code string) (res string, err error) {
	// Check cache
	rc, err := s.rc.GetUrl(ctx, code)
	if !errors.Is(err, redis.Nil) {
		return "", err
	} else if err == nil {
		return rc, nil
	}

	rd, err := s.db.GetURLByCode(ctx, code)
	if err != nil {
		return "", models.UrlNotFound{Message: "Url not found."}
	}

	// Check if the URL is expired in PostgreSQL
	if time.Now().After(*rd.ExpireTime) {
		err := s.db.DeleteURLByID(ctx, rd.ID)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	// Return URL if it hasn't expired
	return rd.OriginalUrl, nil
}
