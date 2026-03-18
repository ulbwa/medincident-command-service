package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minOrganizationNameLength        = 1
	maxOrganizationNameLength        = 255
	minOrganizationDescriptionLength = 1
	maxOrganizationDescriptionLength = 2000
)

func validateOrganizationID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidOrganizationID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidOrganizationID)
	}
	return nil
}

func validateOrganizationName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed != name {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidOrganizationName)
	}

	length := utf8.RuneCountInString(name)
	if length < minOrganizationNameLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidOrganizationName, minOrganizationNameLength)
	}
	if length > maxOrganizationNameLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidOrganizationName, maxOrganizationNameLength)
	}

	return nil
}

func validateOrganizationDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed != description {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidOrganizationDescription)
	}

	length := utf8.RuneCountInString(description)
	if length < minOrganizationDescriptionLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidOrganizationDescription, minOrganizationDescriptionLength)
	}
	if length > maxOrganizationDescriptionLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidOrganizationDescription, maxOrganizationDescriptionLength)
	}

	return nil
}

func validateOrganization(o *Organization) error {
	if err := validateOrganizationID(o.ID); err != nil {
		return err
	}

	if err := validateOrganizationName(o.Name); err != nil {
		return err
	}

	if o.Description != nil {
		if err := validateOrganizationDescription(*o.Description); err != nil {
			return err
		}
	}

	return nil
}
