package tests

import (
	stderrors "errors"
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

func assertInvalidClinicField(t *testing.T, err error, field errors.ClinicField) {
	t.Helper()
	var clinicErr *errors.InvalidClinicError
	require.True(t, stderrors.As(err, &clinicErr))
	assert.Equal(t, field, clinicErr.Field)
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
		assertInvalidClinicField(t, err, errors.ClinicFieldID)
		var uuidErr *errors.InvalidUUIDError
		require.True(t, stderrors.As(err, &uuidErr))
		assert.Equal(t, errors.UUIDValidationReasonRequired, uuidErr.Reason)
	})

	t.Run("InvalidUUIDVersion", func(t *testing.T) {
		t.Parallel()
		orgID := validOrgID()
		addr := validClinicAddress()
		id := uuid.New() // v4

		_, err := model.NewClinic(id, orgID, "Clinic", nil, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldID)
		var uuidErr *errors.InvalidUUIDError
		require.True(t, stderrors.As(err, &uuidErr))
		assert.Equal(t, errors.UUIDValidationReasonInvalidVersion, uuidErr.Reason)
	})

	t.Run("NilOrganizationID", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		addr := validClinicAddress()

		_, err := model.NewClinic(id, uuid.Nil, "Clinic", nil, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldOrganizationID)
	})

	t.Run("InvalidOrganizationUUIDVersion", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		addr := validClinicAddress()
		orgID := uuid.New() // v4

		_, err := model.NewClinic(id, orgID, "Clinic", nil, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldOrganizationID)
	})

	t.Run("EmptyName", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()

		_, err := model.NewClinic(id, orgID, "", nil, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldName)
	})

	t.Run("TooLongName", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		longName := strings.Repeat("A", 256)

		_, err := model.NewClinic(id, orgID, longName, nil, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldName)
	})

	t.Run("NameWithTrailingWhitespace", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()

		_, err := model.NewClinic(id, orgID, "Clinic ", nil, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldName)
	})

	t.Run("EmptyDescription", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		emptyDesc := ""

		_, err := model.NewClinic(id, orgID, "Clinic", &emptyDesc, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldDescription)
	})

	t.Run("TooLongDescription", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		addr := validClinicAddress()
		longDesc := strings.Repeat("A", 2001)

		_, err := model.NewClinic(id, orgID, "Clinic", &longDesc, addr)
		assertInvalidClinicField(t, err, errors.ClinicFieldDescription)
	})

	t.Run("InvalidPhysicalAddress", func(t *testing.T) {
		t.Parallel()
		id := validClinicID()
		orgID := validOrgID()
		invalidAddr := model.Address{Value: "bad", Point: nil}

		_, err := model.NewClinic(id, orgID, "Clinic", nil, invalidAddr)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})
}

func TestClinic_NewAndRestore_AvoidAliasing(t *testing.T) {
	t.Parallel()

	description := "Initial clinic description"
	point := model.GeoPoint{Latitude: 55.7539, Longitude: 37.6208}
	address := model.Address{Value: "Moscow, Clinic Street, 10", Point: &point}

	created, err := model.NewClinic(validClinicID(), validOrgID(), "Clinic", &description, address)
	require.NoError(t, err)

	restored, err := model.RestoreClinic(validClinicID(), validOrgID(), "Clinic", &description, address)
	require.NoError(t, err)

	description = "Mutated outside"
	point.Longitude = 10

	require.NotNil(t, created.Description)
	require.NotNil(t, restored.Description)
	assert.Equal(t, "Initial clinic description", *created.Description)
	assert.Equal(t, "Initial clinic description", *restored.Description)
	require.NotNil(t, created.PhysicalAddress.Point)
	require.NotNil(t, restored.PhysicalAddress.Point)
	assert.Equal(t, 37.6208, created.PhysicalAddress.Point.Longitude)
	assert.Equal(t, 37.6208, restored.PhysicalAddress.Point.Longitude)
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
		assertInvalidClinicField(t, err, errors.ClinicFieldName)
		var tooShortErr *errors.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
		assert.Equal(t, "NewClinic", clinic.Name) // unchanged
	})
}

func TestClinic_SetDescription(t *testing.T) {
	t.Parallel()

	id := validClinicID()
	orgID := validOrgID()
	addr := validClinicAddress()
	clinic, _ := model.NewClinic(id, orgID, "Clinic", nil, addr)

	t.Run("SetDescription", func(t *testing.T) {
		desc := "A modern clinic"
		err := clinic.SetDescription(desc)
		require.NoError(t, err)
		assert.NotNil(t, clinic.Description)
		assert.Equal(t, desc, *clinic.Description)
	})

	t.Run("ClearDescription", func(t *testing.T) {
		clinic.RemoveDescription()
		assert.Nil(t, clinic.Description)
	})

	t.Run("InvalidDescription", func(t *testing.T) {
		err := clinic.SetDescription("")
		assertInvalidClinicField(t, err, errors.ClinicFieldDescription)
		var tooShortErr *errors.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
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
		assertInvalidClinicField(t, err, errors.ClinicFieldPhysicalAddress)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
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
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})
}
