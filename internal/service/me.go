package service

import (
	"context"
	"errors"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (s *service) Me(ctx context.Context) (models.UserResponse, error) {
	var urlsDB []sqlc.GetUrlsByUserRow
	var err error
	var userDB sqlc.User
	var userUUID uuid.UUID

	if token, ok := ctx.Value("token").(string); ok {
		// Parse token
		claims, err := s.jwt.Parse(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return models.UserResponse{}, models.TokenExpiredErr()
			} else if err != nil {
				return models.UserResponse{}, models.TokenFailureErr()
			}
		}

		// Parse user_id from token
		userUUID, err = uuid.Parse(claims["id"].(string))
		if err != nil {
			return models.UserResponse{}, models.CustomerIdParseErr()
		}

		// Get user from db
		userDB, err = s.db.GetUserByUserID(ctx, userUUID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				userDB = sqlc.User{}
			} else {
				return models.UserResponse{}, err
			}
		}

		// Get links from db
		urlsDB, err = s.db.GetUrlsByUser(ctx, sqlc.GetUrlsByUserParams{
			UserID:    &userUUID,
			IpAddress: nil,
		})
		if err != nil {
			return models.UserResponse{}, err
		}
	} else {
		if ip, ok := ctx.Value("ip").(string); ok {
			// Get URLs by IP address
			urlsDB, err = s.db.GetUrlsByUser(ctx, sqlc.GetUrlsByUserParams{
				IpAddress: &ip,
				UserID:    nil,
			})
			if err != nil {
				return models.UserResponse{}, err
			}
		}
	}

	// Convert URLs from DB to response format
	urls := make([]models.UserUrls, len(urlsDB))
	for v, i := range urlsDB {
		urls[v] = models.UserUrls{
			Id:          i.ID,
			Url:         i.ShortenedCode,
			OriginalUrl: i.OriginalUrl,
			ClickCount:  i.TotalClicks,
			IsActive:    i.IsActive,
			IsDeleted:   i.IsDeleted,
		}
	}

	if userDB != (sqlc.User{}) {
		return models.UserResponse{
			Username: userDB.Username,
			Urls:     urls,
		}, nil
	}

	return models.UserResponse{
		Urls: urls,
	}, nil
}
