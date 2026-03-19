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
