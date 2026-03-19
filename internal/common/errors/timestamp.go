package errors

import (
	"fmt"
	"time"
)

type TimestampBeforeMinimumError struct {
	Actual        time.Time
	ExpectedAfter time.Time
}

func (e *TimestampBeforeMinimumError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"expected timestamp after %s, actual %s",
		e.ExpectedAfter.Format(time.RFC3339Nano),
		e.Actual.Format(time.RFC3339Nano),
	)
}

func NewTimestampBeforeMinimumError(actual, expectedAfter time.Time) *TimestampBeforeMinimumError {
	return &TimestampBeforeMinimumError{
		Actual:        actual,
		ExpectedAfter: expectedAfter,
	}
}

type TimestampAfterMaximumError struct {
	Actual         time.Time
	ExpectedBefore time.Time
}

func (e *TimestampAfterMaximumError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"expected timestamp before %s, actual %s",
		e.ExpectedBefore.Format(time.RFC3339Nano),
		e.Actual.Format(time.RFC3339Nano),
	)
}

func NewTimestampAfterMaximumError(actual, expectedBefore time.Time) *TimestampAfterMaximumError {
	return &TimestampAfterMaximumError{
		Actual:         actual,
		ExpectedBefore: expectedBefore,
	}
}
