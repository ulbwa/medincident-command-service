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

func validOrgID() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}

func validAddress() model.Address {
	addr, _ := model.NewAddress("Moscow, Red Square, 1", nil)
	return addr
}

func assertInvalidOrganizationField(t *testing.T, err error, field errors.OrganizationField) {
	t.Helper()
	var organizationErr *errors.InvalidOrganizationError
	require.True(t, stderrors.As(err, &organizationErr))
	assert.Equal(t, field, organizationErr.Field)
}

func TestOrganization_NewOrganization(t *testing.T) {
	t.Parallel()

	t.Run("ValidWithoutDescription", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()

		org, err := model.NewOrganization(id, "MedCorp Inc", nil, addr)
		require.NoError(t, err)
		assert.Equal(t, id, org.ID)
		assert.Equal(t, "MedCorp Inc", org.Name)
		assert.Nil(t, org.Description)
		assert.Equal(t, addr, org.LegalAddress)
	})

	t.Run("ValidWithDescription", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()
		desc := "A leading medical organization"

		org, err := model.NewOrganization(id, "MedCorp Inc", &desc, addr)
		require.NoError(t, err)
		assert.NotNil(t, org.Description)
		assert.Equal(t, desc, *org.Description)
	})

	t.Run("NilID", func(t *testing.T) {
		t.Parallel()
		addr := validAddress()

		_, err := model.NewOrganization(uuid.Nil, "MedCorp Inc", nil, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldID)
		var uuidErr *errors.InvalidUUIDError
		require.True(t, stderrors.As(err, &uuidErr))
		assert.Equal(t, errors.UUIDValidationReasonRequired, uuidErr.Reason)
	})

	t.Run("InvalidUUIDVersion", func(t *testing.T) {
		t.Parallel()
		addr := validAddress()
		id := uuid.New() // v4

		_, err := model.NewOrganization(id, "MedCorp Inc", nil, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldID)
		var uuidErr *errors.InvalidUUIDError
		require.True(t, stderrors.As(err, &uuidErr))
		assert.Equal(t, errors.UUIDValidationReasonInvalidVersion, uuidErr.Reason)

		assert.Equal(t, uuid.Version(7), uuidErr.Details.ExpectedVersion)
		assert.Equal(t, uuid.Version(4), uuidErr.Details.ActualVersion)
	})

	t.Run("EmptyName", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()

		_, err := model.NewOrganization(id, "", nil, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldName)
	})

	t.Run("TooLongName", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()
		longName := strings.Repeat("A", 256)

		_, err := model.NewOrganization(id, longName, nil, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldName)
	})

	t.Run("NameWithLeadingWhitespace", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()

		_, err := model.NewOrganization(id, " MedCorp", nil, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldName)
	})

	t.Run("EmptyDescription", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()
		emptyDesc := ""

		_, err := model.NewOrganization(id, "MedCorp Inc", &emptyDesc, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldDescription)
	})

	t.Run("TooLongDescription", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		addr := validAddress()
		longDesc := strings.Repeat("A", 2001)

		_, err := model.NewOrganization(id, "MedCorp Inc", &longDesc, addr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldDescription)
	})

	t.Run("InvalidLegalAddress", func(t *testing.T) {
		t.Parallel()
		id := validOrgID()
		invalidAddr := model.Address{Value: "bad", Point: nil}

		_, err := model.NewOrganization(id, "MedCorp Inc", nil, invalidAddr)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})
}

func TestOrganization_NewAndRestore_AvoidAliasing(t *testing.T) {
	t.Parallel()

	description := "Initial description"
	point := model.GeoPoint{Latitude: 55.7539, Longitude: 37.6208}
	address := model.Address{Value: "Moscow, Red Square, 1", Point: &point}

	created, err := model.NewOrganization(validOrgID(), "Org", &description, address)
	require.NoError(t, err)

	restored, err := model.RestoreOrganization(validOrgID(), "Org", &description, address)
	require.NoError(t, err)

	description = "Mutated outside"
	point.Latitude = 10

	require.NotNil(t, created.Description)
	require.NotNil(t, restored.Description)
	assert.Equal(t, "Initial description", *created.Description)
	assert.Equal(t, "Initial description", *restored.Description)
	require.NotNil(t, created.LegalAddress.Point)
	require.NotNil(t, restored.LegalAddress.Point)
	assert.Equal(t, 55.7539, created.LegalAddress.Point.Latitude)
	assert.Equal(t, 55.7539, restored.LegalAddress.Point.Latitude)
}

