package service

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignupUser(ctx context.Context, req models.SignupUserRequest) (res models.SignupUserResponse, err error) {
	// Check username already taken
	exists, err := s.db.UserExists(ctx, req.Username)
	if err != nil {
		return models.SignupUserResponse{}, models.InternalServerErr()
	}

	if exists {
		return models.SignupUserResponse{}, models.UserAlreadyExistsErr()
	}

	// Encrypt password
	pw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.SignupUserResponse{}, models.ErrorWhileEncryptingPasswordErr()
	}

	// Save user to db
	savedUser, err := s.db.CreateUser(ctx, sqlc.CreateUserParams{
		Username: req.Username,
		Password: string(pw),
	})
	if err != nil {
		return models.SignupUserResponse{}, models.ErrorWhileCreatingUserErr()
	}

	return models.SignupUserResponse{
		Info: fmt.Sprintf("User %s successfully created.", savedUser.Username),
	}, nil
}
