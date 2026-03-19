package errors

import "fmt"

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
