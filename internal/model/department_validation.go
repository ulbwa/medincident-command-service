//nolint:dupl
package model

import (
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
	return validateUUIDv7(id)
}

func validateDepartmentName(name string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(name); err != nil {
		return err
	}
	if err := validateStringNoConsecutiveSpaces(name); err != nil {
		return err
	}
	if err := validateStringMinLength(name, minDepartmentNameLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(name, maxDepartmentNameLength); err != nil {
		return err
	}

	return nil
}

func validateDepartmentDescription(description string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(description); err != nil {
		return err
	}
	if err := validateStringMinLength(description, minDepartmentDescriptionLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(description, maxDepartmentDescriptionLength); err != nil {
		return err
	}

	return nil
}

func validateDepartment(d *Department) error {
	if err := validateDepartmentID(d.ID); err != nil {
		return errs.NewInvalidDepartmentError(errs.DepartmentFieldID, err)
	}

	if err := validateClinicID(d.ClinicID); err != nil {
		return errs.NewInvalidDepartmentError(errs.DepartmentFieldClinicID, err)
	}

	if err := validateDepartmentName(d.Name); err != nil {
		return errs.NewInvalidDepartmentError(errs.DepartmentFieldName, err)
	}

	if d.Description != nil {
		if err := validateDepartmentDescription(*d.Description); err != nil {
			return errs.NewInvalidDepartmentError(errs.DepartmentFieldDescription, err)
		}
	}

	return nil
}
