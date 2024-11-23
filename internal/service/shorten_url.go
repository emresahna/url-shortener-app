package service

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/google/uuid"
	"time"
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
	duration := time.Second * 0
	if req.ExpireTime != nil {
		expirationTime, err := time.Parse(time.RFC3339, *req.ExpireTime)
		if err != nil {
			fmt.Println("Error parsing expiration date:", err)
			return models.ShortenURLResponse{}, err
		}

		duration = expirationTime.Sub(time.Now())
	}

	err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, duration)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: savedUrl.ShortenedCode,
	}, nil
}
