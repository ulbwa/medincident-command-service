package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
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
	assert.ErrorIs(t, err, errors.ErrInvalidUserGivenName)

	// Invalid Family Name (Empty)
	_, err = model.NewUserName("Ivan", "", nil)
	assert.ErrorIs(t, err, errors.ErrInvalidUserFamilyName)

	// Invalid Middle Name (Empty str pointer)
	_, err = model.NewUserName("Ivan", "Ivanov", ptr(""))
	assert.ErrorIs(t, err, errors.ErrInvalidUserMiddleName)

	// Too long names (> 100)
	longName := string(make([]byte, 101))
	_, err = model.NewUserName(longName, "Ivanov", nil)
	assert.ErrorIs(t, err, errors.ErrInvalidUserGivenName)
}

func TestUser_CreationAndEvents(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("Test", "Testerov", nil)

	// Should create a user and assign UserCreatedEvent
	user, err := model.NewUser(validUserID, un)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, validUserID, user.ID)
	assert.Equal(t, un, user.Name)
	assert.Nil(t, user.CustomName)

	events := user.PopEvents()
	require.Len(t, events, 1)
	createdEvent, ok := events[0].(model.UserCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, validUserID, createdEvent.ID)
	assert.Equal(t, un, createdEvent.Name)

	// Invalid ID validation check
	_, err = model.NewUser(0, un)
	assert.ErrorIs(t, err, errors.ErrInvalidUserID)

	// Invalid snowflake without timestamp
	_, err = model.NewUser(1<<20, un)
	assert.ErrorIs(t, err, errors.ErrInvalidUserID)
}

func TestUser_Name_AvoidAliasing(t *testing.T) {
	t.Parallel()

	middleName := "Ivanovich"
	name := model.UserName{
		GivenName:  "Ivan",
		FamilyName: "Ivanov",
		MiddleName: &middleName,
	}

	user, err := model.NewUser(validUserID, name)
	require.NoError(t, err)
	require.NotNil(t, user.Name.MiddleName)
	assert.Equal(t, "Ivanovich", *user.Name.MiddleName)

	middleName = "Mutated outside"
	assert.Equal(t, "Ivanovich", *user.Name.MiddleName)

	newMiddleName := "Petrovich"
	newName := model.UserName{
		GivenName:  "Petr",
		FamilyName: "Petrov",
		MiddleName: &newMiddleName,
	}

	err = user.UpdateName(newName)
	require.NoError(t, err)
	require.NotNil(t, user.Name.MiddleName)
	assert.Equal(t, "Petrovich", *user.Name.MiddleName)

	newMiddleName = "Changed again"
	assert.Equal(t, "Petrovich", *user.Name.MiddleName)
}

func TestUser_RemoveCustomName(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("Base", "Name", nil)
	user, _ := model.NewUser(validUserID, un)
	user.PopEvents() // Clear initial event

	// Error if already clear
	err := user.RemoveCustomName()
	assert.ErrorIs(t, err, errors.ErrUserCustomNameAlreadyEmpty)

	// Add custom name directly for test
	customName, _ := model.NewUserName("Custom", "Name", nil)
	_ = user.OverrideName(customName)
	user.PopEvents()

	// Clear successfully
	err = user.RemoveCustomName()
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
	user, _ := model.NewUser(validUserID, un)
	customName, _ := model.NewUserName("Custom", "Name", nil)
	user.PopEvents()

	err := user.OverrideName(customName)
	require.NoError(t, err)
	assert.NotNil(t, user.CustomName)
	assert.True(t, user.CustomName.Equals(customName))
	pn := user.PreferredName()
	assert.Equal(t, "Name Custom", (&pn).DisplayName())

	events := user.PopEvents()
	require.Len(t, events, 1)
	assert.IsType(t, model.UserCustomNameUpdatedEvent{}, events[0])

	// Verify idempotency (same name should skip update)
	user.PopEvents()
	err = user.OverrideName(customName)
	require.NoError(t, err)
	assert.Len(t, user.PopEvents(), 0) // No new events

	// Invalid name check
	invalidCustomName := model.UserName{GivenName: "", FamilyName: "F"}
	err = user.OverrideName(invalidCustomName)
	assert.ErrorIs(t, err, errors.ErrInvalidUserGivenName)
}

func TestUser_UpdateName(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("First", "Last", nil)
	user, _ := model.NewUser(validUserID, un)
	newName, _ := model.NewUserName("NewFirst", "NewLast", nil)
	user.PopEvents()

	err := user.UpdateName(newName)
	require.NoError(t, err)
	assert.True(t, user.Name.Equals(newName))

	events := user.PopEvents()
	require.Len(t, events, 1)
	assert.IsType(t, model.UserNameUpdatedEvent{}, events[0])

	// Idempotency check
	user.PopEvents()
	err = user.UpdateName(newName)
	require.NoError(t, err)
	assert.Len(t, user.PopEvents(), 0)

	// Validation
	invalidName := model.UserName{GivenName: "B", FamilyName: ""}
	err = user.UpdateName(invalidName)
	assert.ErrorIs(t, err, errors.ErrInvalidUserFamilyName)
}

