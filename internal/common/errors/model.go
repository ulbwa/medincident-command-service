package errors

import (
	"errors"
	"fmt"
)

// ErrInvariantViolation indicates that a domain rule or invariant was broken.
var ErrInvariantViolation = errors.New("invariant violation")

// User & Identity errors
var (
	ErrInvalidUserID           = fmt.Errorf("%w: invalid user id", ErrInvariantViolation)
	ErrInvalidGivenName        = fmt.Errorf("%w: invalid given name", ErrInvariantViolation)
	ErrInvalidFamilyName       = fmt.Errorf("%w: invalid family name", ErrInvariantViolation)
	ErrInvalidMiddleName       = fmt.Errorf("%w: invalid middle name", ErrInvariantViolation)
	ErrInvalidAdminSince       = fmt.Errorf("%w: invalid admin since", ErrInvariantViolation)
	ErrAdminRoleGrantForbidden = fmt.Errorf("%w: admin role grant forbidden", ErrInvariantViolation)
	ErrCustomNameAlreadyEmpty  = fmt.Errorf("%w: custom name is already empty", ErrInvariantViolation)
)

// Employment errors
var (
	ErrInvalidEmploymentID                   = fmt.Errorf("%w: invalid employment id", ErrInvariantViolation)
	ErrInvalidEmploymentPosition             = fmt.Errorf("%w: invalid employment position", ErrInvariantViolation)
	ErrInvalidEmploymentDeputy               = fmt.Errorf("%w: invalid deputy", ErrInvariantViolation)
	ErrInvalidEmploymentVacation             = fmt.Errorf("%w: invalid vacation state", ErrInvariantViolation)
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
