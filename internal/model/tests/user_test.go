package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

func ptr(s string) *string {
	return &s
}

// Generate valid snowflake mock ID.
// Snowflake timestamp component must be > 0.
const validUserID = int64(1 << 23)

func TestUserName_Formatting(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		givenName           string
		familyName          string
		middleName          *string
		expectedDisplayName string
		expectedShortName   string
	}{
		{
			name:                "Full name with middle name",
			givenName:           "Ivan",
			familyName:          "Ivanov",
			middleName:          ptr("Ivanovich"),
			expectedDisplayName: "Ivanov Ivan Ivanovich",
			expectedShortName:   "Ivanov I.I.",
		},
		{
			name:                "Name without middle name",
			givenName:           "John",
			familyName:          "Doe",
			middleName:          nil,
			expectedDisplayName: "Doe John",
			expectedShortName:   "Doe J.",
		},
		{
			name:                "Cyrillic support",
			givenName:           "Юлия",
			familyName:          "Смирнова",
			middleName:          ptr("Олеговна"),
			expectedDisplayName: "Смирнова Юлия Олеговна",
			expectedShortName:   "Смирнова Ю.О.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := model.UserName{
				GivenName:  tt.givenName,
				FamilyName: tt.familyName,
				MiddleName: tt.middleName,
			}
			assert.Equal(t, tt.expectedDisplayName, un.DisplayName())
			assert.Equal(t, tt.expectedShortName, un.ShortName())
		})
	}
}

func TestUserName_Equals(t *testing.T) {
	t.Parallel()
	un1 := model.UserName{GivenName: "A", FamilyName: "B", MiddleName: ptr("C")}
	un2 := model.UserName{GivenName: "A", FamilyName: "B", MiddleName: ptr("C")}
	un3 := model.UserName{GivenName: "X", FamilyName: "B", MiddleName: ptr("C")}
	un4 := model.UserName{GivenName: "A", FamilyName: "B", MiddleName: ptr("D")}
	un5 := model.UserName{GivenName: "A", FamilyName: "B", MiddleName: nil}

	assert.True(t, un1.Equals(un2))
	assert.False(t, un1.Equals(un3))
	assert.False(t, un1.Equals(un4))
	assert.False(t, un1.Equals(un5))
	assert.False(t, un5.Equals(un1))

	un6 := model.UserName{GivenName: "A", FamilyName: "B", MiddleName: nil} // Nil vs Nil
	assert.True(t, un5.Equals(un6))
}

func TestNewUserName_Validation(t *testing.T) {
	t.Parallel()
	// Valid
	un, err := model.NewUserName("Ivan", "Ivanov", nil)
	require.NoError(t, err)
	assert.Equal(t, "Ivan", un.GivenName)

	// Invalid Given Name (Empty)
	_, err = model.NewUserName("", "Ivanov", nil)
	assert.ErrorIs(t, err, errors.ErrInvalidGivenName)

	// Invalid Family Name (Empty)
	_, err = model.NewUserName("Ivan", "", nil)
	assert.ErrorIs(t, err, errors.ErrInvalidFamilyName)

	// Invalid Middle Name (Empty str pointer)
	_, err = model.NewUserName("Ivan", "Ivanov", ptr(""))
	assert.ErrorIs(t, err, errors.ErrInvalidMiddleName)

	// Too long names (> 100)
	longName := string(make([]byte, 101))
	_, err = model.NewUserName(longName, "Ivanov", nil)
	assert.ErrorIs(t, err, errors.ErrInvalidGivenName)
}

func TestUser_CreationAndEvents(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("Test", "Testerov", nil)

	// Should create a user and assign UserCreatedEvent
	user, err := model.NewUser(validUserID, *un)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, validUserID, user.ID)
	assert.Equal(t, *un, user.Name)
	assert.Nil(t, user.CustomName)

	events := user.PopEvents()
	require.Len(t, events, 1)
	createdEvent, ok := events[0].(model.UserCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, validUserID, createdEvent.ID)
	assert.Equal(t, *un, createdEvent.Name)

	// Invalid ID validation check
	_, err = model.NewUser(0, *un)
	assert.ErrorIs(t, err, errors.ErrInvalidUserID)

	// Invalid snowflake without timestamp
	_, err = model.NewUser(1<<20, *un)
	assert.ErrorIs(t, err, errors.ErrInvalidUserID)
}

func TestUser_ClearCustomName(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("Base", "Name", nil)
	user, _ := model.NewUser(validUserID, *un)
	user.PopEvents() // Clear initial event

	// Error if already clear
	err := user.ClearCustomName()
	assert.ErrorIs(t, err, errors.ErrCustomNameAlreadyEmpty)

	// Add custom name directly for test
	customName, _ := model.NewUserName("Custom", "Name", nil)
	_ = user.OverrideName(*customName)
	user.PopEvents()

	// Clear successfully
	err = user.ClearCustomName()
	require.NoError(t, err)
	assert.Nil(t, user.CustomName)

	events := user.PopEvents()
	require.Len(t, events, 1)
	customNameEvent, ok := events[0].(model.UserCustomNameUpdatedEvent)
	require.True(t, ok)
	assert.Nil(t, customNameEvent.CustomName)
}

func TestUser_OverrideName(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("Base", "Name", nil)
	user, _ := model.NewUser(validUserID, *un)
	customName, _ := model.NewUserName("Custom", "Name", nil)
	user.PopEvents()

	err := user.OverrideName(*customName)
	require.NoError(t, err)
	assert.NotNil(t, user.CustomName)
	assert.True(t, user.CustomName.Equals(*customName))
	pn := user.PreferredName()
	assert.Equal(t, "Name Custom", (&pn).DisplayName())

	events := user.PopEvents()
	require.Len(t, events, 1)
	assert.IsType(t, model.UserCustomNameUpdatedEvent{}, events[0])

	// Verify idempotency (same name should skip update)
	user.PopEvents()
	err = user.OverrideName(*customName)
	require.NoError(t, err)
	assert.Len(t, user.PopEvents(), 0) // No new events

	// Invalid name check
	invalidCustomName := model.UserName{GivenName: "", FamilyName: "F"}
	err = user.OverrideName(invalidCustomName)
	assert.ErrorIs(t, err, errors.ErrInvalidGivenName)
}

func TestUser_UpdateName(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("First", "Last", nil)
	user, _ := model.NewUser(validUserID, *un)
	newName, _ := model.NewUserName("NewFirst", "NewLast", nil)
	user.PopEvents()

	err := user.UpdateName(*newName)
	require.NoError(t, err)
	assert.True(t, user.Name.Equals(*newName))

	events := user.PopEvents()
	require.Len(t, events, 1)
	assert.IsType(t, model.UserNameUpdatedEvent{}, events[0])

	// Idempotency check
	user.PopEvents()
	err = user.UpdateName(*newName)
	require.NoError(t, err)
	assert.Len(t, user.PopEvents(), 0)

	// Validation
	invalidName := model.UserName{GivenName: "B", FamilyName: ""}
	err = user.UpdateName(invalidName)
	assert.ErrorIs(t, err, errors.ErrInvalidFamilyName)
}
