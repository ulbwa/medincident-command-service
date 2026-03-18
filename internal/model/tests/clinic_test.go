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

func validClinicID() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}

func validClinicAddress() model.Address {
	addr, _ := model.NewAddress("Moscow, Clinic Street, 10", nil)
	return addr
}

func TestClinic_NewClinic(t *testing.T) {
	t.Parallel()

	t.Run("ValidWithoutDescription", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()

		clinic, err := model.NewClinic(id, orgID, "Central Clinic", nil, addr)
		require.NoError(t, err)
		assert.Equal(t, id, clinic.ID)
		assert.Equal(t, orgID, clinic.OrganizationID)
		assert.Equal(t, "Central Clinic", clinic.Name)
		assert.Nil(t, clinic.Description)
		assert.Equal(t, addr, clinic.PhysicalAddress)
	})

	t.Run("ValidWithDescription", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		desc := "Main clinic in Moscow"

		clinic, err := model.NewClinic(id, orgID, "Central Clinic", &desc, addr)
		require.NoError(t, err)
		assert.NotNil(t, clinic.Description)
		assert.Equal(t, desc, *clinic.Description)
	})

	t.Run("NilID", func(t *testing.T) {
		t.Parallel()
		orgID := validOrgID()
		addr := validClinicAddress()

		_, err := model.NewClinic(uuid.Nil, orgID, "Clinic", nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicID)
	})

	t.Run("InvalidUUIDVersion", func(t *testing.T) {
		t.Parallel()
		orgID := validOrgID()
		addr := validClinicAddress()
		id := uuid.New() // v4

		_, err := model.NewClinic(id, orgID, "Clinic", nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicID)
	})

	t.Run("NilOrganizationID", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		addr := validClinicAddress()

		_, err := model.NewClinic(id, uuid.Nil, "Clinic", nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidOrganizationID)
	})

	t.Run("InvalidOrganizationUUIDVersion", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		addr := validClinicAddress()
		orgID := uuid.New() // v4

		_, err := model.NewClinic(id, orgID, "Clinic", nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidOrganizationID)
	})

	t.Run("EmptyName", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()

		_, err := model.NewClinic(id, orgID, "", nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicName)
	})

	t.Run("TooLongName", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		longName := strings.Repeat("A", 256)

		_, err := model.NewClinic(id, orgID, longName, nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicName)
	})

	t.Run("NameWithTrailingWhitespace", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()

		_, err := model.NewClinic(id, orgID, "Clinic ", nil, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicName)
	})

	t.Run("EmptyDescription", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		emptyDesc := ""

		_, err := model.NewClinic(id, orgID, "Clinic", &emptyDesc, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicDescription)
	})

	t.Run("TooLongDescription", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		longDesc := strings.Repeat("A", 2001)

		_, err := model.NewClinic(id, orgID, "Clinic", &longDesc, addr)
		assert.ErrorIs(t, err, errors.ErrInvalidClinicDescription)
	})

	t.Run("InvalidPhysicalAddress", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		invalidAddr := model.Address{Value: "bad", Point: nil}

		_, err := model.NewClinic(id, orgID, "Clinic", nil, invalidAddr)
		assert.ErrorIs(t, err, errors.ErrInvalidAddressValue)
	})
}

func TestClinic_UpdateName(t *testing.T) {
	t.Parallel()

	id := validClinicID()
	orgID := validOrgID()
	addr := validClinicAddress()
	clinic, _ := model.NewClinic(id, orgID, "OldClinic", nil, addr)

	t.Run("UpdateToNewName", func(t *testing.T) {
		err := clinic.UpdateName("NewClinic")
		require.NoError(t, err)
		assert.Equal(t, "NewClinic", clinic.Name)
	})

	t.Run("Idempotent", func(t *testing.T) {
		err := clinic.UpdateName("NewClinic")
		require.NoError(t, err)
		assert.Equal(t, "NewClinic", clinic.Name)
	})

	t.Run("InvalidName", func(t *testing.T) {
		err := clinic.UpdateName("")
		assert.ErrorIs(t, err, errors.ErrInvalidClinicName)
		assert.Equal(t, "NewClinic", clinic.Name) // unchanged
	})
}

func TestClinic_UpdateDescription(t *testing.T) {
	t.Parallel()

	id := validClinicID()
	orgID := validOrgID()
	addr := validClinicAddress()
	clinic, _ := model.NewClinic(id, orgID, "Clinic", nil, addr)

	t.Run("SetDescription", func(t *testing.T) {
		desc := "A modern clinic"
		err := clinic.UpdateDescription(desc)
		require.NoError(t, err)
		assert.NotNil(t, clinic.Description)
		assert.Equal(t, desc, *clinic.Description)
	})

	t.Run("ClearDescription", func(t *testing.T) {
		clinic.RemoveDescription()
		assert.Nil(t, clinic.Description)
	})

	t.Run("InvalidDescription", func(t *testing.T) {
		err := clinic.UpdateDescription("")
		assert.ErrorIs(t, err, errors.ErrInvalidClinicDescription)
	})
}

func TestClinic_UpdatePhysicalAddress(t *testing.T) {
	t.Parallel()

	id := validClinicID()
	orgID := validOrgID()
	addr1 := validClinicAddress()
	clinic, _ := model.NewClinic(id, orgID, "Clinic", nil, addr1)

	t.Run("UpdateAddress", func(t *testing.T) {
		addr2, _ := model.NewAddress("Saint Petersburg, Hospital Lane, 5", nil)
		err := clinic.UpdatePhysicalAddress(addr2)
		require.NoError(t, err)
		assert.Equal(t, addr2, clinic.PhysicalAddress)
	})

	t.Run("Idempotent", func(t *testing.T) {
		currentAddr := clinic.PhysicalAddress
		err := clinic.UpdatePhysicalAddress(currentAddr)
		require.NoError(t, err)
		assert.Equal(t, currentAddr, clinic.PhysicalAddress)
	})

	t.Run("InvalidAddress", func(t *testing.T) {
		invalidAddr := model.Address{Value: "bad", Point: nil}
		previousAddr := clinic.PhysicalAddress

		err := clinic.UpdatePhysicalAddress(invalidAddr)
		assert.ErrorIs(t, err, errors.ErrInvalidAddressValue)
		assert.Equal(t, previousAddr, clinic.PhysicalAddress)
	})
}

func TestClinic_RestoreClinic(t *testing.T) {
	t.Parallel()

	id := validClinicID()
	orgID := validOrgID()
	addr := validClinicAddress()
	desc := "Restored clinic"

	clinic, err := model.RestoreClinic(id, orgID, "RestoredClinic", &desc, addr)
	require.NoError(t, err)
	assert.Equal(t, id, clinic.ID)
	assert.Equal(t, orgID, clinic.OrganizationID)
	assert.Equal(t, "RestoredClinic", clinic.Name)
	assert.NotNil(t, clinic.Description)
	assert.Equal(t, desc, *clinic.Description)

	t.Run("InvalidPhysicalAddress", func(t *testing.T) {
		t.Parallel()
		invalidAddr := model.Address{Value: "bad", Point: nil}

		_, err := model.RestoreClinic(validClinicID(), validOrgID(), "RestoredClinic", nil, invalidAddr)
		assert.ErrorIs(t, err, errors.ErrInvalidAddressValue)
	})
}
