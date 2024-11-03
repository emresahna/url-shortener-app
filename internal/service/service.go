package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/redis"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
)

type Service interface {
	CreateUser(context.Context, models.CreateUserRequest) (models.CreateUserResponse, error)
	ShortenURL(context.Context, models.ShortenURLRequest) (models.ShortenURLResponse, error)
	GetUser(context.Context, string) (models.GetUserResponse, error)
}

type service struct {
	db *sqlc.Queries
	rc redis.Store
}

func NewService(db *sqlc.Queries, rc redis.Store) Service {
	return &service{
		db: db,
		rc: rc,
	}
}
