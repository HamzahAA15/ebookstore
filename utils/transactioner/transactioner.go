package transactioner

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionProvider struct {
	db *sqlx.DB
}

func NewTransactionProvider(db *sqlx.DB) TransactionProvider {
	return TransactionProvider{
		db: db,
	}
}

func (p *TransactionProvider) NewTransaction(ctx context.Context) (*sqlx.Tx, error) {
	return p.db.BeginTxx(ctx, nil)
}
