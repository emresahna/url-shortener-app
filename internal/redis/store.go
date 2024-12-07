package redis

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	ipAddrUsage = "free-tier-usage-%s"
)

type Store interface {
	Ping(context.Context) error
	SetUrlWithExpire(context.Context, string, string, time.Duration) error
	GetUrl(context.Context, string) (string, error)
	IncreaseClick(context.Context, string) error
	SetIpAddrUsage(context.Context, string) error
	IncreaseIpAddrUsage(context.Context, string) error
	GetIpAddrUsage(context.Context, string) (int, error)
}

type store struct {
	rcc *redis.Client
}

func NewRedisClient(cfg configs.RedisConfig, db int) (Store, error) {
	rcc := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		DB:           db,
		PoolSize:     20,
		MinIdleConns: 5,
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
	err := c.rcc.Set(ctx, shortenUrl, originalUrl, expireTime).Err()
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

func (c *store) SetIpAddrUsage(ctx context.Context, ip string) error {
	err := c.rcc.SetEx(ctx, fmt.Sprintf(ipAddrUsage, ip), 1, time.Hour*24*30).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *store) IncreaseIpAddrUsage(ctx context.Context, ip string) error {
	err := c.rcc.Incr(ctx, fmt.Sprintf(ipAddrUsage, ip)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *store) GetIpAddrUsage(ctx context.Context, ip string) (int, error) {
	count, err := c.rcc.Get(ctx, fmt.Sprintf(ipAddrUsage, ip)).Int()
	if err != nil {
		return 0, err
	}
	return count, nil
}
