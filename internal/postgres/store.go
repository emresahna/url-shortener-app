package postgres

import (
	"context"
	"fmt"

	"github.com/emresahna/url-shortener-app/configs"
	"github.com/jackc/pgx/v4"
)

func New(ctx context.Context, cfg configs.Postgres) (*pgx.Conn, error) {
	connDSN := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	config, err := pgx.ParseConfig(connDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	if cfg.DetailedLogging {
		config.LogLevel = pgx.LogLevelTrace
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
