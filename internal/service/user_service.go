package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserLogin authenticates a user and returns login response with tokens
func (s *service) UserLogin(ctx context.Context, req models.LoginUserRequest) (res models.LoginUserResponse, err error) {
	// Get user from db
	user, err := s.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// User not found.
		return models.LoginUserResponse{}, models.UserNotFoundErr()
	}

	// User found, check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return models.LoginUserResponse{}, models.PasswordDoesntMatchErr()
	}

	// Credentials are correct, create access and refresh token
	auth, err := s.jwt.Create(user)
	if err != nil {
		// Error when creating token
		return models.LoginUserResponse{}, models.ErrorWhileCreatingTokenErr()
	}

	return auth, nil
}

// UserSignup creates a new user account
func (s *service) UserSignup(ctx context.Context, req models.SignupUserRequest) (res models.SignupUserResponse, err error) {
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

// UserMe returns user profile information and their URLs
func (s *service) UserMe(ctx context.Context) (models.UserResponse, error) {
	var urlsDB []sqlc.GetUrlsByUserRow
	var err error
	var userDB sqlc.User
	var userUUID uuid.UUID

	if token, ok := ctx.Value(models.TokenKey).(string); ok {
		// Parse token
		claims, err := s.jwt.Parse(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return models.UserResponse{}, models.TokenExpiredErr()
			} else {
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
		if ip, ok := ctx.Value(models.IpKey).(string); ok {
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
