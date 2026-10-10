# mate-golang

Small, reusable Go libraries for Mate projects.

`mate-golang` collects focused utilities that can be shared across services.
Each library has its own public package, with no dependency on service-specific
models, configuration, or application architecture.

## Libraries

| Package | Purpose |
| --- | --- |
| [`pgxdt`](./pgxdt) | Context-aware PostgreSQL execution and transactions using pgx v5. |

### pgxdt

- **Transparent execution:** `Exec`, `Query`, and `QueryRow` use the pool normally
  and the attached transaction inside a `WithTx` callback.
- **Transaction callbacks:** commit on success and attempt rollback on failure.
- **Nested reuse:** nested callbacks share the existing transaction without
  creating savepoints.
- **Small interfaces:** inject `Pgxdt` for database operations or `Transactor` for
  transaction orchestration.
- **Native errors:** callback and pgx errors are returned unchanged.

## Usage

Requires Go **1.26.5** or newer. Import the package:

```go
import "github.com/MateWorkspace/mate-golang/pgxdt"
```

Given a caller context and a `*pgxpool.Pool`, create both helpers from the same
pool and pass the callback context to every database operation:

```go
db := pgxdt.NewPgxdt(pool)
transactor := pgxdt.NewTransactor(pool)

err := transactor.WithTx(ctx, func(txCtx context.Context) error {
	_, err := db.Exec(txCtx,
		"UPDATE accounts SET balance = balance - $1 WHERE id = $2",
		amount, senderID,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(txCtx,
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		amount, recipientID,
	)
	return err
})
if err != nil {
	return err
}
```

The caller owns the pool lifecycle. Nested transaction errors must propagate to
the outer callback when the transaction should roll back. An attached transaction
takes precedence over the pool; pool ownership is not checked.

### Local development

Until a version is published, use a sibling checkout from a consuming project.
Add the package import first, then run:

```sh
go mod edit -require=github.com/MateWorkspace/mate-golang@v0.0.0
go mod edit -replace=github.com/MateWorkspace/mate-golang=../mate-golang
go mod tidy
```

Once releases are published, depend on a tagged version. All packages share the
module's release cycle.

## Project layout

```text
mate-golang/
├── pgxdt/       # PostgreSQL execution and transaction helpers
├── AGENTS.md    # Layout and contribution conventions
├── go.mod
└── go.sum
```

Future libraries belong in sibling package directories at the module root.

## Checks

Run from the repository root:

```sh
go build ./...
go test ./...
go vet ./...
```

Format Go changes with `gofmt`. The current package does not yet contain test
files.
