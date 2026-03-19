package model

import (
	stderrors "errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func TestValidateUserEmployments_DuplicateOrganization(t *testing.T) {
	t.Parallel()

	organizationID := uuid.Must(uuid.NewV7())
	clinicID1 := uuid.Must(uuid.NewV7())
	clinicID2 := uuid.Must(uuid.NewV7())
	departmentID1 := uuid.Must(uuid.NewV7())
	departmentID2 := uuid.Must(uuid.NewV7())
	assignedAt := time.Now().UTC()

	employment1, err := RestoreEmployment(
		uuid.Must(uuid.NewV7()),
		int64(1<<23),
		organizationID,
		clinicID1,
		departmentID1,
		nil,
		assignedAt,
		nil,
		nil,
	)
	require.NoError(t, err)

	employment2, err := RestoreEmployment(
		uuid.Must(uuid.NewV7()),
		int64(1<<23),
		organizationID,
		clinicID2,
		departmentID2,
		nil,
		assignedAt,
		nil,
		nil,
	)
	require.NoError(t, err)

	u := &User{
		ID:          int64(1 << 23),
		Employments: []*Employment{employment1, employment2},
	}

	err = validateUserEmployments(u)
	require.Error(t, err)

	var invalidUserErr *errs.InvalidUserError
	require.True(t, stderrors.As(err, &invalidUserErr))
	require.Equal(t, errs.UserFieldEmployments, invalidUserErr.Field)

	var itemErr *errs.InvalidCollectionItemError
	require.True(t, stderrors.As(invalidUserErr.Reason, &itemErr))
	require.Equal(t, 1, itemErr.Index)

	var invalidEmploymentErr *errs.InvalidEmploymentError
	require.True(t, stderrors.As(itemErr.Reason, &invalidEmploymentErr))
	require.Equal(t, errs.EmploymentFieldOrganizationID, invalidEmploymentErr.Field)

	var duplicateErr *errs.ValueDuplicateError
	require.True(t, stderrors.As(invalidEmploymentErr.Reason, &duplicateErr))
	require.Equal(t, organizationID, duplicateErr.Value)
}
