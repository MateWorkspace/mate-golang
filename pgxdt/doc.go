// Package pgxdt routes PostgreSQL operations through a pgx pool or a transaction
// carried by context. WithTx callbacks must pass their context to all database
// operations. NewTransactor rejects an attached transaction; NewPgxdt reuses it
// for repository operations without creating savepoints. Propagate errors from
// either helper when the outer transaction should roll back.
package pgxdt
