package model

import (
	"strings"
	"unicode/utf8"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func validateStringMinLength(s string, minLength int) error {
	if utf8.RuneCountInString(s) < minLength {
		return errs.NewStringTooShortError(minLength, utf8.RuneCountInString(s), s)
	}
	return nil
}

func validateStringMaxLength(s string, maxLength int) error {
	if utf8.RuneCountInString(s) > maxLength {
		return errs.NewStringTooLongError(maxLength, utf8.RuneCountInString(s), s)
	}
	return nil
}

func validateStringNoLeadingOrTrailingWhitespace(s string) error {
	expected := strings.TrimSpace(s)
	if expected != s {
		return errs.NewStringLeadingOrTrailingWhitespaceError(s, expected)
	}
	return nil
}

func validateStringNoConsecutiveSpaces(s string) error {
	if !strings.Contains(s, "  ") {
		return nil
	}
	expected := strings.Join(strings.Fields(s), " ")
	return errs.NewStringConsecutiveSpacesError(s, expected)
}
