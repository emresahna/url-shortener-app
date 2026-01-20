package worker

import (
	"time"

	"github.com/emresahna/url-shortener-app/internal/logger"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"go.uber.org/zap"
)

func (w *worker) DeleteExpiredUrls() {
	ctx := w.ctx
	pubSub := w.rc.PSubscribe(ctx, "__keyevent@0__:expired")
	ch := pubSub.Channel()

	defer func() { _ = pubSub.Close() }()

	for msg := range ch {
		now := time.Now()
		logger.Log.Info("Initiating soft deletion for expired URL.", zap.String("url", msg.Payload))
		err := w.db.DeleteExpiredUrlByShortCode(ctx, sqlc.DeleteExpiredUrlByShortCodeParams{
			DeletedAt:     &now,
			ShortenedCode: msg.Payload,
		})
		if err != nil {
			logger.Log.Error(
				"Error during soft deletion of URL",
				zap.String("url", msg.Payload),
				zap.Error(err),
			)
		}
		logger.Log.Info(
			"Soft deletion completed successfully for URL.",
			zap.String("url", msg.Payload),
		)
	}
}
