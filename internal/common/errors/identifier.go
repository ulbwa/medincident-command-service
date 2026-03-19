package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type UUIDValidationReason string

const (
	UUIDValidationReasonRequired       UUIDValidationReason = "required"
	UUIDValidationReasonInvalidVariant UUIDValidationReason = "invalid_variant"
	UUIDValidationReasonInvalidVersion UUIDValidationReason = "invalid_version"
)

type UUIDValidationDetails struct {
	ExpectedVersion uuid.Version
	ActualVersion   uuid.Version
	ExpectedVariant uuid.Variant
	ActualVariant   uuid.Variant
	Value           string
}

type InvalidUUIDError struct {
	Reason  UUIDValidationReason
	Details UUIDValidationDetails
}

func (e *InvalidUUIDError) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch e.Reason {
	case UUIDValidationReasonRequired:
		return "invalid uuid: value is required"
	case UUIDValidationReasonInvalidVariant:
		return fmt.Sprintf("invalid uuid: expected variant %d, got %d", e.Details.ExpectedVariant, e.Details.ActualVariant)
	case UUIDValidationReasonInvalidVersion:
		return fmt.Sprintf("invalid uuid: expected version %d, got %d", e.Details.ExpectedVersion, e.Details.ActualVersion)
	default:
		return "invalid uuid"
	}
}

func NewUUIDRequiredError(value uuid.UUID) *InvalidUUIDError {
	return &InvalidUUIDError{
		Reason: UUIDValidationReasonRequired,
		Details: UUIDValidationDetails{
			Value: value.String(),
		},
	}
}

func NewUUIDInvalidVariantError(value uuid.UUID, expected, actual uuid.Variant) *InvalidUUIDError {
	return &InvalidUUIDError{
		Reason: UUIDValidationReasonInvalidVariant,
		Details: UUIDValidationDetails{
			ExpectedVariant: expected,
			ActualVariant:   actual,
			ActualVersion:   value.Version(),
			Value:           value.String(),
		},
	}
}

func NewUUIDInvalidVersionError(value uuid.UUID, expected, actual uuid.Version) *InvalidUUIDError {
	return &InvalidUUIDError{
		Reason: UUIDValidationReasonInvalidVersion,
		Details: UUIDValidationDetails{
			ExpectedVersion: expected,
			ActualVersion:   actual,
			ActualVariant:   value.Variant(),
			Value:           value.String(),
		},
	}
}

type SnowflakeValidationReason string

const (
	SnowflakeValidationReasonMustBePositive       SnowflakeValidationReason = "must_be_positive"
	SnowflakeValidationReasonInvalidTimestampPart SnowflakeValidationReason = "invalid_timestamp_component"
)

type SnowflakeValidationDetails struct {
	ActualValue        int64
	MinValueExclusive  int64
	TimestampComponent int64
}

type InvalidSnowflakeIDError struct {
	Reason  SnowflakeValidationReason
	Details SnowflakeValidationDetails
}

func (e *InvalidSnowflakeIDError) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch e.Reason {
	case SnowflakeValidationReasonMustBePositive:
		return fmt.Sprintf("invalid snowflake id: expected value > %d, got %d", e.Details.MinValueExclusive, e.Details.ActualValue)
	case SnowflakeValidationReasonInvalidTimestampPart:
		return fmt.Sprintf("invalid snowflake id: timestamp component must be > 0, got %d", e.Details.TimestampComponent)
	default:
		return "invalid snowflake id"
	}
}

func NewSnowflakeMustBePositiveError(actualValue int64) *InvalidSnowflakeIDError {
	return &InvalidSnowflakeIDError{
		Reason: SnowflakeValidationReasonMustBePositive,
		Details: SnowflakeValidationDetails{
			ActualValue:       actualValue,
			MinValueExclusive: 0,
		},
	}
}

func NewSnowflakeInvalidTimestampComponentError(actualValue, timestampComponent int64) *InvalidSnowflakeIDError {
	return &InvalidSnowflakeIDError{
		Reason: SnowflakeValidationReasonInvalidTimestampPart,
		Details: SnowflakeValidationDetails{
			ActualValue:        actualValue,
			MinValueExclusive:  0,
			TimestampComponent: timestampComponent,
		},
	}
}
