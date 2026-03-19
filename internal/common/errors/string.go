package errors

import "fmt"

type StringTooShortError struct {
	MinLength    int
	ActualLength int
	ActualValue  string
}

func (e *StringTooShortError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("string is too short: expected minimum length %d symbols, got %d symbols", e.MinLength, e.ActualLength)
}

func NewStringTooShortError(minLength, actualLength int, actualValue string) *StringTooShortError {
	return &StringTooShortError{
		MinLength:    minLength,
		ActualLength: actualLength,
		ActualValue:  actualValue,
	}
}

type StringTooLongError struct {
	MaxLength    int
	ActualLength int
	ActualValue  string
}

func (e *StringTooLongError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("string is too long: expected maximum length %d symbols, got %d symbols", e.MaxLength, e.ActualLength)
}

func NewStringTooLongError(maxLength, actualLength int, actualValue string) *StringTooLongError {
	return &StringTooLongError{
		MaxLength:    maxLength,
		ActualLength: actualLength,
		ActualValue:  actualValue,
	}
}

type StringLeadingOrTrailingWhitespaceError struct {
	ActualValue   string
	ExpectedValue string
}

func (e *StringLeadingOrTrailingWhitespaceError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("string has leading or trailing whitespace: expected %q, got %q", e.ExpectedValue, e.ActualValue)
}

func NewStringLeadingOrTrailingWhitespaceError(actualValue, expectedValue string) *StringLeadingOrTrailingWhitespaceError {
	return &StringLeadingOrTrailingWhitespaceError{ActualValue: actualValue, ExpectedValue: expectedValue}
}

type StringConsecutiveSpacesError struct {
	ActualValue   string
	ExpectedValue string
}

func (e *StringConsecutiveSpacesError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("string has consecutive spaces: expected %q, got %q", e.ExpectedValue, e.ActualValue)
}

func NewStringConsecutiveSpacesError(actualValue, expectedValue string) *StringConsecutiveSpacesError {
	return &StringConsecutiveSpacesError{ActualValue: actualValue, ExpectedValue: expectedValue}
}
