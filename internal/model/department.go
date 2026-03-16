package model

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func validateDepartmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidDepartmentID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidDepartmentID)
	}
	return nil
}
