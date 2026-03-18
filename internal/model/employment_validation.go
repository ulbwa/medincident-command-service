package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minEmploymentPositionLength = 1
	maxEmploymentPositionLength = 255
)

func validateEmploymentPosition(position *string) error {
	if position == nil {
		return nil
	}
	if strings.TrimSpace(*position) != *position {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidEmploymentPosition)
	}
	positionLen := utf8.RuneCountInString(*position)
	if positionLen < minEmploymentPositionLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidEmploymentPosition, minEmploymentPositionLength)
	}
	if positionLen > maxEmploymentPositionLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidEmploymentPosition, maxEmploymentPositionLength)
	}
	return nil
}

func validateDeputyID(deputyID int64) error {
	if err := validateUserID(deputyID); err != nil {
		return fmt.Errorf("%w: invalid deputy id", errs.ErrInvalidEmploymentDeputy)
	}
	return nil
}

func validateEmploymentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidEmploymentID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidEmploymentID)
	}
	return nil
}

func validateEmploymentVacation(vacation *EmploymentVacation) error {
	if vacation == nil {
		return nil
	}
	if vacation.StartsAt.IsZero() {
		return fmt.Errorf("%w: start date is required", errs.ErrInvalidEmploymentVacation)
	}
	if vacation.EndsAt != nil && vacation.EndsAt.Before(vacation.StartsAt) {
		return fmt.Errorf("%w: end date must be after start date", errs.ErrInvalidEmploymentVacation)
	}
	return nil
}

func validateEmployment(e *Employment) error {
	if err := validateEmploymentID(e.ID); err != nil {
		return err
	}
	if err := validateUserID(e.UserID); err != nil {
		return err
	}
	if err := validateOrganizationID(e.OrganizationID); err != nil {
		return err
	}
	if err := validateClinicID(e.ClinicID); err != nil {
		return err
	}
	if err := validateDepartmentID(e.DepartmentID); err != nil {
		return err
	}
	if err := validateEmploymentPosition(e.Position); err != nil {
		return err
	}
	if e.Deputy != nil {
		if err := validateDeputyID(e.Deputy.ID); err != nil {
			return err
		}
	}
	if err := validateEmploymentVacation(e.Vacation); err != nil {
		return err
	}
	return nil
}
