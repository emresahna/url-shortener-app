package service

import (
	"context"
	"errors"

	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (s *service) TokenRefresh(ctx context.Context, req models.RefreshTokenRequest) (res models.LoginUserResponse, err error) {
	// Parse refresh token
	claims, err := s.jwt.Parse(req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return models.LoginUserResponse{}, models.TokenExpiredErr()
		} else if err != nil {
			return models.LoginUserResponse{}, models.TokenFailureErr()
		}
	}

	// Parse user_id from token
	userUUID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return models.LoginUserResponse{}, models.CustomerIdParseErr()
	}

	// Get user from db
	userDB, err := s.db.GetUserByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			userDB = sqlc.User{}
		} else {
			return models.LoginUserResponse{}, err
		}
	}

	// Credentials are correct, create access and refresh token
	auth, err := s.jwt.Create(userDB)
	if err != nil {
		// Error when creating token
		return models.LoginUserResponse{}, models.ErrorWhileCreatingTokenErr()
	}

	return auth, nil
}
