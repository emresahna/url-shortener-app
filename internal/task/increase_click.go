package task

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/hibiken/asynq"
	"strings"
)

func (ts *task) IncreaseClickTask(ctx context.Context, t *asynq.Task) error {
	keys, _ := ts.rc.Keys(ctx, "clicks:*").Result()

	for _, key := range keys {
		shortCode := strings.Split(key, ":")

		urlID, _ := ts.db.GetIDByShortCode(ctx, shortCode[1])

		count, _ := ts.rc.Get(ctx, key).Int64()

		incrementReq := sqlc.IncrementClickCountParams{
			UrlID:       urlID,
			TotalClicks: count,
		}

		_ = ts.db.IncrementClickCount(ctx, incrementReq)

		_ = ts.rc.Del(ctx, key).Err()
	}

	return nil
}
