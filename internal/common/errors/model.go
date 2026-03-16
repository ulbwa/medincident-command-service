package errors

import (
	"errors"
	"fmt"
)

// ErrInvariantViolation indicates that a domain rule or invariant was broken.
var ErrInvariantViolation = errors.New("invariant violation")

// User & Identity errors
var (
	ErrInvalidUserID          = fmt.Errorf("%w: invalid user id", ErrInvariantViolation)
	ErrInvalidGivenName       = fmt.Errorf("%w: invalid given name", ErrInvariantViolation)
	ErrInvalidFamilyName      = fmt.Errorf("%w: invalid family name", ErrInvariantViolation)
	ErrInvalidMiddleName      = fmt.Errorf("%w: invalid middle name", ErrInvariantViolation)
	ErrCustomNameAlreadyEmpty = fmt.Errorf("%w: custom name is already empty", ErrInvariantViolation)
)

// Employment errors
var (
	ErrInvalidEmploymentPosition = fmt.Errorf("%w: invalid employment position", ErrInvariantViolation)
	ErrInvalidEmploymentDeputy   = fmt.Errorf("%w: invalid deputy", ErrInvariantViolation)
	ErrInvalidEmploymentVacation = fmt.Errorf("%w: invalid vacation state", ErrInvariantViolation)
)

// Clinic errors
var ErrInvalidClinicID = fmt.Errorf("%w: invalid clinic id", ErrInvariantViolation)

// Department errors
var ErrInvalidDepartmentID = fmt.Errorf("%w: invalid department id", ErrInvariantViolation)

// Organization errors
var ErrInvalidOrganizationID = fmt.Errorf("%w: invalid organization id", ErrInvariantViolation)
