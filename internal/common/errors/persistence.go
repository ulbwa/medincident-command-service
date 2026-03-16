package errors

import (
	"errors"
	"fmt"
)

// ErrPersistence is the base error for all persistence and database/transaction related failures.
var ErrPersistence = errors.New("persistence error")

var (
	ErrNotFound                     = fmt.Errorf("%w: entity not found", ErrPersistence)
	ErrTransactionAlreadyExists     = fmt.Errorf("%w: transaction already exists and cannot be created again", ErrPersistence)
	ErrTransactionDead              = fmt.Errorf("%w: transaction is dead and operation cannot be performed", ErrPersistence)
	ErrTransactionAlreadyCommitted  = fmt.Errorf("%w: transaction already committed and cannot be rolled back", ErrPersistence)
	ErrTransactionAlreadyRolledBack = fmt.Errorf("%w: transaction already rolled back and cannot be committed", ErrPersistence)
	ErrSavepointNotFound            = fmt.Errorf("%w: savepoint with given name not found", ErrPersistence)
)
