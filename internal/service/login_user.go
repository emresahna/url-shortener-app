package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) LoginUser(ctx context.Context, req models.LoginUserRequest) (res models.LoginUserResponse, err error) {
	// Get user from db
	user, err := s.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// User not found.
		return models.LoginUserResponse{}, err
	}

	// User found, check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Password not correct.
		return models.LoginUserResponse{}, err
	}

	// Credentials are correct, create access and refresh token
	auth, err := s.jwt.Create(user)
	if err != nil {
		// Error when creating token
		return models.LoginUserResponse{}, err
	}

	return auth, nil
}
