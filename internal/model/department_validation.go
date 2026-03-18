//nolint:dupl
package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minDepartmentNameLength        = 1
	maxDepartmentNameLength        = 255
	minDepartmentDescriptionLength = 1
	maxDepartmentDescriptionLength = 2000
)

func validateDepartmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidDepartmentID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidDepartmentID)
	}
	return nil
}

func validateDepartmentName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed != name {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidDepartmentName)
	}

	length := utf8.RuneCountInString(name)
	if length < minDepartmentNameLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidDepartmentName, minDepartmentNameLength)
	}
	if length > maxDepartmentNameLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidDepartmentName, maxDepartmentNameLength)
	}

	return nil
}

func validateDepartmentDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed != description {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidDepartmentDescription)
	}

	length := utf8.RuneCountInString(description)
	if length < minDepartmentDescriptionLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidDepartmentDescription, minDepartmentDescriptionLength)
	}
	if length > maxDepartmentDescriptionLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidDepartmentDescription, maxDepartmentDescriptionLength)
	}

	return nil
}

func validateDepartment(d *Department) error {
	if err := validateDepartmentID(d.ID); err != nil {
		return err
	}

	if err := validateClinicID(d.ClinicID); err != nil {
		return err
	}

	if err := validateDepartmentName(d.Name); err != nil {
		return err
	}

	if d.Description != nil {
		if err := validateDepartmentDescription(*d.Description); err != nil {
			return err
		}
	}

	return nil
}
