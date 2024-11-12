package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/google/uuid"
)

func (s *service) ShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	t := ctx.Value("token").(string)
	c, err := s.jwt.Parse(t)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Parse user_id from token
	userUUID, err := uuid.Parse(c["id"].(string))
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	// Save to db
	savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
		OriginalUrl:   req.OriginalUrl,
		ShortenedCode: shortenUrl,
		UserID:        &userUUID,
		ExpireTime:    req.ExpireTime,
	})
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Create click count record
	err = s.db.InsertClickCount(context.Background(), savedUrl.ID)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Save to cache
	err = s.rcc.SetUrl(ctx, shortenUrl, req.OriginalUrl)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: savedUrl.ShortenedCode,
	}, nil
}
