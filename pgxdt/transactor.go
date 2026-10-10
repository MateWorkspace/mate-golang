package pgxdt

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type transactor struct {
	pool *pgxpool.Pool
}

// NewTransactor wraps pool for transaction orchestration without owning it.
// WithTx rejects contexts with an attached transaction so successful calls own
// their commit boundary. Use NewPgxdt for reusable repository transactions.
func NewTransactor(pool *pgxpool.Pool) Transactor {
	return &transactor{pool: pool}
}

func (t *transactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if getTx(ctx) != nil {
		return ErrTransactionAlreadyActive
	}

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	wrapped := &txImpl{tx: tx}
	ctx = putTx(ctx, wrapped)

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
