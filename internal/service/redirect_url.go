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

	// Check cache
	original, err = s.rc.GetUrl(ctx, code)
	if !errors.Is(err, redis.Nil) && err != nil {
		return "", err
	}

	// Cache not hit, check database
	if original == "" {
		original, err = s.db.GetURLByCode(ctx, code)
		if !errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
	}

	// If url does not appear in db and cache then return not found
	if original == "" {
		return "", models.UrlNotFoundErr()
	}

	err = s.arc.IncreaseClick(ctx, code)
	if err != nil {
		return "", err
	}

	// Return URL
	return original, nil
}
