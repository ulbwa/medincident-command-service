package errors

import (
	"errors"
	"fmt"
)

type UserNameField string

const (
	UserNameFieldGivenName  UserNameField = "givenName"
	UserNameFieldFamilyName UserNameField = "familyName"
	UserNameFieldMiddleName UserNameField = "middleName"
)

type InvalidUserNameError struct {
	Field  UserNameField
	Reason error
}

func (e *InvalidUserNameError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid user name field %s: %s", e.Field, e.Reason)
}

func (e *InvalidUserNameError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidUserNameError(field UserNameField, reason error) *InvalidUserNameError {
	return &InvalidUserNameError{Field: field, Reason: reason}
}

type AdminRoleField string

const (
	AdminRoleFieldGrantedAt AdminRoleField = "grantedAt"
	AdminRoleFieldGrantedBy AdminRoleField = "grantedBy"
)

type InvalidAdminRoleError struct {
	Field  AdminRoleField
	Reason error
}

func (e *InvalidAdminRoleError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid admin role field %s: %s", e.Field, e.Reason)
}

func (e *InvalidAdminRoleError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidAdminRoleError(field AdminRoleField, reason error) *InvalidAdminRoleError {
	return &InvalidAdminRoleError{Field: field, Reason: reason}
}

type UserField string

const (
	UserFieldID          UserField = "id"
	UserFieldName        UserField = "name"
	UserFieldAdminRole   UserField = "adminRole"
	UserFieldCustomName  UserField = "customName"
	UserFieldEmployments UserField = "employments"
)

type InvalidUserError struct {
	Field  UserField
	Reason error
}

func (e *InvalidUserError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid user field %s: %s", e.Field, e.Reason)
}

func (e *InvalidUserError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidUserError(field UserField, reason error) *InvalidUserError {
	return &InvalidUserError{Field: field, Reason: reason}
}

var (
	ErrAdminRoleGrantActorNotAdmin      = errors.New("admin role grant forbidden: actor is not admin")
	ErrAdminRoleGrantInsufficientTenure = errors.New("admin role grant forbidden: actor has insufficient tenure")
	ErrUserAlreadyAdmin                 = errors.New("user is already admin")
	ErrUserNotAdmin                     = errors.New("user is not admin")
	ErrAdminSelfRevokeForbidden         = errors.New("admin cannot revoke own role")
)
