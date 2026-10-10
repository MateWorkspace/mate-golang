# Context-aware PostgreSQL helper

`pgxdt` routes execution through a pgx pool or the transaction attached to the
caller context. It is independent of service models, logging, and configuration.

## Layout

- `pgxdt.go`: public execution and transaction interfaces.
- `db_impl.go`: pool-backed execution and the `NewPgxdt` constructor.
- `transactor.go`: transaction-only implementation and `NewTransactor`.
- `tx_impl.go`: execution through an existing pgx transaction.
- `key.go`: private context key and transaction access helpers.

## Behavior and style

Repositories receive `Pgxdt`; transaction orchestration can receive the smaller
`Transactor` interface. Always pass the callback context to database operations.
`WithTx` starts a transaction when one is absent, commits after callback success,
and attempts rollback on exit. Callback and pgx errors are returned unchanged.
Nested calls reuse the attached transaction and do not create savepoints; return
nested errors to the outer callback when the transaction must roll back.

The context key is package-private. Use the same pool for the database helper and
transactor in a service: an attached transaction takes precedence over the pool,
and this helper does not validate pool ownership. Callers own the pool lifecycle.

Preserve the existing public API and per-responsibility file layout. Use explicit
unexported implementation structs, context-first methods, and early returns.
Keep SQL building and service error classification in consuming projects.
