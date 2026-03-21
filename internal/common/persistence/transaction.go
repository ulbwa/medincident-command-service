package persistence

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/rs/zerolog"
)

type TransactionFactory interface {
	Begin(ctx context.Context) (context.Context, Transaction, error)
}

type Transaction interface {
	io.Closer
	Savepoint(ctx context.Context, name string) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	RollbackSavepoint(ctx context.Context, name string) error
}

func WithinTransaction(ctx context.Context, tx Transaction, fn func() error) (err error) {
	defer func() {
		if closeErr := tx.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close transaction")
		}
	}()

	if err = fn(); err != nil {
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()

		rollbackErr := tx.Rollback(rollbackCtx)
		if rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
