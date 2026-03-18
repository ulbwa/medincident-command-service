package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minAddressLength = 5
	maxAddressLength = 500
)

func validateAddressValue(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed != value {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidAddressValue)
	}

	length := utf8.RuneCountInString(value)
	if length < minAddressLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidAddressValue, minAddressLength)
	}
	if length > maxAddressLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidAddressValue, maxAddressLength)
	}

	return nil
}
