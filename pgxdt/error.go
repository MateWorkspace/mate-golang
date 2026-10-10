package pgxdt

import "errors"

// ErrTransactionAlreadyActive indicates that a strict transactor received a
// context with an attached transaction. Its callback has not been invoked.
var ErrTransactionAlreadyActive = errors.New("pgxdt: transaction already active")
