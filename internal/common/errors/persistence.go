package errors

import (
	"errors"
)

// ErrPersistence is the base error for all persistence and database/transaction related failures.
var ErrPersistence = errors.New("persistence error")
