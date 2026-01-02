package shared

import "context"

// UnitOfWork defines the interface for transaction management.
type UnitOfWork interface {
	// Begin starts a new transaction and returns a context with the transaction.
	Begin(ctx context.Context) (context.Context, error)

	// Commit commits the current transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the current transaction.
	Rollback(ctx context.Context) error
}

// txKey is the context key for storing transactions.
type txKey struct{}

// TxKey returns the context key for transactions.
func TxKey() interface{} {
	return txKey{}
}