func TestOrganization_UpdateName(t *testing.T) {
	t.Parallel()

	id := validOrgID()
	addr := validAddress()
	org, _ := model.NewOrganization(id, "OldName", nil, addr)

	t.Run("UpdateToNewName", func(t *testing.T) {
		err := org.UpdateName("NewName")
		require.NoError(t, err)
		assert.Equal(t, "NewName", org.Name)
	})

	t.Run("Idempotent", func(t *testing.T) {
		err := org.UpdateName("NewName")
		require.NoError(t, err)
		assert.Equal(t, "NewName", org.Name)
	})

	t.Run("InvalidName", func(t *testing.T) {
		err := org.UpdateName("")
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldName)
		var tooShortErr *errors.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
		assert.Equal(t, "NewName", org.Name) // unchanged
	})
}

func TestOrganization_SetDescription(t *testing.T) {
	t.Parallel()

	id := validOrgID()
	addr := validAddress()
	org, _ := model.NewOrganization(id, "Org", nil, addr)

	t.Run("SetDescription", func(t *testing.T) {
		desc := "A great organization"
		err := org.SetDescription(desc)
		require.NoError(t, err)
		assert.NotNil(t, org.Description)
		assert.Equal(t, desc, *org.Description)
	})

	t.Run("ClearDescription", func(t *testing.T) {
		org.RemoveDescription()
		assert.Nil(t, org.Description)
	})

	t.Run("IdempotentRemove", func(t *testing.T) {
		org.RemoveDescription()
		assert.Nil(t, org.Description)
	})

	t.Run("InvalidDescription", func(t *testing.T) {
		err := org.SetDescription("")
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldDescription)
		var tooShortErr *errors.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
	})
}

func TestOrganization_UpdateLegalAddress(t *testing.T) {
	t.Parallel()

	id := validOrgID()
	addr1 := validAddress()
	org, _ := model.NewOrganization(id, "Org", nil, addr1)

	t.Run("UpdateAddress", func(t *testing.T) {
		addr2, _ := model.NewAddress("Saint Petersburg, Palace Square, 1", nil)
		err := org.UpdateLegalAddress(addr2)
		require.NoError(t, err)
		assert.Equal(t, addr2, org.LegalAddress)
	})

	t.Run("Idempotent", func(t *testing.T) {
		currentAddr := org.LegalAddress
		err := org.UpdateLegalAddress(currentAddr)
		require.NoError(t, err)
		assert.Equal(t, currentAddr, org.LegalAddress)
	})

	t.Run("InvalidAddress", func(t *testing.T) {
		invalidAddr := model.Address{Value: "bad", Point: nil}
		previousAddr := org.LegalAddress

		err := org.UpdateLegalAddress(invalidAddr)
		assertInvalidOrganizationField(t, err, errors.OrganizationFieldLegalAddress)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
		assert.Equal(t, previousAddr, org.LegalAddress)
	})
}

func TestOrganization_RestoreOrganization(t *testing.T) {
	t.Parallel()

	id := validOrgID()
	addr := validAddress()
	desc := "Restored org"

	org, err := model.RestoreOrganization(id, "RestoredOrg", &desc, addr)
	require.NoError(t, err)
	assert.Equal(t, id, org.ID)
	assert.Equal(t, "RestoredOrg", org.Name)
	assert.NotNil(t, org.Description)
	assert.Equal(t, desc, *org.Description)

	t.Run("InvalidLegalAddress", func(t *testing.T) {
		t.Parallel()
		invalidAddr := model.Address{Value: "bad", Point: nil}

		_, err := model.RestoreOrganization(validOrgID(), "RestoredOrg", nil, invalidAddr)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})
}
