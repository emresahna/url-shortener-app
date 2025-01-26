package configs

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"io/fs"
	"time"
)

type Config struct {
	Http
	Postgres
	Redis
	Auth
	Cors
}

type Http struct {
	Address        string        `env:"SERVER_ADDRESS"`
	WriteTimeout   time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	ReadTimeout    time.Duration `env:"SERVER_READ_TIMEOUT"`
	IdleTimeout    time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	MaxHeaderBytes int           `env:"SERVER_MAX_HEADER_BYTES"`
}

type Postgres struct {
	Host            string `env:"POSTGRES_HOST"`
	Port            uint16 `env:"POSTGRES_PORT"`
	Database        string `env:"POSTGRES_DB"`
	User            string `env:"POSTGRES_USER"`
	Password        string `env:"POSTGRES_PASS"`
	DetailedLogging bool   `env:"POSTGRES_DETAILED_LOGGING"`
}

type Redis struct {
	Address      string `env:"REDIS_ADDRESS"`
	CacheDB      int    `env:"REDIS_CACHE_DB"`
	AnalyticDB   int    `env:"REDIS_ANALYTIC_DB"`
	SchedulerDB  int    `env:"REDIS_SCHEDULER_DB"`
	PoolSize     int    `env:"REDIS_POOL_SIZE"`
	MinIdleConns int    `env:"REDIS_MIN_IDLE_CONNS"`
}

type Auth struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
}

type Cors struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
	AllowedMethods []string `env:"ALLOWED_METHODS"`
	AllowedHeaders []string `env:"ALLOWED_HEADERS"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		return nil, err
	}

	var hc Http
	if err = env.Parse(&hc); err != nil {
		return nil, err
	}

	var pc Postgres
	if err = env.Parse(&pc); err != nil {
		return nil, err
	}

	var rc Redis
	if err = env.Parse(&rc); err != nil {
		return nil, err
	}

	var ac Auth
	if err = env.Parse(&ac); err != nil {
		return nil, err
	}

	var cc Cors
	if err = env.Parse(&ac); err != nil {
		return nil, err
	}

	return &Config{
		Http:     hc,
		Postgres: pc,
		Redis:    rc,
		Auth:     ac,
		Cors:     cc,
	}, nil
}
