package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DBStore struct {
	pool *pgxpool.Pool
}

func NewDBStore(ctx context.Context, dsn string) (*DBStore, error) {
	pcfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pcfg.MaxConns = 20
	pcfg.MinConns = 2
	pcfg.MaxConnLifetime = 30 * time.Minute
	pcfg.MaxConnIdleTime = 5 * time.Minute
	pcfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, err
	}

	return &DBStore{pool}, nil
}

func (db *DBStore) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DBStore) Pool() *pgxpool.Pool {
	return db.pool
}
