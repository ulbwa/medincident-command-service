package model

import (
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
	return validateUUIDv7(id)
}

func validateOrganizationName(name string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(name); err != nil {
		return err
	}
	if err := validateStringNoConsecutiveSpaces(name); err != nil {
		return err
	}
	if err := validateStringMinLength(name, minOrganizationNameLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(name, maxOrganizationNameLength); err != nil {
		return err
	}

	return nil
}

func validateOrganizationDescription(description string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(description); err != nil {
		return err
	}
	if err := validateStringMinLength(description, minOrganizationDescriptionLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(description, maxOrganizationDescriptionLength); err != nil {
		return err
	}

	return nil
}

func validateOrganization(o *Organization) error {
	if err := validateOrganizationID(o.ID); err != nil {
		return errs.NewInvalidOrganizationError(errs.OrganizationFieldID, err)
	}

	if err := validateOrganizationName(o.Name); err != nil {
		return errs.NewInvalidOrganizationError(errs.OrganizationFieldName, err)
	}

	if o.Description != nil {
		if err := validateOrganizationDescription(*o.Description); err != nil {
			return errs.NewInvalidOrganizationError(errs.OrganizationFieldDescription, err)
		}
	}

	if err := validateAddress(o.LegalAddress); err != nil {
		return errs.NewInvalidOrganizationError(errs.OrganizationFieldLegalAddress, err)
	}

	return nil
}
