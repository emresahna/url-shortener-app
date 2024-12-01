package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/EmreSahna/url-shortener-app/internal/validator"
	"log"
	"time"
)

func (s *service) LimitedShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Validate url
	if !validator.ValidateURL(req.OriginalUrl) {
		return models.ShortenURLResponse{}, models.UrlNotValidErr()
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	duration := time.Minute * 10

	// Save to redis
	redisCh := make(chan error, 1)
	go func() {
		log.Printf("Starting to save limited shorten URL %s to Redis...", shortenUrl)
		err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, duration)
		if err != nil {
			redisCh <- models.SaveToCacheErr()
		}
		redisCh <- err
	}()

	// Save to db
	dbCh := make(chan error, 1)
	go func() {
		log.Printf("Starting to save limited shorten URL %s to PostgreSQL...", shortenUrl)
		savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
			OriginalUrl:   req.OriginalUrl,
			ShortenedCode: shortenUrl,
		})
		if err != nil {
			dbCh <- models.CreateURLErr()
		}

		err = s.db.InsertClickCount(context.Background(), savedUrl.ID)
		if err != nil {
			dbCh <- models.CreateClickCountErr()
		}
		dbCh <- nil
	}()

	if err = <-dbCh; err != nil {
		return models.ShortenURLResponse{}, err
	}
	log.Printf("Successfully saved limited shorten URL %s to Redis.", shortenUrl)

	if err = <-redisCh; err != nil {
		return models.ShortenURLResponse{}, err
	}
	log.Printf("Successfully saved limited shorten URL %s to PostgreSQL.", shortenUrl)

	// Return
	return models.ShortenURLResponse{
		Url: shortenUrl,
	}, nil
}
