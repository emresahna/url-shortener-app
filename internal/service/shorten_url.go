package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/google/uuid"
)

func (s *service) ShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Parse user_id from token
	userId := ctx.Value("userID").(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Shorten url
	shortenUrl := "abcdef"

	// Save to db
	savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
		OriginalUrl:   req.OriginalUrl,
		ShortenedCode: shortenUrl,
		UserID:        userUUID,
	})
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Save to cache
	err = s.rc.SetUrl(ctx, shortenUrl, req.OriginalUrl)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: savedUrl.ShortenedCode,
	}, nil
}
