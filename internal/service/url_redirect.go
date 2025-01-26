package service

import (
	"context"
	"errors"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/redis/go-redis/v9"
	"log"
)

func (s *service) UrlRedirect(ctx context.Context, code string) (res string, err error) {
	var original string

	// Check redis
	log.Printf("Checking Redis...")
	original, err = s.rcc.GetUrl(ctx, code)
	if !errors.Is(err, redis.Nil) && err != nil {
		return "", models.InternalServerErr()
	}

	// Redis not hit, check database
	if original == "" {
		log.Printf("Checking DB...")
		original, err = s.db.GetURLByCode(ctx, code)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Printf("URL not found...")
				return "", models.UrlNotFoundErr()
			} else {
				return "", models.InternalServerErr()
			}
		}
	}

	// Increase click
	log.Printf("Increasing click count..")
	err = s.rca.IncreaseClick(ctx, code)
	if err != nil {
		return "", models.InternalServerErr()
	}

	// Return URL
	return original, nil
}
