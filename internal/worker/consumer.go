package worker

import (
	"context"

	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/redis/go-redis/v9"
)

type Worker interface {
	DeleteExpiredUrls()
}

type worker struct {
	db  *sqlc.Queries
	rc  *redis.Client
	ctx context.Context
}

func New(db *sqlc.Queries, rc *redis.Client, ctx context.Context) Worker {
	return &worker{
		db:  db,
		rc:  rc,
		ctx: ctx,
	}
}
