package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/google/uuid"
	"time"
)

func (s *service) RemoveUrl(ctx context.Context, id string) (res models.RemoveUrlResponse, err error) {
	// Parse url id
	urlUUID, err := uuid.Parse(id)
	if err != nil {
		return models.RemoveUrlResponse{}, err
	}

	// Get url from db
	code, err := s.db.GetURLByID(ctx, urlUUID)
	if err != nil {
		return models.RemoveUrlResponse{}, err
	}

	// Delete from redis
	if err = s.rcc.DeleteUrl(ctx, code); err != nil {
		return models.RemoveUrlResponse{}, err
	}

	// Soft delete from db
	now := time.Now()
	if err = s.db.DeleteExpiredUrlByShortCode(ctx, sqlc.DeleteExpiredUrlByShortCodeParams{
		DeletedAt:     &now,
		ShortenedCode: code,
	}); err != nil {
		return models.RemoveUrlResponse{}, err
	}

	return models.RemoveUrlResponse{
		Message: "Url successfully removed.",
	}, nil
}
