package postgres

import (
	"context"
	"errors"

	"hub/internal/application/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoTransaction      = errors.New("no transaction in context")
	ErrTransactionStarted = errors.New("transaction already started")
)

// UnitOfWork implements shared.UnitOfWork using PostgreSQL transactions.
type UnitOfWork struct {
	pool *pgxpool.Pool
}

// NewUnitOfWork creates a new UnitOfWork.
func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

// Begin starts a new transaction.
func (u *UnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	if tx := u.getTx(ctx); tx != nil {
		return ctx, ErrTransactionStarted
	}

	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, shared.TxKey(), tx), nil
}

// Commit commits the current transaction.
func (u *UnitOfWork) Commit(ctx context.Context) error {
	tx := u.getTx(ctx)
	if tx == nil {
		return ErrNoTransaction
	}
	return tx.Commit(ctx)
}

// Rollback rolls back the current transaction.
func (u *UnitOfWork) Rollback(ctx context.Context) error {
	tx := u.getTx(ctx)
	if tx == nil {
		return ErrNoTransaction
	}
	return tx.Rollback(ctx)
}

// getTx extracts the transaction from context.
func (u *UnitOfWork) getTx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(shared.TxKey()).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}

// GetTxOrPool returns the transaction from context or the pool.
func GetTxOrPool(ctx context.Context, pool *pgxpool.Pool) Querier {
	if tx, ok := ctx.Value(shared.TxKey()).(pgx.Tx); ok {
		return tx
	}
	return pool
}

// Querier is an interface for database operations.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
