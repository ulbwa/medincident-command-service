package errors

import "errors"

var (
	ErrNotFound                     = errors.New("entity not found")
	ErrTransactionAlreadyExists     = errors.New("transaction already exists and cannot be created again")
	ErrTransactionDead              = errors.New("transaction is dead and operation cannot be performed")
	ErrTransactionAlreadyCommitted  = errors.New("transaction already committed and cannot be rolled back")
	ErrTransactionAlreadyRolledBack = errors.New("transaction already rolled back and cannot be committed")
	ErrSavepointNotFound            = errors.New("savepoint with given name not found")
)
