package errors

import "fmt"

type UserNameField string

const (
	UserNameFieldGivenName  UserNameField = "givenName"
	UserNameFieldFamilyName UserNameField = "familyName"
	UserNameFieldMiddleName UserNameField = "middleName"
)

type UserNameValidationReason string

const (
	UserNameValidationReasonTooShort                    UserNameValidationReason = "too_short"
	UserNameValidationReasonTooLong                     UserNameValidationReason = "too_long"
	UserNameValidationReasonLeadingOrTrailingWhitespace UserNameValidationReason = "leading_or_trailing_whitespace"
)

type UserNameValidationDetails struct {
	MinLength                      int
	MaxLength                      int
	ActualLength                   int
	ActualValue                    string
	HasLeadingOrTrailingWhitespace bool
}

type InvalidUserNameError struct {
	Field   UserNameField
	Reason  UserNameValidationReason
	Details UserNameValidationDetails
}

func (e *InvalidUserNameError) Error() string {
	if e == nil {
		return "<nil>"
	}

	base := fmt.Sprintf("invalid user name field %s", e.Field)

	switch e.Reason {
	case UserNameValidationReasonTooShort:
		return fmt.Sprintf("%s: expected minimum length %d symbols, got %d symbols", base, e.Details.MinLength, e.Details.ActualLength)
	case UserNameValidationReasonTooLong:
		return fmt.Sprintf("%s: expected maximum length %d symbols, got %d symbols", base, e.Details.MaxLength, e.Details.ActualLength)
	case UserNameValidationReasonLeadingOrTrailingWhitespace:
		return fmt.Sprintf("%s: leading or trailing whitespace is not allowed", base)
	default:
		return base
	}
}

func NewUserNameTooShortError(field UserNameField, minLength, actualLength int, actualValue string) *InvalidUserNameError {
	return &InvalidUserNameError{
		Field:  field,
		Reason: UserNameValidationReasonTooShort,
		Details: UserNameValidationDetails{
			MinLength:    minLength,
			ActualLength: actualLength,
			ActualValue:  actualValue,
		},
	}
}

func NewUserNameTooLongError(field UserNameField, maxLength, actualLength int, actualValue string) *InvalidUserNameError {
	return &InvalidUserNameError{
		Field:  field,
		Reason: UserNameValidationReasonTooLong,
		Details: UserNameValidationDetails{
			MaxLength:    maxLength,
			ActualLength: actualLength,
			ActualValue:  actualValue,
		},
	}
}

func NewUserNameLeadingOrTrailingWhitespaceError(field UserNameField, actualValue string) *InvalidUserNameError {
	return &InvalidUserNameError{
		Field:  field,
		Reason: UserNameValidationReasonLeadingOrTrailingWhitespace,
		Details: UserNameValidationDetails{
			ActualValue:                    actualValue,
			HasLeadingOrTrailingWhitespace: true,
		},
	}
}
