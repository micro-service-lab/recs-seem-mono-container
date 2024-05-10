// Package pgadapter provides a Postgres adapter for the store package.
package pgadapter

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
)

// PgAdapter is a Postgres adapter for the store package.
type PgAdapter struct {
	mu      sync.RWMutex
	query   *query.Queries
	pool    *pgxpool.Pool
	clocker clock.Clock
	qtxMap  map[store.Sd]*query.Queries
	txMap   map[store.Sd]pgx.Tx
}

var _ store.Store = (*PgAdapter)(nil)

// NewPgAdapter creates a new PgAdapter.
func NewPgAdapter(ctx context.Context, clocker clock.Clock, cfg config.Config) (*PgAdapter, error) {
	return NewPgAdapterFromConnStr(ctx, cfg.DBUrl, clocker)
}

// NewPgAdapterFromConnStr creates a new PgAdapter from a connection string.
func NewPgAdapterFromConnStr(ctx context.Context, connStr string, clocker clock.Clock) (*PgAdapter, error) {
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	q := query.New(pool)
	return &PgAdapter{
		pool:    pool,
		query:   q,
		clocker: clocker,
		qtxMap:  make(map[store.Sd]*query.Queries),
		txMap:   make(map[store.Sd]pgx.Tx),
	}, nil
}

// Begin starts a transaction.
func (a *PgAdapter) Begin(ctx context.Context) (store.Sd, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return store.Sd(uuid.Nil), fmt.Errorf("failed to begin transaction: %w", err)
	}
	id := store.Sd(uuid.New())
	a.txMap[id] = tx
	a.qtxMap[id] = a.query.WithTx(tx)
	return id, nil
}

// Commit commits a transaction.
func (a *PgAdapter) Commit(ctx context.Context, id store.Sd) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.qtxMap[id]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	delete(a.qtxMap, id)
	tx, ok := a.txMap[id]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	delete(a.txMap, id)
	err := tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back a transaction.
func (a *PgAdapter) Rollback(ctx context.Context, id store.Sd) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.qtxMap[id]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	delete(a.qtxMap, id)
	tx, ok := a.txMap[id]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	delete(a.txMap, id)
	err := tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

// Cleanup cleans up the store.
func (a *PgAdapter) Cleanup(_ context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for id := range a.qtxMap {
		delete(a.qtxMap, id)
	}
	for id := range a.txMap {
		delete(a.txMap, id)
	}
	a.pool.Close()
	return nil
}
