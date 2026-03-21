package model

import (
	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minEmploymentPositionLength = 1
	maxEmploymentPositionLength = 255
)

func validateEmploymentPosition(position string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(position); err != nil {
		return err
	}
	if err := validateStringNoConsecutiveSpaces(position); err != nil {
		return err
	}
	if err := validateStringMinLength(position, minEmploymentPositionLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(position, maxEmploymentPositionLength); err != nil {
		return err
	}
	return nil
}

func validateEmploymentDeputy(deputy EmploymentDeputy) error {
	if err := validateUserID(deputy.ID); err != nil {
		return errs.NewInvalidEmploymentDeputyError(errs.EmploymentDeputyFieldID, err)
	}
	return nil
}

func validateDeputyID(deputyID uuid.UUID) error {
	return validateEmploymentDeputy(EmploymentDeputy{ID: deputyID})
}

func validateEmploymentID(id uuid.UUID) error {
	return validateUUIDv7(id)
}

func validateEmploymentVacation(vacation *EmploymentVacation) error {
	if vacation == nil {
		return nil
	}
	if vacation.StartsAt.IsZero() {
		return errs.NewInvalidEmploymentVacationError(errs.EmploymentVacationFieldStartsAt, errs.NewValueRequiredError())
	}
	if vacation.EndsAt != nil {
		if err := validateTimestampNotBefore(*vacation.EndsAt, vacation.StartsAt); err != nil {
			return errs.NewInvalidEmploymentVacationError(errs.EmploymentVacationFieldEndsAt, err)
		}
	}
	return nil
}

func validateEmployment(e *Employment) error {
	if err := validateEmploymentID(e.ID); err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldID, err)
	}
	if err := validateUserID(e.UserID); err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldUserID, err)
	}
	if e.AssignedAt.IsZero() {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldAssignedAt, errs.NewValueRequiredError())
	}
	if err := validateOrganizationID(e.OrganizationID); err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldOrganizationID, err)
	}
	if err := validateClinicID(e.ClinicID); err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldClinicID, err)
	}
	if err := validateDepartmentID(e.DepartmentID); err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldDepartmentID, err)
	}
	if e.Position != nil {
		if err := validateEmploymentPosition(*e.Position); err != nil {
			return errs.NewInvalidEmploymentError(errs.EmploymentFieldPosition, err)
		}
	}
	if e.Deputy != nil {
		if err := validateEmploymentDeputy(*e.Deputy); err != nil {
			return errs.NewInvalidEmploymentError(errs.EmploymentFieldDeputy, err)
		}
		if e.Deputy.ID == e.UserID {
			return errs.NewInvalidEmploymentError(errs.EmploymentFieldDeputy, errs.ErrUserCannotBeOwnDeputy)
		}
	}
	if err := validateEmploymentVacation(e.Vacation); err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldVacation, err)
	}
	return nil
}
