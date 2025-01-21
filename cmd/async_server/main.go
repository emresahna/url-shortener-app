package main

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/EmreSahna/url-shortener-app/internal/postgres"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/EmreSahna/url-shortener-app/internal/task"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	// load environment file
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// initialize postgres client
	ctx := context.Background()
	db, err := postgres.NewDBClient(ctx, cfg.PostgresConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	// initialize sqlc client
	sc := sqlc.New(db)

	// initialize redis client cache
	ra := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConfig.Address,
		DB:   cfg.RedisConfig.AnalyticDB,
	})

	pubSub := ra.PSubscribe(ctx, "__keyevent@0__:expired")
	ch := pubSub.Channel()

	defer pubSub.Close()

	go func() {
		for msg := range ch {
			now := time.Now()
			log.Printf("Initiating soft deletion for expired URL: %s\n", msg.Payload)
			err := sc.DeleteExpiredUrlByShortCode(context.TODO(), sqlc.DeleteExpiredUrlByShortCodeParams{
				DeletedAt:     &now,
				ShortenedCode: msg.Payload,
			})
			if err != nil {
				log.Printf("Error during soft deletion of URL %s: %v\n", msg.Payload, err)
			}
			log.Printf("Soft deletion completed successfully for URL: %s\n", msg.Payload)
		}
	}()

	// initialize redis client analytics
	rc := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConfig.Address,
		DB:   cfg.RedisConfig.AnalyticDB,
	})

	ts := task.NewTask(sc, rc)

	// initialize async server
	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr: cfg.RedisConfig.Address,
		DB:   cfg.RedisConfig.SchedulerDB,
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
		DB:   cfg.RedisConfig.SchedulerDB,
	}, nil)

	_, err = scheduler.Register("@every 10s", asynq.NewTask("click:increase", nil))

	go func() {
		if err = scheduler.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
