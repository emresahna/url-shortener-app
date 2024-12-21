package task

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/hibiken/asynq"
	"log"
	"strings"
)

func (ts *task) IncreaseClickTask(ctx context.Context, t *asynq.Task) error {
	keys, _ := ts.rc.Keys(ctx, "clicks:*").Result()

	for _, key := range keys {
		count, err := ts.rc.Get(ctx, key).Int64()
		if err != nil {
			log.Printf("Key not found in Redis: %s\n", key)
		}

		parts := strings.Split(key, ":")
		if len(parts) < 2 {
			log.Printf("Key parts are missing: %s\n", key)
		}

		shortCode := parts[1]

		incrementReq := sqlc.IncrementClickCountParams{
			ShortenedCode: shortCode,
			TotalClicks:   count,
		}

		err = ts.db.IncrementClickCount(ctx, incrementReq)
		if err != nil {
			log.Printf("Error while incrementing click counts: %s\n", key)
		}

		err = ts.rc.Del(ctx, key).Err()
		if err != nil {
			log.Printf("Error while deleting click counts: %s\n", key)
		}
	}

	return nil
}
