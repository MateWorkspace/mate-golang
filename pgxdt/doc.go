// Package pgxdt routes PostgreSQL operations through a pgx pool or a transaction
// carried by context. WithTx callbacks must pass their context to all database
// operations. Nested callbacks reuse the transaction without creating savepoints.
package pgxdt
