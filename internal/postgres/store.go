package postgres

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDBClient(cfg configs.PostgresConfig) (*pgxpool.Pool, error) {
	connDSN := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=%d pool_min_conns=%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.MaxConn, cfg.MinConn)

	conn, err := pgxpool.Connect(context.Background(), connDSN)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
