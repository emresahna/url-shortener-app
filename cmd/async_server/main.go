package main

import (
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/EmreSahna/url-shortener-app/internal/postgres"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/EmreSahna/url-shortener-app/internal/task"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	// load environment file
	cfg, err := configs.LoadConfig()
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

	// initialize redis client
	rc := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConfig.Address,
		DB:   1,
	})

	ts := task.NewTask(sc, rc)

	// initialize async server
	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr: cfg.RedisConfig.Address,
		DB:   2,
	}, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("click:increase", ts.IncreaseClickTask)

	go func() {
		if err = srv.Run(mux); err != nil {
			log.Fatal(err)
		}
	}()

	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{
		Addr: cfg.RedisConfig.Address,
		DB:   2,
	}, nil)

	_, err = scheduler.Register("@every 5s", asynq.NewTask("click:increase", nil))

	go func() {
		if err = scheduler.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
