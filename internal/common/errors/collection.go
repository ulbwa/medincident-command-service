package errors

import "fmt"

type InvalidCollectionItemError struct {
	Index  int
	Reason error
}

func (e *InvalidCollectionItemError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid collection item at index %d: %s", e.Index, e.Reason)
}

func (e *InvalidCollectionItemError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidCollectionItemError(index int, reason error) *InvalidCollectionItemError {
	return &InvalidCollectionItemError{Index: index, Reason: reason}
}
