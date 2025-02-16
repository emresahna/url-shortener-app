package task

import (
	"errors"
	"strings"

	"github.com/EmreSahna/url-shortener-app/internal/logger"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (ts *task) IncreaseClickTask() {
	ctx := ts.ctx
	keys, _ := ts.rc.Keys(ctx, "clicks:*").Result()

	for _, key := range keys {
		count, err := ts.rc.Get(ctx, key).Int64()
		if errors.Is(redis.Nil, err) {
			logger.Log.Error("Key not found in Redis", zap.String("key", key))
		} else if err != nil {
			logger.Log.Error("Error while getting key from Redis", zap.String("key", key), zap.Error(err))
		}

		parts := strings.Split(key, ":")
		if len(parts) < 2 {
			logger.Log.Error("Key parts are missing.", zap.String("key", key))
		}

		shortCode := parts[1]

		incrementReq := sqlc.IncrementClickCountParams{
			ShortenedCode: shortCode,
			TotalClicks:   count,
		}

		err = ts.db.IncrementClickCount(ctx, incrementReq)
		if err != nil {
			logger.Log.Error("Error while incrementing click counts.", zap.String("key", key), zap.Error(err))
		}

		err = ts.rc.Del(ctx, key).Err()
		if err != nil {
			logger.Log.Error("Error while deleting click counts.", zap.String("key", key), zap.Error(err))
		}

		logger.Log.Info("Increasing click counts is success.", zap.String("key", key))
	}
}