func TestUser_AdminStatus(t *testing.T) {
	t.Parallel()
	un, _ := model.NewUserName("Admin", "User", nil)
	user, _ := model.NewUser(validUserID, un)
	user.PopEvents()

	// Default is not admin
	assert.Nil(t, user.AdminRole)
	assert.False(t, user.IsAdmin())

	now := time.Now().UTC()
	oldAdminTime := now.Add(-model.AdminMinimumTenure - 1*time.Hour)

	// Create an admin actor with sufficient tenure
	actorUn, _ := model.NewUserName("Actor", "User", nil)
	actor, _ := model.NewUser(int64(2<<23), actorUn)
	actor.AdminRole = &model.AdminRole{
		GrantedAt: oldAdminTime,
		GrantedBy: int64(3 << 23),
	}
	actor.PopEvents()

	// Promote
	err := user.GrantAdminRole(actor)
	require.NoError(t, err)
	assert.NotNil(t, user.AdminRole)
	assert.False(t, user.AdminRole.GrantedAt.IsZero())
	assert.Equal(t, actor.ID, user.AdminRole.GrantedBy)
	assert.True(t, user.IsAdmin())

	events := user.PopEvents()
	require.Len(t, events, 1)
	grantEvent, ok := events[0].(model.UserGrantedAdminRoleEvent)
	require.True(t, ok)
	assert.Equal(t, actor.ID, grantEvent.GrantedBy)

	// Promote again (idempotent)
	err = user.GrantAdminRole(actor)
	assert.ErrorIs(t, err, errors.ErrUserAlreadyAdmin)
	assert.NotNil(t, user.AdminRole)
	assert.False(t, user.AdminRole.GrantedAt.IsZero())
	assert.Equal(t, actor.ID, user.AdminRole.GrantedBy) // Should still be old ID
	assert.Empty(t, user.PopEvents())

	// Demote
	err = user.RevokeAdminRole(actor)
	require.NoError(t, err)
	assert.Nil(t, user.AdminRole)
	assert.False(t, user.IsAdmin())

	events = user.PopEvents()
	require.Len(t, events, 1)
	revokeEvent, ok := events[0].(model.UserRevokedAdminRoleEvent)
	require.True(t, ok)
	assert.False(t, revokeEvent.RevokedAt.IsZero())
	assert.Equal(t, actor.ID, revokeEvent.RevokedBy)

	// Demote again (idempotent)
	err = user.RevokeAdminRole(actor)
	require.NoError(t, err)
	assert.Nil(t, user.AdminRole)
	assert.Empty(t, user.PopEvents())
}

func TestUser_RevokeAdminRole_SelfRevokeForbidden(t *testing.T) {
	t.Parallel()

	adminName, _ := model.NewUserName("Admin", "User", nil)
	admin, _ := model.NewUser(validUserID, adminName)
	admin.AdminRole = &model.AdminRole{
		GrantedAt: time.Now().UTC().Add(-model.AdminMinimumTenure - time.Hour),
		GrantedBy: int64(2 << 23),
	}
	admin.PopEvents()

	err := admin.RevokeAdminRole(admin)
	assert.ErrorIs(t, err, errors.ErrAdminSelfRevokeForbidden)
	assert.NotNil(t, admin.AdminRole)
	assert.Empty(t, admin.PopEvents())
}

