package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	configs "github.com/ESE-MONDAY/relay-service/internal/config"
)

func NewPool(cfg *configs.Config) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(
		ctx,
		cfg.DatabaseURL,
	)

	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
