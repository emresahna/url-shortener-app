package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/EmreSahna/url-shortener-app/internal/validator"
	"time"
)

func (s *service) LimitedShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Validate url
	if !validator.ValidateURL(req.OriginalUrl) {
		return models.ShortenURLResponse{}, models.UrlNotValidErr()
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	// Save to db
	savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
		OriginalUrl:   req.OriginalUrl,
		ShortenedCode: shortenUrl,
	})
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Create click count record async
	clickCh := make(chan error, 1)
	go func() {
		err = s.db.InsertClickCount(context.Background(), savedUrl.ID)
		if err != nil {
			clickCh <- models.CreateClickCountErr()
		}
		clickCh <- err
	}()

	// Save to cache
	cacheCh := make(chan error, 1)
	go func() {
		duration := time.Minute * 10
		err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, duration)
		if err != nil {
			cacheCh <- models.SaveToCacheErr()
		}
		cacheCh <- err
	}()

	if err = <-clickCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	if err = <-cacheCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: savedUrl.ShortenedCode,
	}, nil
}
