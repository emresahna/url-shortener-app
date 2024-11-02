package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
)

func (s *service) ShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Parse user_id from token
	user_id := 4

	// Shorten url
	shorten_url := "abcdef"

	// Save to db
	saved_url, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
		OriginalUrl:   req.OriginalUrl,
		ShortenedCode: shorten_url,
		UserID:        int32(user_id),
	})
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Save to cache
	err = s.rc.SetUrl(ctx, shorten_url, req.OriginalUrl)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: saved_url.ShortenedCode,
	}, nil
}
