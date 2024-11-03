package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/google/uuid"
)

func (s *service) GetUser(ctx context.Context, userId string) (res models.GetUserResponse, err error) {
	// Parse uuid
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.GetUserResponse{}, err
	}

	user, err := s.db.GetUserByUserID(ctx, userUUID)
	if err != nil {
		return models.GetUserResponse{}, err
	}

	return models.GetUserResponse{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
