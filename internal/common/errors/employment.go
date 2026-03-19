package errors

import (
	"errors"
	"fmt"
)

type EmploymentField string

const (
	EmploymentFieldID             EmploymentField = "id"
	EmploymentFieldUserID         EmploymentField = "userID"
	EmploymentFieldOrganizationID EmploymentField = "organizationID"
	EmploymentFieldClinicID       EmploymentField = "clinicID"
	EmploymentFieldDepartmentID   EmploymentField = "departmentID"
	EmploymentFieldPosition       EmploymentField = "position"
	EmploymentFieldDeputy         EmploymentField = "deputy"
	EmploymentFieldVacation       EmploymentField = "vacation"
	EmploymentFieldAssignedAt     EmploymentField = "assignedAt"
)

type InvalidEmploymentError struct {
	Field  EmploymentField
	Reason error
}

func (e *InvalidEmploymentError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid employment field %s: %s", e.Field, e.Reason)
}

func (e *InvalidEmploymentError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidEmploymentError(field EmploymentField, reason error) *InvalidEmploymentError {
	return &InvalidEmploymentError{Field: field, Reason: reason}
}

type EmploymentVacationField string

const (
	EmploymentVacationFieldStartsAt EmploymentVacationField = "startsAt"
	EmploymentVacationFieldEndsAt   EmploymentVacationField = "endsAt"
)

type InvalidEmploymentVacationError struct {
	Field  EmploymentVacationField
	Reason error
}

func (e *InvalidEmploymentVacationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid employment vacation field %s: %s", e.Field, e.Reason)
}

func (e *InvalidEmploymentVacationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidEmploymentVacationError(field EmploymentVacationField, reason error) *InvalidEmploymentVacationError {
	return &InvalidEmploymentVacationError{Field: field, Reason: reason}
}

type EmploymentDeputyField string

const (
	EmploymentDeputyFieldID EmploymentDeputyField = "id"
)

type InvalidEmploymentDeputyError struct {
	Field  EmploymentDeputyField
	Reason error
}

func (e *InvalidEmploymentDeputyError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid employment deputy field %s: %s", e.Field, e.Reason)
}

func (e *InvalidEmploymentDeputyError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidEmploymentDeputyError(field EmploymentDeputyField, reason error) *InvalidEmploymentDeputyError {
	return &InvalidEmploymentDeputyError{Field: field, Reason: reason}
}

var (
	ErrUserCannotBeOwnDeputy                 = errors.New("user cannot be their own deputy")
	ErrEmploymentVacationAlreadyExists       = errors.New("employment vacation already exists")
	ErrEmploymentAlreadyExistsInOrganization = errors.New("employment already exists in organization")
	ErrEmploymentNotFound                    = errors.New("employment not found")
)
