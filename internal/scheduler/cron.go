package scheduler

import (
	"context"

	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/redis/go-redis/v9"
)

type Scheduler interface {
	IncreaseClicks()
}

type scheduler struct {
	db  *sqlc.Queries
	rc  *redis.Client
	ctx context.Context
}

func New(db *sqlc.Queries, rc *redis.Client, ctx context.Context) Scheduler {
	return &scheduler{
		db:  db,
		rc:  rc,
		ctx: ctx,
	}
}
