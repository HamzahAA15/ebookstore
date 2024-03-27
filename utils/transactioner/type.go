package transactioner

import (
	"context"
	"database/sql"
)

type TxxProvider interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type ITransactionProvider interface {
	NewTransaction(ctx context.Context) (TxxProvider, error)
}
