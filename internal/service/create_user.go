package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
)

func (s *service) CreateUser(ctx context.Context, req models.CreateUserRequest) (res models.CreateUserResponse, err error) {
	// Save user to db
	savedUser, err := s.db.CreateUser(ctx, req.Username)
	if err != nil {
		return models.CreateUserResponse{}, err
	}

	return models.CreateUserResponse{
		Username: savedUser.Username,
	}, nil
}
