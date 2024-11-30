package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/google/uuid"
)

func (s *service) Me(ctx context.Context, userId string) (models.UserResponse, error) {
	// Parse sent userId
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get user from db
	userDB, err := s.db.GetUserByUserID(ctx, userUUID)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get links from db
	urlsDB, err := s.db.GetUrlsByUserID(ctx, &userUUID)
	if err != nil {
		return models.UserResponse{}, err
	}

	urls := make([]models.UserUrls, len(urlsDB))
	for v, i := range urlsDB {
		urls[v] = models.UserUrls{
			Url:         i.ShortenedCode,
			OriginalUrl: i.OriginalUrl,
			ClickCount:  i.TotalClicks,
		}
	}

	return models.UserResponse{
		Username: userDB.Username,
		Urls:     urls,
	}, nil
}
