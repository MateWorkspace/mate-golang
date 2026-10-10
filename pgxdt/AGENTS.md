# Context-aware PostgreSQL helper

`pgxdt` routes execution through a pgx pool or the transaction attached to the
caller context. It is independent of service models, logging, and configuration.

## Layout

- `pgxdt.go`: public execution and transaction interfaces.
- `db_impl.go`: pool-backed execution and the `NewPgxdt` constructor.
- `transactor.go`: transaction-only implementation and `NewTransactor`.
- `tx_impl.go`: execution through an existing pgx transaction.
- `key.go`: private context key and transaction access helpers.
- `error.go`: the exported active-transaction rejection sentinel.
- `transactor_test.go`: strict transactor and reusable repository behavior tests.

## Behavior and style

Repositories receive `Pgxdt`; transaction orchestration can receive the smaller
`Transactor` interface backed by `NewTransactor`. Always pass the callback context
to database operations. Both constructors' `WithTx` starts a transaction when one
is absent, commits after callback success, and attempts rollback on exit.
Callback and pgx errors are returned unchanged.

`NewTransactor` is strict: an attached transaction returns
`ErrTransactionAlreadyActive` before pool access or callback invocation. It never
starts an independent nested transaction or commits/rolls back the existing one.
`NewPgxdt` and transaction-backed `WithTx` retain reuse for repository composition
and do not create savepoints. A Pgxdt can satisfy Transactor structurally but does
not provide strict behavior; application wiring must choose `NewTransactor`.
Return nested errors to the outer callback when the transaction must roll back;
rejection alone does not mark the outer transaction for rollback.

The context key is package-private. Use the same pool for the database helper and
transactor in a service: an attached transaction takes precedence over the pool,
and this helper does not validate pool ownership. Callers own the pool lifecycle.

Preserve the existing public API and per-responsibility file layout. Use explicit
unexported implementation structs, context-first methods, and early returns.
Keep SQL building and service error classification in consuming projects.

Tests use attached transaction wrappers with nil pools/underlying transactions to
verify rejection and reuse before database access. Run `go test ./...`,
`go vet ./...`, and `go build ./...` from the module root.
