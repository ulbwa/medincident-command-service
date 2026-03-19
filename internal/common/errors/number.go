package errors

import "fmt"

type NumberOutOfRangeError struct {
	Min    float64
	Max    float64
	Actual float64
}

func (e *NumberOutOfRangeError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("number is out of range: expected [%v, %v], got %v", e.Min, e.Max, e.Actual)
}

func NewNumberOutOfRangeError(min, max, actual float64) *NumberOutOfRangeError {
	return &NumberOutOfRangeError{Min: min, Max: max, Actual: actual}
}
