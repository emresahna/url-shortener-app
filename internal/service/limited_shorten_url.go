package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"time"
)

func (s *service) LimitedShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
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

	// Save to cache
	duration := time.Minute * 10
	err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, duration)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: savedUrl.ShortenedCode,
	}, nil
}
