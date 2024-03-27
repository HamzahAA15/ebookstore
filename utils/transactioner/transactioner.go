package transactioner

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionProvider struct {
	db *sqlx.DB
}

func NewTransactionProvider(db *sqlx.DB) ITransactionProvider {
	return &TransactionProvider{
		db: db,
	}
}

func (p *TransactionProvider) NewTransaction(ctx context.Context) (TxxProvider, error) {
	return p.db.BeginTxx(ctx, nil)
}
