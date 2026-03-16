package errors

import "errors"

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrIdentityNotFound  = errors.New("identity not found")
	ErrIdentityNotHuman  = errors.New("identity profile is not a human")
	ErrUserAlreadyExists = errors.New("user already exists")
)
