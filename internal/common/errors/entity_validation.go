package errors

import "fmt"

type OrganizationField string

const (
	OrganizationFieldID           OrganizationField = "id"
	OrganizationFieldName         OrganizationField = "name"
	OrganizationFieldDescription  OrganizationField = "description"
	OrganizationFieldLegalAddress OrganizationField = "legalAddress"
)

type InvalidOrganizationError struct {
	Field  OrganizationField
	Reason error
}

func (e *InvalidOrganizationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid organization field %s: %s", e.Field, e.Reason)
}

func (e *InvalidOrganizationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidOrganizationError(field OrganizationField, reason error) *InvalidOrganizationError {
	return &InvalidOrganizationError{Field: field, Reason: reason}
}

type ClinicField string

const (
	ClinicFieldID              ClinicField = "id"
	ClinicFieldOrganizationID  ClinicField = "organizationID"
	ClinicFieldName            ClinicField = "name"
	ClinicFieldDescription     ClinicField = "description"
	ClinicFieldPhysicalAddress ClinicField = "physicalAddress"
)

type InvalidClinicError struct {
	Field  ClinicField
	Reason error
}

func (e *InvalidClinicError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid clinic field %s: %s", e.Field, e.Reason)
}

func (e *InvalidClinicError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidClinicError(field ClinicField, reason error) *InvalidClinicError {
	return &InvalidClinicError{Field: field, Reason: reason}
}

type DepartmentField string

const (
	DepartmentFieldID          DepartmentField = "id"
	DepartmentFieldClinicID    DepartmentField = "clinicID"
	DepartmentFieldName        DepartmentField = "name"
	DepartmentFieldDescription DepartmentField = "description"
)

type InvalidDepartmentError struct {
	Field  DepartmentField
	Reason error
}

func (e *InvalidDepartmentError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid department field %s: %s", e.Field, e.Reason)
}

func (e *InvalidDepartmentError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidDepartmentError(field DepartmentField, reason error) *InvalidDepartmentError {
	return &InvalidDepartmentError{Field: field, Reason: reason}
}
