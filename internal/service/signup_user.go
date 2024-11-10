package service

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignupUser(ctx context.Context, req models.SignupUserRequest) (res models.SignupUserResponse, err error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// Error while hashing password
		return models.SignupUserResponse{}, err
	}

	p := sqlc.CreateUserParams{
		Username: req.Username,
		Password: string(pw),
	}

	// Save user to db
	savedUser, err := s.db.CreateUser(ctx, p)
	if err != nil {
		// Error while creating user
		return models.SignupUserResponse{}, err
	}

	return models.SignupUserResponse{
		Info: fmt.Sprintf("User %s successfully created.", savedUser.Username),
	}, nil
}
