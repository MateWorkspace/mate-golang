package pgxdt

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTransactorRejectsExistingTransaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Nil pool and underlying transaction make accidental database access panic.
	// An attached transaction must be rejected before accessing either one.
	attached := &txImpl{}
	txCtx := putTx(ctx, attached)
	deadline, _ := txCtx.Deadline()
	transactor := NewTransactor(nil)
	called := false

	err := transactor.WithTx(txCtx, func(context.Context) error {
		called = true
		return nil
	})
	if !errors.Is(err, ErrTransactionAlreadyActive) {
		t.Fatalf("expected active-transaction error, got %v", err)
	}
	if called {
		t.Fatal("rejected transaction invoked callback")
	}
	if getTx(txCtx) != attached {
		t.Fatal("rejection changed the attached transaction")
	}
	if got, ok := txCtx.Deadline(); !ok || !got.Equal(deadline) {
		t.Fatal("rejection changed the context deadline")
	}
	if err := txCtx.Err(); err != nil {
		t.Fatalf("rejection canceled caller context: %v", err)
	}

	cancel()
	if !errors.Is(txCtx.Err(), context.Canceled) {
		t.Fatal("transaction context lost parent cancellation")
	}
}

func TestPgxdtReusesExistingTransaction(t *testing.T) {
	callbackErr := errors.New("callback failed")

	for _, test := range []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "callback failure", err: callbackErr},
	} {
		t.Run(test.name, func(t *testing.T) {
			attached := &txImpl{}
			txCtx := putTx(context.Background(), attached)
			db := NewPgxdt(nil)
			called := false

			err := db.WithTx(txCtx, func(received context.Context) error {
				called = true
				if received != txCtx || getTx(received) != attached {
					t.Fatal("repository callback did not reuse caller transaction context")
				}
				return test.err
			})
			if !called {
				t.Fatal("repository transaction skipped callback")
			}
			if err != test.err {
				t.Fatalf("callback error changed: got %v, want %v", err, test.err)
			}
		})
	}
}

func TestNestedTransactorErrorPropagatesThroughRepository(t *testing.T) {
	attached := &txImpl{}
	txCtx := putTx(context.Background(), attached)
	db := NewPgxdt(nil)
	transactor := NewTransactor(nil)
	nestedCalled := false

	err := db.WithTx(txCtx, func(received context.Context) error {
		return transactor.WithTx(received, func(context.Context) error {
			nestedCalled = true
			return nil
		})
	})
	if !errors.Is(err, ErrTransactionAlreadyActive) {
		t.Fatalf("repository swallowed nested rejection: %v", err)
	}
	if nestedCalled {
		t.Fatal("strict nested callback was invoked")
	}
	if getTx(txCtx) != attached {
		t.Fatal("nested rejection replaced caller transaction")
	}
}
