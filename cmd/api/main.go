package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/emresahna/url-shortener-app/configs"
	"github.com/emresahna/url-shortener-app/internal/auth"
	_ "github.com/emresahna/url-shortener-app/internal/docs"
	"github.com/emresahna/url-shortener-app/internal/handler"
	"github.com/emresahna/url-shortener-app/internal/logger"
	"github.com/emresahna/url-shortener-app/internal/postgres"
	"github.com/emresahna/url-shortener-app/internal/redis"
	"github.com/emresahna/url-shortener-app/internal/service"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"go.uber.org/zap"
)

func main() {
	// initialize global logger
	logger.Init()
	defer logger.Log.Sync()

	logger.Log.Info("Application starting...")

	// load environment file
	cfg, err := configs.Load()
	if err != nil {
		logger.Log.Fatal("Error while loading config", zap.Error(err))
	}

	// initialize jwt client
	jwt, err := auth.NewJWT(cfg.Auth)
	if err != nil {
		logger.Log.Fatal("Error while initializing jwt client", zap.Error(err))
	}

	// initialize redis client for cache
	rcc, err := redis.New(cfg.Redis, cfg.Redis.CacheDB)
	if err != nil {
		logger.Log.Fatal("Error while initializing redis for cache", zap.Error(err))
	}

	// initialize redis client for analytics
	rca, err := redis.New(cfg.Redis, cfg.Redis.AnalyticDB)
	if err != nil {
		logger.Log.Fatal("Error while initializing redis for analytics", zap.Error(err))
	}

	// initialize postgres client
	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Log.Fatal("Error while initializing postgres", zap.Error(err))
	}
	defer db.Close(ctx)

	// initialize sqlc client
	sc := sqlc.New(db)

	// initialize service
	serv := service.New(sc, rcc, jwt, rca)

	// initialize handler
	h := handler.NewHTTP(serv, cfg.Cors)

	// initialize http server
	hs := http.Server{
		Handler:        h,
		Addr:           cfg.Http.Address,
		ReadTimeout:    cfg.Http.ReadTimeout,
		WriteTimeout:   cfg.Http.WriteTimeout,
		IdleTimeout:    cfg.Http.IdleTimeout,
		MaxHeaderBytes: cfg.Http.MaxHeaderBytes,
	}

	go func() {
		logger.Log.Info("HTTP server running.", zap.String("address", cfg.Http.Address))
		if err = hs.ListenAndServe(); err != nil {
			logger.Log.Fatal("Error while starting http server", zap.Error(err))
		}
	}()

	logger.Log.Info("API started...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	logger.Log.Info("Shutting down gracefully...")
}
