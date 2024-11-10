package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func (s *service) RedirectUrl(ctx context.Context, code string) (res string, err error) {
	// Check cache
	rc, err := s.rc.GetUrl(ctx, code)

	// Get from db
	if errors.Is(err, redis.Nil) {
		rd, err := s.db.GetURLByCode(ctx, code)
		if err != nil {
			return "", err
		}

		// Check if the URL is expired in PostgreSQL
		if time.Now().After(rd.ExpireTime.Time) {
			err := s.db.DeleteURLByID(ctx, rd.ID)
			if err != nil {
				return "", err
			}
			return "", nil
		}

		// Return URL if it hasn't expired
		return rd.OriginalUrl, nil
	}

	if err != nil {
		return "", err
	}

	// Return
	return rc, nil
}
