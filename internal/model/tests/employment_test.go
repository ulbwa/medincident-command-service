package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

func validEmploymentUserID() int64 {
	return int64(1 << 23)
}

func validEmploymentDepartmentID() uuid.UUID {
	return uuid.MustParse("11111111-1111-7111-8111-111111111111")
}

func validEmploymentOrganizationID() uuid.UUID {
	return uuid.MustParse("33333333-3333-7333-8333-333333333333")
}

func validEmploymentClinicID() uuid.UUID {
	return uuid.MustParse("55555555-5555-7555-8555-555555555555")
}

func TestEmployment_DeputyAndVacation(t *testing.T) {
	t.Parallel()

	assignedAt := time.Now().UTC()
	employment, err := model.NewEmployment(validEmploymentUserID(), validEmploymentOrganizationID(), validEmploymentClinicID(), validEmploymentDepartmentID(), nil, assignedAt)
	require.NoError(t, err)
	require.NotNil(t, employment)
	assert.NotEqual(t, uuid.Nil, employment.ID)
	assert.Equal(t, validEmploymentUserID(), employment.UserID)
	assert.Equal(t, validEmploymentOrganizationID(), employment.OrganizationID)
	assert.Equal(t, validEmploymentClinicID(), employment.ClinicID)

	t.Run("AssignAndRemoveDeputy", func(t *testing.T) {
		err := employment.AssignDeputy(int64(2 << 23))
		require.NoError(t, err)
		assert.True(t, employment.HasDeputy())
		require.NotNil(t, employment.Deputy)
		assert.Equal(t, int64(2<<23), employment.Deputy.ID)

		employment.RemoveDeputy()
		assert.False(t, employment.HasDeputy())
		assert.Nil(t, employment.Deputy)
	})

	t.Run("AssignDeputyInvalidIDForbidden", func(t *testing.T) {
		err := employment.AssignDeputy(0)
		assert.ErrorIs(t, err, errs.ErrInvalidEmploymentDeputy)
	})

	t.Run("GrantAndEndVacation", func(t *testing.T) {
		endAt := time.Now().UTC().Add(72 * time.Hour)
		err := employment.GrantVacation(&endAt)
		require.NoError(t, err)
		require.NotNil(t, employment.Vacation)
		assert.True(t, employment.IsOnVacation())
		require.NotNil(t, employment.Vacation.EndsAt)
		assert.Equal(t, endAt, *employment.Vacation.EndsAt)

		employment.EndVacation()
		assert.False(t, employment.IsOnVacation())
		assert.Nil(t, employment.Vacation)
	})

	t.Run("ScheduleVacationDeferred", func(t *testing.T) {
		startsAt := time.Now().UTC().Add(48 * time.Hour)
		endsAt := startsAt.Add(72 * time.Hour)

		err := employment.ScheduleVacation(startsAt, &endsAt)
		require.NoError(t, err)
		require.NotNil(t, employment.Vacation)
		assert.True(t, employment.HasScheduledVacation())
		assert.False(t, employment.IsOnVacation())
	})
}

func TestRestoreEmployment_VacationInvariant(t *testing.T) {
	t.Parallel()

	assignedAt := time.Now().UTC()
	startsAt := assignedAt.Add(48 * time.Hour)
	invalidEndsAt := assignedAt.Add(24 * time.Hour)
	vacation := &model.EmploymentVacation{
		StartsAt: startsAt,
		EndsAt:   &invalidEndsAt,
	}

	_, err := model.RestoreEmployment(
		uuid.MustParse("11111111-1111-7111-8111-111111111111"),
		validEmploymentUserID(),
		validEmploymentOrganizationID(),
		validEmploymentClinicID(),
		validEmploymentDepartmentID(),
		nil,
		assignedAt,
		nil,
		vacation,
	)
	assert.ErrorIs(t, err, errs.ErrInvalidEmploymentVacation)
}

func TestRestoreEmploymentDeputy_Validation(t *testing.T) {
	t.Parallel()

	deputy, err := model.RestoreEmploymentDeputy(int64(2 << 23))
	require.NoError(t, err)
	assert.Equal(t, int64(2<<23), deputy.ID)

	_, err = model.RestoreEmploymentDeputy(0)
	assert.ErrorIs(t, err, errs.ErrInvalidEmploymentDeputy)
}

func TestRestoreEmploymentVacation_CopyEndsAt(t *testing.T) {
	t.Parallel()

	startsAt := time.Now().UTC().Add(24 * time.Hour)
	endsAt := startsAt.Add(48 * time.Hour)

	vacation, err := model.RestoreEmploymentVacation(startsAt, &endsAt)
	require.NoError(t, err)
	require.NotNil(t, vacation.EndsAt)
	assert.Equal(t, endsAt, *vacation.EndsAt)

	originalEndsAt := endsAt
	endsAt = endsAt.Add(24 * time.Hour)
	assert.Equal(t, originalEndsAt, *vacation.EndsAt)
}
