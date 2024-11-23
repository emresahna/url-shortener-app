package service

import (
	"context"
	"errors"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/redis/go-redis/v9"
)

func (s *service) RedirectUrl(ctx context.Context, code string) (res string, err error) {
	var original string

	// Check redis
	original, err = s.rcc.GetUrl(ctx, code)
	if !errors.Is(err, redis.Nil) && err != nil {
		return "", models.InternalServerErr()
	}

	// Redis not hit, check database
	if original == "" {
		original, err = s.db.GetURLByCode(ctx, code)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", models.UrlNotFoundErr()
			} else {
				return "", models.InternalServerErr()
			}
		}
	}

	// Increase click
	err = s.rca.IncreaseClick(ctx, code)
	if err != nil {
		return "", models.InternalServerErr()
	}

	// Return URL
	return original, nil
}
