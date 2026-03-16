package model

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

// nolint:unused
func validateClinicID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidClinicID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidClinicID)
	}
	return nil
}
