package database

import (
	"context"
	"hub/internal/logger"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type (
	Database interface {
		Pool() *pgxpool.Pool
	}
	database struct {
		pool   *pgxpool.Pool
		logger *logger.Logger
	}
)

func NewDatabase(dns string, minConns, maxConst int, logger *logger.Logger) Database {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dns)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse database config")
	}

	poolConfig.MinConns = int32(minConns)
	poolConfig.MaxConns = int32(maxConst)
	poolConfig.MinIdleConns = int32(2)
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create database connection pool")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.WithError(err).Fatal("Failed to ping database")
	}

	return &database{
		pool:   pool,
		logger: logger,
	}
}

func (d *database) Pool() *pgxpool.Pool {
	return d.pool
}
