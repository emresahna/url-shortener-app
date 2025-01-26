package main

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/EmreSahna/url-shortener-app/internal/auth"
	"github.com/EmreSahna/url-shortener-app/internal/handler"
	"github.com/EmreSahna/url-shortener-app/internal/postgres"
	"github.com/EmreSahna/url-shortener-app/internal/redis"
	"github.com/EmreSahna/url-shortener-app/internal/service"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"log"
	"net/http"
)

func main() {
	// load environment file
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// initialize jwt client
	jwt, err := auth.NewJWTAuth(cfg.Auth)
	if err != nil {
		log.Fatal(err)
	}

	// initialize redis client for cache
	rcc, err := redis.NewRedisClient(cfg.Redis, cfg.Redis.CacheDB)
	if err != nil {
		log.Fatal(err)
	}

	// initialize redis client for analytics
	rca, err := redis.NewRedisClient(cfg.Redis, cfg.Redis.AnalyticDB)
	if err != nil {
		log.Fatal(err)
	}

	// initialize postgres client
	ctx := context.Background()
	db, err := postgres.NewDBClient(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	// initialize sqlc client
	sc := sqlc.New(db)

	// initialize service
	serv := service.NewService(sc, rcc, jwt, rca)

	// initialize handler
	h := handler.NewHttpHandler(serv, cfg.Cors)

	// initialize http server
	hs := http.Server{
		Handler:        h,
		Addr:           cfg.Http.Address,
		ReadTimeout:    cfg.Http.ReadTimeout,
		WriteTimeout:   cfg.Http.WriteTimeout,
		IdleTimeout:    cfg.Http.IdleTimeout,
		MaxHeaderBytes: cfg.Http.MaxHeaderBytes,
	}

	fmt.Printf("Server running on %s", cfg.Http.Address)
	if err = hs.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
