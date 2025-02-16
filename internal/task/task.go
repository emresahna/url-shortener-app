package task

import (
	"context"

	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/redis/go-redis/v9"
)

type Task interface {
	IncreaseClickTask()
	DeleteExpiredUrl()
}

type task struct {
	db  *sqlc.Queries
	rc  *redis.Client
	ctx context.Context
}

func NewTask(db *sqlc.Queries, rc *redis.Client, ctx context.Context) Task {
	return &task{
		db:  db,
		rc:  rc,
		ctx: ctx,
	}
}
