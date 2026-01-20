package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/emresahna/url-shortener-app/configs"
	"github.com/emresahna/url-shortener-app/internal/logger"
	"github.com/emresahna/url-shortener-app/internal/postgres"
	"github.com/emresahna/url-shortener-app/internal/scheduler"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	// initialize global logger
	logger.Init()
	defer func() { _ = logger.Log.Sync() }()

	logger.Log.Info("Application starting...")

	// ctx
	ctx := context.Background()

	// load environment file
	cfg, err := configs.Load()
	if err != nil {
		logger.Log.Fatal("Error while loading config", zap.Error(err))
	}

	// initialize postgres client
	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Log.Fatal("Error while connecting postgres", zap.Error(err))
	}
	defer func() { _ = db.Close(ctx) }()

	// initialize sqlc client
	sc := sqlc.New(db)

	// initialize redis client cache
	ra := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
		DB:   cfg.AnalyticDB,
	})
	defer func() { _ = ra.Close() }()

	// Ping Redis
	if _, err := ra.Ping(ctx).Result(); err != nil {
		logger.Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	s := scheduler.New(sc, ra, ctx)

	// initialize cron
	c := cron.New()
	_, err = c.AddFunc("@every 10s", s.IncreaseClicks)
	if err != nil {
		logger.Log.Fatal("Failed to schedule IncreaseClickTask", zap.Error(err))
	}
	c.Start()
	defer c.Stop()

	logger.Log.Info("Scheduler server started...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	logger.Log.Info("Shutting down gracefully...")
}
