package tests

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

func validDeptID() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}

func TestDepartment_NewDepartment(t *testing.T) {
	t.Parallel()

	t.Run("ValidWithoutDescription", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()

		dept, err := model.NewDepartment(id, clinicID, "Cardiology", nil)
		require.NoError(t, err)
		assert.Equal(t, id, dept.ID)
		assert.Equal(t, clinicID, dept.ClinicID)
		assert.Equal(t, "Cardiology", dept.Name)
		assert.Nil(t, dept.Description)
	})

	t.Run("ValidWithDescription", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()
		desc := "Cardiology department for heart diseases"

		dept, err := model.NewDepartment(id, clinicID, "Cardiology", &desc)
		require.NoError(t, err)
		assert.NotNil(t, dept.Description)
		assert.Equal(t, desc, *dept.Description)
	})

	t.Run("NilID", func(t *testing.T) {
		t.Parallel()
		clinicID := validClinicID()

		_, err := model.NewDepartment(uuid.Nil, clinicID, "Department", nil)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentID)
	})

	t.Run("InvalidUUIDVersion", func(t *testing.T) {
		t.Parallel()
		clinicID := validClinicID()
		id := uuid.New() // v4

		_, err := model.NewDepartment(id, clinicID, "Department", nil)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentID)
	})

	t.Run("NilClinicID", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()

		_, err := model.NewDepartment(id, uuid.Nil, "Department", nil)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicID)
	})

	t.Run("InvalidClinicUUIDVersion", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := uuid.New() // v4

		_, err := model.NewDepartment(id, clinicID, "Department", nil)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicID)
	})

	t.Run("EmptyName", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()

		_, err := model.NewDepartment(id, clinicID, "", nil)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentName)
	})

	t.Run("TooLongName", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()
		longName := strings.Repeat("A", 256)

		_, err := model.NewDepartment(id, clinicID, longName, nil)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentName)
	})

	t.Run("NameWithLeadingWhitespace", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()

		_, err := model.NewDepartment(id, clinicID, " Department", nil)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentName)
	})

	t.Run("EmptyDescription", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()
		emptyDesc := ""

		_, err := model.NewDepartment(id, clinicID, "Department", &emptyDesc)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentDescription)
	})

	t.Run("TooLongDescription", func(t *testing.T) {
		t.Parallel()
		id := validDeptID()
		clinicID := validClinicID()
		longDesc := strings.Repeat("A", 2001)

		_, err := model.NewDepartment(id, clinicID, "Department", &longDesc)
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentDescription)
	})
}

func TestDepartment_UpdateName(t *testing.T) {
	t.Parallel()

	id := validDeptID()
	clinicID := validClinicID()
	dept, _ := model.NewDepartment(id, clinicID, "OldDept", nil)

	t.Run("UpdateToNewName", func(t *testing.T) {
		err := dept.UpdateName("NewDept")
		require.NoError(t, err)
		assert.Equal(t, "NewDept", dept.Name)
	})

	t.Run("Idempotent", func(t *testing.T) {
		err := dept.UpdateName("NewDept")
		require.NoError(t, err)
		assert.Equal(t, "NewDept", dept.Name)
	})

	t.Run("InvalidName", func(t *testing.T) {
		err := dept.UpdateName("")
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentName)
		assert.Equal(t, "NewDept", dept.Name) // unchanged
	})
}

func TestDepartment_UpdateDescription(t *testing.T) {
	t.Parallel()

	id := validDeptID()
	clinicID := validClinicID()
	dept, _ := model.NewDepartment(id, clinicID, "Dept", nil)

	t.Run("SetDescription", func(t *testing.T) {
		desc := "A specialized department"
		err := dept.UpdateDescription(desc)
		require.NoError(t, err)
		assert.NotNil(t, dept.Description)
		assert.Equal(t, desc, *dept.Description)
	})

	t.Run("ClearDescription", func(t *testing.T) {
		dept.RemoveDescription()
		assert.Nil(t, dept.Description)
	})

	t.Run("InvalidDescription", func(t *testing.T) {
		err := dept.UpdateDescription("")
		assert.ErrorIs(t, err, errors.ErrInvalidDepartmentDescription)
	})
}

func TestDepartment_RestoreDepartment(t *testing.T) {
	t.Parallel()

	id := validDeptID()
	clinicID := validClinicID()
	desc := "Restored department"

	dept, err := model.RestoreDepartment(id, clinicID, "RestoredDept", &desc)
	require.NoError(t, err)
	assert.Equal(t, id, dept.ID)
	assert.Equal(t, clinicID, dept.ClinicID)
	assert.Equal(t, "RestoredDept", dept.Name)
	assert.NotNil(t, dept.Description)
	assert.Equal(t, desc, *dept.Description)
}
