package errors

import "fmt"

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
