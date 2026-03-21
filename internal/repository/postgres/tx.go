// Package postgres provides shared PostgreSQL infrastructure: sqlx-backed
// TransactionFactory and context helpers used by all repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/jmoiron/sqlx"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
)

// txContextKey is the unexported type used as a context key for the active transaction.
type txContextKey struct{}

// savepointNameRegexp restricts savepoint names to safe SQL identifiers, preventing
// injection via fmt.Sprintf since PostgreSQL does not support parameterised SAVEPOINT names.
var savepointNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`)

// TxFactory creates PostgreSQL transactions backed by a sqlx.DB pool.
type TxFactory struct {
	db *sqlx.DB
}

// NewTxFactory returns a persistence.TransactionFactory backed by the given pool.
func NewTxFactory(db *sqlx.DB) *TxFactory {
	return &TxFactory{db: db}
}

// Begin opens a read-committed transaction, stores it in the context, and returns
// the enriched context together with the Transaction handle.
func (f *TxFactory) Begin(ctx context.Context) (context.Context, persistence.Transaction, error) {
	sqlTx, err := f.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction: %w", err)
	}

	t := &transaction{tx: sqlTx}
	txCtx := context.WithValue(ctx, txContextKey{}, t)

	return txCtx, t, nil
}

// TxFromContext returns the active *sqlx.Tx stored in ctx, or nil if none is present.
// Repository implementations use this to execute within the caller's transaction.
func TxFromContext(ctx context.Context) *sqlx.Tx {
	t, _ := ctx.Value(txContextKey{}).(*transaction)
	if t == nil {
		return nil
	}

	return t.tx
}

// transaction implements persistence.Transaction for a PostgreSQL sqlx connection.
type transaction struct {
	tx        *sqlx.Tx
	done      bool
	committed bool
}

// SQLTx exposes the underlying *sqlx.Tx for infrastructure components that
// need direct access to the database transaction (e.g. the outbox dispatcher).
func (t *transaction) SQLTx() *sqlx.Tx {
	return t.tx
}

// Close is idempotent: a no-op if the transaction was already committed or rolled back.
// Called by persistence.WithinTransaction in a defer to release the connection.
func (t *transaction) Close() error {
	if t.done {
		return nil
	}

	// Connection was neither committed nor explicitly rolled back — roll back silently.
	if err := t.tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		return fmt.Errorf("close (implicit rollback): %w", err)
	}

	t.done = true

	return nil
}

// Commit commits the transaction. Returns an error if already finished.
func (t *transaction) Commit(ctx context.Context) error {
	if t.done {
		if t.committed {
			return errs.ErrTransactionAlreadyCommitted
		}

		return errs.ErrTransactionAlreadyRolledBack
	}

	if err := t.tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	t.done = true
	t.committed = true

	return nil
}

// Rollback aborts the transaction. Returns an error if already finished.
func (t *transaction) Rollback(ctx context.Context) error {
	if t.done {
		if t.committed {
			return errs.ErrTransactionAlreadyCommitted
		}

		return errs.ErrTransactionAlreadyRolledBack
	}

	if err := t.tx.Rollback(); err != nil {
		return fmt.Errorf("rollback: %w", err)
	}

	t.done = true

	return nil
}

// Savepoint creates a named savepoint within the transaction.
// The name must match [a-zA-Z_][a-zA-Z0-9_]{0,62} to prevent SQL injection
// (PostgreSQL does not support parameterised savepoint names).
func (t *transaction) Savepoint(ctx context.Context, name string) error {
	if t.done {
		return errs.ErrTransactionDead
	}

	if !savepointNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid savepoint name %q: must be a safe SQL identifier", name)
	}

	if _, err := t.tx.ExecContext(ctx, "SAVEPOINT "+name); err != nil {
		return fmt.Errorf("savepoint %s: %w", name, err)
	}

	return nil
}

// RollbackSavepoint rolls back to a previously created savepoint.
func (t *transaction) RollbackSavepoint(ctx context.Context, name string) error {
	if t.done {
		return errs.ErrTransactionDead
	}

	if !savepointNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid savepoint name %q: must be a safe SQL identifier", name)
	}

	if _, err := t.tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+name); err != nil {
		return fmt.Errorf("rollback savepoint %s: %w", name, err)
	}

	return nil
}
