package redis

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/redis/go-redis/v9"
)

type Store interface {
	Ping(context.Context) error
	SetUrl(context.Context, string, string) error
}

type store struct {
	rcc *redis.Client
}

func NewRedisClient(cfg configs.RedisConfig) (Store, error) {
	rcc := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
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

func (c *store) SetUrl(ctx context.Context, shortenUrl string, originalUrl string) error {
	err := c.rcc.Set(ctx, shortenUrl, originalUrl, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
