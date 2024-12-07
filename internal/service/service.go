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
	Refresh(context.Context, models.RefreshTokenRequest) (models.LoginUserResponse, error)
	Me(ctx context.Context) (models.UserResponse, error)
	ShortenURL(context.Context, models.ShortenURLRequest) (models.ShortenURLResponse, error)
	LimitedShortenURL(context.Context, models.ShortenURLRequest) (models.ShortenURLResponse, error)
	RedirectUrl(context.Context, string) (string, error)
}

type service struct {
	db  *sqlc.Queries
	rcc redis.Store
	jwt auth.Auth
	rca redis.Store
}

func NewService(db *sqlc.Queries, rcc redis.Store, jwt auth.Auth, rca redis.Store) Service {
	return &service{
		db:  db,
		rcc: rcc,
		jwt: jwt,
		rca: rca,
	}
}
