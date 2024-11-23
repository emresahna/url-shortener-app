package redis

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/redis/go-redis/v9"
	"time"
)

type Store interface {
	Ping(context.Context) error
	SetUrlWithExpire(context.Context, string, string, time.Duration) error
	GetUrl(context.Context, string) (string, error)
	IncreaseClick(context.Context, string) error
}

type store struct {
	rcc *redis.Client
}

func NewRedisClient(cfg configs.RedisConfig, db int) (Store, error) {
	rcc := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		DB:   db,
	})

	s := &store{
		rcc: rcc,
	}

	return s, nil
}

func (c *store) Ping(ctx context.Context) error {
	err := c.rcc.Ping(ctx).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *store) SetUrlWithExpire(ctx context.Context, shortenUrl string, originalUrl string, expireTime time.Duration) error {
	err := c.rcc.SetEx(ctx, shortenUrl, originalUrl, expireTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *store) GetUrl(ctx context.Context, code string) (string, error) {
	result, err := c.rcc.Get(ctx, code).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *store) IncreaseClick(ctx context.Context, code string) error {
	err := c.rcc.Incr(ctx, "clicks:"+code).Err()
	if err != nil {
		return err
	}
	return nil
}
