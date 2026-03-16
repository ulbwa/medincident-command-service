package model

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

// nolint:unused
func validateOrganizationID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: required", errs.ErrInvalidOrganizationID)
	}
	if id.Version() != 7 {
		return fmt.Errorf("%w: must be a UUIDv7", errs.ErrInvalidOrganizationID)
	}
	return nil
}
