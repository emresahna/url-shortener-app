package postgres

import (
	"context"
	"fmt"
	"github.com/EmreSahna/url-shortener-app/configs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

func NewDBClient(cfg configs.PostgresConfig) (*pgxpool.Pool, error) {
	connDSN := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=%d pool_min_conns=%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.MaxConn, cfg.MinConn)

	config, err := pgxpool.ParseConfig(connDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		log.Printf("Acquiring connection: %p", conn)
		return true
	}

	config.AfterRelease = func(conn *pgx.Conn) bool {
		log.Printf("Releasing connection: %p", conn)
		return true
	}

	config.HealthCheckPeriod = time.Minute * 1
	config.ConnConfig.ConnectTimeout = time.Second * 5

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
