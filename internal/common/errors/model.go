package errors

import (
	"errors"
	"fmt"
)

// ErrInvariantViolation indicates that a domain rule or invariant was broken.
var ErrInvariantViolation = errors.New("invariant violation")

// User & Identity errors
var (
	ErrInvalidUserID              = fmt.Errorf("%w: invalid user id", ErrInvariantViolation)
	ErrInvalidUserGivenName       = fmt.Errorf("%w: invalid given name", ErrInvariantViolation)
	ErrInvalidUserFamilyName      = fmt.Errorf("%w: invalid family name", ErrInvariantViolation)
	ErrInvalidUserMiddleName      = fmt.Errorf("%w: invalid middle name", ErrInvariantViolation)
	ErrInvalidUserEmployment      = fmt.Errorf("%w: invalid user employment", ErrInvariantViolation)
	ErrInvalidAdminRoleSince      = fmt.Errorf("%w: invalid admin since", ErrInvariantViolation)
	ErrInvalidAdminRoleGranterID  = fmt.Errorf("%w: invalid admin role granter id", ErrInvariantViolation)
	ErrAdminRoleGrantForbidden    = fmt.Errorf("%w: admin role grant forbidden", ErrInvariantViolation)
	ErrAdminSelfRevokeForbidden   = fmt.Errorf("%w: admin cannot revoke own role", ErrInvariantViolation)
	ErrUserAlreadyAdmin           = fmt.Errorf("%w: user is already admin", ErrInvariantViolation)
	ErrUserCustomNameAlreadyEmpty = fmt.Errorf("%w: custom name is already empty", ErrInvariantViolation)
)

// Employment errors
var (
	ErrInvalidEmploymentID                   = fmt.Errorf("%w: invalid employment id", ErrInvariantViolation)
	ErrInvalidEmploymentUserID               = fmt.Errorf("%w: invalid employment user id", ErrInvariantViolation)
	ErrInvalidEmploymentPosition             = fmt.Errorf("%w: invalid employment position", ErrInvariantViolation)
	ErrInvalidEmploymentDeputy               = fmt.Errorf("%w: invalid deputy", ErrInvariantViolation)
	ErrInvalidEmploymentVacation             = fmt.Errorf("%w: invalid vacation state", ErrInvariantViolation)
	ErrInvalidEmploymentAssignedAt           = fmt.Errorf("%w: invalid employment assigned at", ErrInvariantViolation)
	ErrEmploymentVacationTooFarInFuture      = fmt.Errorf("%w: vacation starts too far in the future", ErrInvariantViolation)
	ErrEmploymentVacationAlreadyExists       = fmt.Errorf("%w: vacation already exists", ErrInvariantViolation)
	ErrUserAlreadyEmployed                   = fmt.Errorf("%w: user is already employed", ErrInvariantViolation)
	ErrUserNotEmployed                       = fmt.Errorf("%w: user is not employed", ErrInvariantViolation)
	ErrEmploymentNotFound                    = fmt.Errorf("%w: employment not found", ErrInvariantViolation)
	ErrEmploymentAlreadyExistsInOrganization = fmt.Errorf("%w: employment already exists in organization", ErrInvariantViolation)
)

// Clinic errors
var ErrInvalidClinicID = fmt.Errorf("%w: invalid clinic id", ErrInvariantViolation)

// Department errors
var ErrInvalidDepartmentID = fmt.Errorf("%w: invalid department id", ErrInvariantViolation)

// Organization errors
var (
	ErrInvalidOrganizationID          = fmt.Errorf("%w: invalid organization id", ErrInvariantViolation)
	ErrInvalidOrganizationName        = fmt.Errorf("%w: invalid organization name", ErrInvariantViolation)
	ErrInvalidOrganizationDescription = fmt.Errorf("%w: invalid organization description", ErrInvariantViolation)
)

// Clinic errors (extended)
var (
	ErrInvalidClinicName        = fmt.Errorf("%w: invalid clinic name", ErrInvariantViolation)
	ErrInvalidClinicDescription = fmt.Errorf("%w: invalid clinic description", ErrInvariantViolation)
)

// Department errors (extended)
var (
	ErrInvalidDepartmentName        = fmt.Errorf("%w: invalid department name", ErrInvariantViolation)
	ErrInvalidDepartmentDescription = fmt.Errorf("%w: invalid department description", ErrInvariantViolation)
)

// Address errors
var (
	ErrInvalidAddressValue = fmt.Errorf("%w: invalid address value", ErrInvariantViolation)
	ErrInvalidLatitude     = fmt.Errorf("%w: invalid latitude", ErrInvariantViolation)
	ErrInvalidLongitude    = fmt.Errorf("%w: invalid longitude", ErrInvariantViolation)
)
