package task

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Task interface {
	IncreaseClickTask(context.Context, *asynq.Task) error
}

type task struct {
	db *sqlc.Queries
	rc *redis.Client
}

func NewTask(db *sqlc.Queries, rc *redis.Client) Task {
	return &task{
		db: db,
		rc: rc,
	}
}
