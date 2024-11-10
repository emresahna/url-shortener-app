package main

import (
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
	jwt, err := auth.NewJWTAuth()
	if err != nil {
		log.Fatal(err)
	}

	// initialize redis client
	rc, err := redis.NewRedisClient(cfg.RedisConfig)
	if err != nil {
		log.Fatal(err)
	}

	// initialize postgres client
	db, err := postgres.NewDBClient(cfg.PostgresConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// initialize sqlc client
	sc := sqlc.New(db)

	// initialize service
	serv := service.NewService(sc, rc, jwt)

	// initialize handler
	h := handler.NewHandler(serv)

	// initialize http server
	hs := http.Server{
		Handler:        h,
		Addr:           cfg.HttpConfig.Address,
		ReadTimeout:    cfg.HttpConfig.ReadTimeout,
		WriteTimeout:   cfg.HttpConfig.WriteTimeout,
		IdleTimeout:    cfg.HttpConfig.IdleTimeout,
		MaxHeaderBytes: cfg.HttpConfig.MaxHeaderBytes,
	}

	fmt.Printf("Server running on %s", cfg.HttpConfig.Address)
	if err = hs.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
