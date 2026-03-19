//nolint:dupl
package model

import (
	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minClinicNameLength        = 1
	maxClinicNameLength        = 255
	minClinicDescriptionLength = 1
	maxClinicDescriptionLength = 2000
)

func validateClinicID(id uuid.UUID) error {
	return validateUUIDv7(id)
}

func validateClinicName(name string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(name); err != nil {
		return err
	}
	if err := validateStringNoConsecutiveSpaces(name); err != nil {
		return err
	}
	if err := validateStringMinLength(name, minClinicNameLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(name, maxClinicNameLength); err != nil {
		return err
	}

	return nil
}

func validateClinicDescription(description string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(description); err != nil {
		return err
	}
	if err := validateStringMinLength(description, minClinicDescriptionLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(description, maxClinicDescriptionLength); err != nil {
		return err
	}

	return nil
}

func validateClinic(c *Clinic) error {
	if err := validateClinicID(c.ID); err != nil {
		return errs.NewInvalidClinicError(errs.ClinicFieldID, err)
	}

	if err := validateOrganizationID(c.OrganizationID); err != nil {
		return errs.NewInvalidClinicError(errs.ClinicFieldOrganizationID, err)
	}

	if err := validateClinicName(c.Name); err != nil {
		return errs.NewInvalidClinicError(errs.ClinicFieldName, err)
	}

	if c.Description != nil {
		if err := validateClinicDescription(*c.Description); err != nil {
			return errs.NewInvalidClinicError(errs.ClinicFieldDescription, err)
		}
	}

	if err := validateAddress(c.PhysicalAddress); err != nil {
		return errs.NewInvalidClinicError(errs.ClinicFieldPhysicalAddress, err)
	}

	return nil
}
