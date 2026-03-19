package model

import (
	stderrors "errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func TestValidateUUID_Generic(t *testing.T) {
	t.Parallel()

	t.Run("NilUUID", func(t *testing.T) {
		err := validateUUID(uuid.Nil)
		require.Error(t, err)

		var uuidErr *errs.InvalidUUIDError
		require.True(t, stderrors.As(err, &uuidErr))
		assert.Equal(t, errs.UUIDValidationReasonRequired, uuidErr.Reason)
		assert.Equal(t, uuid.Nil.String(), uuidErr.Details.Value)
	})

	t.Run("ValidRFC4122", func(t *testing.T) {
		id := uuid.New()
		err := validateUUID(id)
		require.NoError(t, err)
	})
}

func TestValidateUUIDVersion(t *testing.T) {
	t.Parallel()

	v4 := uuid.New()
	err := validateUUIDVersion(v4, 7)
	require.Error(t, err)

	var uuidErr *errs.InvalidUUIDError
	require.True(t, stderrors.As(err, &uuidErr))
	assert.Equal(t, errs.UUIDValidationReasonInvalidVersion, uuidErr.Reason)
	assert.Equal(t, uuid.Version(7), uuidErr.Details.ExpectedVersion)
	assert.Equal(t, uuid.Version(4), uuidErr.Details.ActualVersion)
	assert.Equal(t, v4.String(), uuidErr.Details.Value)
}

func TestValidateSnowflakeID(t *testing.T) {
	t.Parallel()

	t.Run("NonPositive", func(t *testing.T) {
		err := validateSnowflakeID(0)
		require.Error(t, err)

		var snowflakeErr *errs.InvalidSnowflakeIDError
		require.True(t, stderrors.As(err, &snowflakeErr))
		assert.Equal(t, errs.SnowflakeValidationReasonMustBePositive, snowflakeErr.Reason)
		assert.Equal(t, int64(0), snowflakeErr.Details.ActualValue)
	})

	t.Run("InvalidTimestampComponent", func(t *testing.T) {
		value := int64(1 << 20)
		err := validateSnowflakeID(value)
		require.Error(t, err)

		var snowflakeErr *errs.InvalidSnowflakeIDError
		require.True(t, stderrors.As(err, &snowflakeErr))
		assert.Equal(t, errs.SnowflakeValidationReasonInvalidTimestampPart, snowflakeErr.Reason)
		assert.Equal(t, int64(0), snowflakeErr.Details.TimestampComponent)
		assert.Equal(t, value, snowflakeErr.Details.ActualValue)
	})

	t.Run("Valid", func(t *testing.T) {
		err := validateSnowflakeID(int64(1 << 23))
		require.NoError(t, err)
	})
}
