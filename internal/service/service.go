package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/auth"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/redis"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
)

type Service interface {
	TokenRefresh(context.Context, models.RefreshTokenRequest) (models.LoginUserResponse, error)
	UserMe(context.Context) (models.UserResponse, error)
	UserSignup(context.Context, models.SignupUserRequest) (models.SignupUserResponse, error)
	UserLogin(context.Context, models.LoginUserRequest) (models.LoginUserResponse, error)
	UrlShortenUser(context.Context, models.ShortenURLRequest) (models.ShortenURLResponse, error)
	UrlRemove(context.Context, string) (models.RemoveUrlResponse, error)
	UrlShortenGuest(context.Context, models.ShortenURLRequest) (models.ShortenURLResponse, error)
	UrlRedirect(context.Context, string) (string, error)
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
