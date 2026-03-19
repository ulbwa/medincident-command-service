package errors

import "errors"

var (
	ErrAdminRoleGrantActorNotAdmin           = errors.New("admin role grant forbidden: actor is not admin")
	ErrAdminRoleGrantInsufficientTenure      = errors.New("admin role grant forbidden: actor has insufficient tenure")
	ErrUserAlreadyAdmin                      = errors.New("user is already admin")
	ErrUserNotAdmin                          = errors.New("user is not admin")
	ErrAdminSelfRevokeForbidden              = errors.New("admin cannot revoke own role")
	ErrEmploymentAlreadyExistsInOrganization = errors.New("employment already exists in organization")
	ErrEmploymentNotFound                    = errors.New("employment not found")
)
