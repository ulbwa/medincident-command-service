package model

import (
	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func validateUUID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errs.NewUUIDRequiredError(id)
	}
	if id.Variant() != uuid.RFC4122 {
		return errs.NewUUIDInvalidVariantError(id, uuid.RFC4122, id.Variant())
	}
	return nil
}

func validateUUIDVersion(id uuid.UUID, version uuid.Version) error {
	if err := validateUUID(id); err != nil {
		return err
	}
	if id.Version() != version {
		return errs.NewUUIDInvalidVersionError(id, version, id.Version())
	}
	return nil
}

func validateUUIDv7(id uuid.UUID) error {
	return validateUUIDVersion(id, 7)
}

func validateSnowflakeID(id int64) error {
	if id <= 0 {
		return errs.NewSnowflakeMustBePositiveError(id)
	}
	timestampComponent := id >> 22
	if timestampComponent <= 0 {
		return errs.NewSnowflakeInvalidTimestampComponentError(id, timestampComponent)
	}
	return nil
}