func TestRestoreUser_AdminRoleInvariant(t *testing.T) {
	t.Parallel()

	un, _ := model.NewUserName("Admin", "User", nil)
	now := time.Now().UTC()
	granterID := int64(3 << 23)

	t.Run("NonAdminUserRestored", func(t *testing.T) {
		t.Parallel()

		user, err := model.RestoreUser(validUserID, un, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Nil(t, user.AdminRole)
	})

	t.Run("AdminWithoutGrantedAt", func(t *testing.T) {
		t.Parallel()

		adminRole := &model.AdminRole{
			GrantedAt: time.Time{},
			GrantedBy: granterID,
		}
		_, err := model.RestoreUser(validUserID, un, nil, adminRole)
		assert.ErrorIs(t, err, errors.ErrInvalidAdminRoleSince)
	})

	t.Run("ValidAdminRole", func(t *testing.T) {
		t.Parallel()

		adminRole := &model.AdminRole{
			GrantedAt: now,
			GrantedBy: granterID,
		}
		user, err := model.RestoreUser(validUserID, un, nil, adminRole)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotNil(t, user.AdminRole)
		assert.Equal(t, now, user.AdminRole.GrantedAt)
		assert.Equal(t, granterID, user.AdminRole.GrantedBy)
	})
}

func TestUser_EmploymentAsEntity_MultipleEmployments(t *testing.T) {
	t.Parallel()

	un, _ := model.NewUserName("Employee", "User", nil)
	user, _ := model.NewUser(validUserID, un)
	user.PopEvents()

	organizationID := uuid.MustParse("33333333-3333-7333-8333-333333333333")
	secondOrganizationID := uuid.MustParse("44444444-4444-7444-8444-444444444444")
	clinicID := uuid.MustParse("55555555-5555-7555-8555-555555555555")
	secondClinicID := uuid.MustParse("66666666-6666-7666-8666-666666666666")
	departmentID := uuid.MustParse("11111111-1111-7111-8111-111111111111")
	secondDepartmentID := uuid.MustParse("22222222-2222-7222-8222-222222222222")

	firstEmploymentID, err := user.AssignEmployment(organizationID, clinicID, departmentID, nil)
	require.NoError(t, err)
	require.Len(t, user.Employments, 1)
	assert.Equal(t, departmentID, user.Employments[0].DepartmentID)
	assert.Equal(t, clinicID, user.Employments[0].ClinicID)
	assert.Equal(t, organizationID, user.Employments[0].OrganizationID)
	assert.Equal(t, firstEmploymentID, user.Employments[0].ID)
	assert.False(t, user.Employments[0].AssignedAt.IsZero())
	assert.True(t, user.IsEmployeeOfOrganization(organizationID))
	assert.True(t, user.IsEmployeeOfClinic(clinicID))
	assert.True(t, user.IsEmployeeOfDepartment(departmentID))

	secondEmploymentID, err := user.AssignEmployment(secondOrganizationID, secondClinicID, secondDepartmentID, ptr("Second role"))
	require.NoError(t, err)
	require.Len(t, user.Employments, 2)
	assert.Equal(t, secondEmploymentID, user.Employments[1].ID)
	assert.True(t, user.IsEmployeeOfOrganization(secondOrganizationID))
	assert.True(t, user.IsEmployeeOfClinic(secondClinicID))
	assert.True(t, user.IsEmployeeOfDepartment(secondDepartmentID))

	events := user.PopEvents()
	require.Len(t, events, 2)
	firstEmployed, ok := events[0].(model.UserEmployedEvent)
	require.True(t, ok)
	assert.Equal(t, firstEmploymentID, firstEmployed.EmploymentID)
	assert.Equal(t, organizationID, firstEmployed.OrganizationID)
	assert.Equal(t, clinicID, firstEmployed.ClinicID)
	assert.Equal(t, departmentID, firstEmployed.DepartmentID)
	assert.False(t, firstEmployed.AssignedAt.IsZero())
	secondEmployed, ok := events[1].(model.UserEmployedEvent)
	require.True(t, ok)
	assert.Equal(t, secondEmploymentID, secondEmployed.EmploymentID)
	assert.Equal(t, secondOrganizationID, secondEmployed.OrganizationID)
	assert.Equal(t, secondClinicID, secondEmployed.ClinicID)
	assert.Equal(t, secondDepartmentID, secondEmployed.DepartmentID)
	assert.False(t, secondEmployed.AssignedAt.IsZero())

	firstEmployment := user.Employments[0]
	firstEmployment.PopEvents()

	err = firstEmployment.AssignDeputy(int64(2 << 23))
	require.NoError(t, err)
	assert.True(t, user.Employments[0].HasDeputy())

	events = firstEmployment.PopEvents()
	require.Len(t, events, 1)
	deputyAssigned, ok := events[0].(model.EmploymentDeputyAssignedEvent)
	require.True(t, ok)
	assert.Equal(t, firstEmploymentID, deputyAssigned.EmploymentID)
	assert.Equal(t, int64(2<<23), deputyAssigned.DeputyID)

	err = firstEmployment.AssignDeputy(validUserID)
	assert.ErrorIs(t, err, errors.ErrInvalidEmploymentDeputy)

	startsAt := time.Now().UTC().Add(24 * time.Hour)
	endsAt := startsAt.Add(24 * time.Hour)
	err = firstEmployment.ScheduleVacation(startsAt, &endsAt)
	require.NoError(t, err)
	assert.True(t, user.Employments[0].HasScheduledVacation())

	events = firstEmployment.PopEvents()
	require.Len(t, events, 1)
	vacationScheduled, ok := events[0].(model.EmploymentVacationScheduledEvent)
	require.True(t, ok)
	assert.Equal(t, firstEmploymentID, vacationScheduled.EmploymentID)

	firstEmployment.EndVacation()
	assert.Nil(t, user.Employments[0].Vacation)

	events = firstEmployment.PopEvents()
	require.Len(t, events, 1)
	vacationEnded, ok := events[0].(model.EmploymentVacationEndedEvent)
	require.True(t, ok)
	assert.Equal(t, firstEmploymentID, vacationEnded.EmploymentID)

	err = user.Dismiss(firstEmploymentID)
	require.NoError(t, err)
	require.Len(t, user.Employments, 1)
	assert.Equal(t, secondEmploymentID, user.Employments[0].ID)
	assert.False(t, user.IsEmployeeOfOrganization(organizationID))
	assert.False(t, user.IsEmployeeOfClinic(clinicID))
	assert.False(t, user.IsEmployeeOfDepartment(departmentID))

	events = user.PopEvents()
	require.Len(t, events, 1)
	dismissed, ok := events[0].(model.UserDismissedEvent)
	require.True(t, ok)
	assert.Equal(t, firstEmploymentID, dismissed.EmploymentID)
}

func TestUser_AssignEmployment_SameOrganizationForbidden(t *testing.T) {
	t.Parallel()

	un, _ := model.NewUserName("Employee", "User", nil)
	user, _ := model.NewUser(validUserID, un)
	user.PopEvents()

	organizationID := uuid.MustParse("33333333-3333-7333-8333-333333333333")
	clinicID := uuid.MustParse("55555555-5555-7555-8555-555555555555")
	firstDepartmentID := uuid.MustParse("11111111-1111-7111-8111-111111111111")
	secondDepartmentID := uuid.MustParse("22222222-2222-7222-8222-222222222222")

	_, err := user.AssignEmployment(organizationID, clinicID, firstDepartmentID, nil)
	require.NoError(t, err)

	_, err = user.AssignEmployment(organizationID, clinicID, secondDepartmentID, ptr("Second role"))
	assert.ErrorIs(t, err, errors.ErrEmploymentAlreadyExistsInOrganization)
	require.Len(t, user.Employments, 1)
	require.Len(t, user.PopEvents(), 1)
}

func TestUser_AssignEmployment_ClonesPositionPointer(t *testing.T) {
	t.Parallel()

	un, _ := model.NewUserName("Employee", "User", nil)
	user, _ := model.NewUser(validUserID, un)
	user.PopEvents()

	organizationID := uuid.MustParse("33333333-3333-7333-8333-333333333333")
	clinicID := uuid.MustParse("55555555-5555-7555-8555-555555555555")
	departmentID := uuid.MustParse("11111111-1111-7111-8111-111111111111")

	position := "Doctor"
	_, err := user.AssignEmployment(organizationID, clinicID, departmentID, &position)
	require.NoError(t, err)

	position = "Mutated outside"

	require.Len(t, user.Employments, 1)
	require.NotNil(t, user.Employments[0].Position)
	assert.Equal(t, "Doctor", *user.Employments[0].Position)
	assert.NotSame(t, &position, user.Employments[0].Position)

	events := user.PopEvents()
	require.Len(t, events, 1)
	employed, ok := events[0].(model.UserEmployedEvent)
	require.True(t, ok)
	require.NotNil(t, employed.Position)
	assert.Equal(t, "Doctor", *employed.Position)
	assert.NotSame(t, &position, employed.Position)
	assert.NotSame(t, user.Employments[0].Position, employed.Position)
}

func TestUser_CanGrantAdminRole(t *testing.T) {
	t.Parallel()

	un, _ := model.NewUserName("Admin", "User", nil)
	actor, _ := model.NewUser(validUserID, un)
	now := time.Now().UTC()
	granterID := int64(4 << 23)

	t.Run("NotAdmin", func(t *testing.T) {
		t.Parallel()

		err := actor.CanGrantAdminRole()
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantForbidden)
	})

	t.Run("AdminLessThan72Hours", func(t *testing.T) {
		t.Parallel()

		recentAdminSince := now.Add(-model.AdminMinimumTenure + 1*time.Hour)
		adminRole := &model.AdminRole{
			GrantedAt: recentAdminSince,
			GrantedBy: granterID,
		}
		admin, err := model.RestoreUser(validUserID, un, nil, adminRole)
		require.NoError(t, err)

		err = admin.CanGrantAdminRole()
		assert.ErrorIs(t, err, errors.ErrAdminRoleGrantForbidden)
	})

	t.Run("AdminMoreThan72Hours", func(t *testing.T) {
		t.Parallel()

		oldAdminSince := now.Add(-model.AdminMinimumTenure - 1*time.Hour)
		adminRole := &model.AdminRole{
			GrantedAt: oldAdminSince,
			GrantedBy: granterID,
		}
		admin, err := model.RestoreUser(validUserID, un, nil, adminRole)
		require.NoError(t, err)

		err = admin.CanGrantAdminRole()
		require.NoError(t, err)
	})
}
