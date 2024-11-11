package configs

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"io/fs"
	"time"
)

type Config struct {
	HttpConfig     HttpConfig
	PostgresConfig PostgresConfig
	RedisConfig    RedisConfig
	AuthConfig     AuthConfig
}

type HttpConfig struct {
	Address        string        `env:"SERVER_ADDRESS"`
	WriteTimeout   time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	ReadTimeout    time.Duration `env:"SERVER_READ_TIMEOUT"`
	IdleTimeout    time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	MaxHeaderBytes int           `env:"SERVER_MAX_HEADER_BYTES"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     uint16 `env:"POSTGRES_PORT"`
	Database string `env:"POSTGRES_DB"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASS"`
	MaxConn  int32  `env:"POSTGRES_MAX_CONN"`
	MinConn  int32  `env:"POSTGRES_MIN_CONN"`
}

type RedisConfig struct {
	Address string `env:"REDIS_ADDRESS"`
}

type AuthConfig struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		return nil, err
	}

	var hc HttpConfig
	if err = env.Parse(&hc); err != nil {
		return nil, err
	}

	var pc PostgresConfig
	if err = env.Parse(&pc); err != nil {
		return nil, err
	}

	var rc RedisConfig
	if err = env.Parse(&rc); err != nil {
		return nil, err
	}

	var ac AuthConfig
	if err = env.Parse(&ac); err != nil {
		return nil, err
	}

	return &Config{
		HttpConfig:     hc,
		PostgresConfig: pc,
		RedisConfig:    rc,
		AuthConfig:     ac,
	}, nil
}
