package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/EmreSahna/url-shortener-app/internal/logger"
	"github.com/EmreSahna/url-shortener-app/internal/postgres"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/EmreSahna/url-shortener-app/internal/task"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	// initialize global logger
	logger.InitLogger()
	defer logger.Log.Sync()

	logger.Log.Info("Application starting...")

	// ctx
	ctx := context.Background()

	// load environment file
	cfg, err := configs.LoadConfig()
	if err != nil {
		logger.Log.Fatal("Error while loading config", zap.Error(err))
	}

	// initialize postgres client
	db, err := postgres.NewDBClient(ctx, cfg.Postgres)
	if err != nil {
		logger.Log.Fatal("Error while loading config", zap.Error(err))
	}
	defer db.Close(ctx)

	// initialize sqlc client
	sc := sqlc.New(db)

	// initialize redis client cache
	ra := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
		DB:   cfg.Redis.AnalyticDB,
	})
	defer ra.Close()

	// Ping Redis
	if _, err := ra.Ping(ctx).Result(); err != nil {
		logger.Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	ts := task.NewTask(sc, ra, ctx)

	go ts.DeleteExpiredUrl()

	// initialize cron
	c := cron.New()
	_, err = c.AddFunc("@every 10s", ts.IncreaseClickTask)
	if err != nil {
		logger.Log.Fatal("Failed to schedule IncreaseClickTask", zap.Error(err))
	}
	c.Start()
	defer c.Stop()

	logger.Log.Info("Async server started...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	logger.Log.Info("Shutting down gracefully...")
}
