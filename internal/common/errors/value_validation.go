package errors

import "fmt"

type ValueRequiredError struct{}

func (e *ValueRequiredError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return "value is required"
}

func NewValueRequiredError() *ValueRequiredError {
	return &ValueRequiredError{}
}

type ValueDuplicateError struct {
	Value any
}

func (e *ValueDuplicateError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("value is duplicated: %v", e.Value)
}

func NewValueDuplicateError(value any) *ValueDuplicateError {
	return &ValueDuplicateError{Value: value}
}

type ValueMismatchError struct {
	Expected any
	Actual   any
}

func (e *ValueMismatchError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("value mismatch: expected %v, got %v", e.Expected, e.Actual)
}

func NewValueMismatchError(expected, actual any) *ValueMismatchError {
	return &ValueMismatchError{
		Expected: expected,
		Actual:   actual,
	}
}
