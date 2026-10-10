package pgxdt

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Pgxdt executes queries and runs callbacks within context-bound transactions.
type Pgxdt interface {
	Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (rows pgx.Rows, err error)
	QueryRow(ctx context.Context, sql string, args ...any) (row pgx.Row)
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

// Transactor runs callbacks within a transaction, reusing one when present.
type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
