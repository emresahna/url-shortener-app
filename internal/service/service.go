package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/auth"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/redis"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
)

type Service interface {
	SignupUser(context.Context, models.SignupUserRequest) (models.SignupUserResponse, error)
	LoginUser(context.Context, models.LoginUserRequest) (models.LoginUserResponse, error)
	ShortenURL(context.Context, models.ShortenURLRequest) (models.ShortenURLResponse, error)
	RedirectUrl(context.Context, string) (string, error)
}

type service struct {
	db  *sqlc.Queries
	rc  redis.Store
	jwt auth.Auth
}

func NewService(db *sqlc.Queries, rc redis.Store, jwt auth.Auth) Service {
	return &service{
		db:  db,
		rc:  rc,
		jwt: jwt,
	}
}
