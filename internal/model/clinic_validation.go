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
	minClinicNameLength        = 1
	maxClinicNameLength        = 255
	minClinicDescriptionLength = 1
	maxClinicDescriptionLength = 2000
)

func validateClinicID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidClinicID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidClinicID)
	}
	return nil
}

func validateClinicName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed != name {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidClinicName)
	}

	length := utf8.RuneCountInString(name)
	if length < minClinicNameLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidClinicName, minClinicNameLength)
	}
	if length > maxClinicNameLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidClinicName, maxClinicNameLength)
	}

	return nil
}

func validateClinicDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed != description {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidClinicDescription)
	}

	length := utf8.RuneCountInString(description)
	if length < minClinicDescriptionLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidClinicDescription, minClinicDescriptionLength)
	}
	if length > maxClinicDescriptionLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidClinicDescription, maxClinicDescriptionLength)
	}

	return nil
}

func validateClinic(c *Clinic) error {
	if err := validateClinicID(c.ID); err != nil {
		return err
	}

	if err := validateOrganizationID(c.OrganizationID); err != nil {
		return err
	}

	if err := validateClinicName(c.Name); err != nil {
		return err
	}

	if c.Description != nil {
		if err := validateClinicDescription(*c.Description); err != nil {
			return err
		}
	}

	if err := validateAddress(c.PhysicalAddress); err != nil {
		return err
	}

	return nil
}
