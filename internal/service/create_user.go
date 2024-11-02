package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
)

func (s *service) CreateUser(ctx context.Context, req models.CreateUserRequest) (res models.CreateUserResponse, err error) {
	// Save user to db
	saved_user, err := s.db.CreateUser(ctx, req.Username)
	if err != nil {
		return models.CreateUserResponse{}, err
	}

	return models.CreateUserResponse{
		Username: saved_user.Username,
	}, nil
}
