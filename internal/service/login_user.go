package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (s *service) LoginUser(ctx context.Context, req models.LoginUserRequest) (res models.LoginUserResponse, err error) {
	// Get user from db
	log.Printf("Checking user record for that username...")
	user, err := s.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// User not found.
		return models.LoginUserResponse{}, models.UserNotFoundErr()
	}

	// User found, check password
	log.Printf("User found, comparing password...")
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return models.LoginUserResponse{}, models.PasswordDoesntMatchErr()
	}

	// Credentials are correct, create access and refresh token
	log.Printf("Credentials success, creating tokens...")
	auth, err := s.jwt.Create(user)
	if err != nil {
		// Error when creating token
		return models.LoginUserResponse{}, models.ErrorWhileCreatingTokenErr()
	}

	return auth, nil
}
