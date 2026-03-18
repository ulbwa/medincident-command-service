package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const AdminMinimumTenure = 72 * time.Hour

type UserName struct {
	GivenName  string
	FamilyName string
	MiddleName *string
}

func (u UserName) DisplayName() string {
	if u.MiddleName != nil {
		return fmt.Sprintf("%s %s %s", u.FamilyName, u.GivenName, *u.MiddleName)
	}
	return fmt.Sprintf("%s %s", u.FamilyName, u.GivenName)
}

func (u UserName) ShortName() string {
	if u.GivenName == "" {
		return ""
	}
	givenInitial := []rune(u.GivenName)[0]
	if u.MiddleName != nil && *u.MiddleName != "" {
		middleInitial := []rune(*u.MiddleName)[0]
		return fmt.Sprintf("%s %c.%c.", u.FamilyName, givenInitial, middleInitial)
	}
	return fmt.Sprintf("%s %c.", u.FamilyName, givenInitial)
}

// Equals safely compares two UserName structures, handling the MiddleName pointer.
func (u UserName) Equals(other UserName) bool {
	if u.GivenName != other.GivenName || u.FamilyName != other.FamilyName {
		return false
	}

	if (u.MiddleName == nil && other.MiddleName != nil) || (u.MiddleName != nil && other.MiddleName == nil) {
		return false
	}

	if u.MiddleName != nil && other.MiddleName != nil && *u.MiddleName != *other.MiddleName {
		return false
	}

	return true
}

func NewUserName(givenName, familyName string, middleName *string) (UserName, error) {
	name := UserName{
		GivenName:  givenName,
		FamilyName: familyName,
		MiddleName: middleName,
	}
	if err := validateUserName(name); err != nil {
		return UserName{}, err
	}
	return name, nil
}

type AdminRole struct {
	GrantedAt time.Time
	GrantedBy int64
}

type User struct {
	Entity
	ID          int64
	Name        UserName
	CustomName  *UserName
	AdminRole   *AdminRole
	Employments []*Employment
}

func NewUser(id int64, name UserName) (*User, error) {
	u := &User{
		ID:          id,
		Name:        name,
		AdminRole:   nil,
		Employments: make([]*Employment, 0),
	}
	if err := validateUser(u); err != nil {
		return nil, err
	}

	u.recordEvent(UserCreatedEvent{
		ID:   u.ID,
		Name: name,
	})

	return u, nil
}

func RestoreUser(id int64, name UserName, customName *UserName, adminRole *AdminRole) (*User, error) {
	u := &User{
		ID:          id,
		Name:        name,
		CustomName:  customName,
		AdminRole:   adminRole,
		Employments: make([]*Employment, 0),
	}
	if err := validateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) PreferredName() UserName {
	if u.CustomName != nil {
		return *u.CustomName
	}
	return u.Name
}

func (u *User) IsAdmin() bool {
	return u.AdminRole != nil
}

func (u *User) OverrideName(customName UserName) error {
	if u.CustomName != nil && u.CustomName.Equals(customName) {
		return nil
	}

	if err := validateUserName(customName); err != nil {
		return err
	}

	u.CustomName = &customName

	u.recordEvent(UserCustomNameUpdatedEvent{
		ID:         u.ID,
		CustomName: u.CustomName,
	})

	return nil
}

func (u *User) RemoveCustomName() error {
	if u.CustomName == nil {
		return errors.ErrUserCustomNameAlreadyEmpty
	}

	u.CustomName = nil
	u.recordEvent(UserCustomNameUpdatedEvent{
		ID:         u.ID,
		CustomName: nil,
	})

	return nil
}

func (u *User) UpdateName(name UserName) error {
	if u.Name.Equals(name) {
		return nil
	}

	if err := validateUserName(name); err != nil {
		return err
	}

	u.Name = name
	u.recordEvent(UserNameUpdatedEvent{
		ID:   u.ID,
		Name: name,
	})

	return nil
}

func (u *User) CanGrantAdminRole() error {
	if !u.IsAdmin() {
		return fmt.Errorf("%w: actor is not admin", errors.ErrAdminRoleGrantForbidden)
	}

	if u.AdminRole.GrantedAt.IsZero() {
		return fmt.Errorf("%w: admin since must not be zero time", errors.ErrAdminRoleGrantForbidden)
	}

	if time.Since(u.AdminRole.GrantedAt) < AdminMinimumTenure {
		return fmt.Errorf("%w: actor admin tenure is less than %.0f hours", errors.ErrAdminRoleGrantForbidden, AdminMinimumTenure.Hours())
	}

	return nil
}

func (u *User) GrantAdminRole(actor *User) error {
	if actor == nil {
		return fmt.Errorf("%w: actor is nil", errors.ErrAdminRoleGrantForbidden)
	}

	if err := actor.CanGrantAdminRole(); err != nil {
		return err
	}

	if u.IsAdmin() {
		return nil // idempotent
	}

	at := time.Now().UTC()

	u.AdminRole = &AdminRole{
		GrantedAt: at,
		GrantedBy: actor.ID,
	}
	u.recordEvent(UserGrantedAdminRoleEvent{
		ID:        u.ID,
		GrantedAt: at,
		GrantedBy: actor.ID,
	})

	return nil
}

func (u *User) RevokeAdminRole(actor *User) error {
	if actor == nil {
		return fmt.Errorf("%w: actor is nil", errors.ErrAdminRoleGrantForbidden)
	}

	if err := actor.CanGrantAdminRole(); err != nil {
		return err
	}

	if !u.IsAdmin() {
		return nil // idempotent
	}

	revokedAt := time.Now().UTC()

	u.AdminRole = nil
	u.recordEvent(UserRevokedAdminRoleEvent{
		ID:        u.ID,
		RevokedAt: revokedAt,
		RevokedBy: actor.ID,
	})

	return nil
}

func (u *User) findEmploymentIndex(employmentID uuid.UUID) int {
	for index, employment := range u.Employments {
		if employment != nil && employment.ID == employmentID {
			return index
		}
	}
	return -1
}

func (u *User) IsEmployeeOfOrganization(organizationID uuid.UUID) bool {
	for _, employment := range u.Employments {
		if employment != nil && employment.OrganizationID == organizationID {
			return true
		}
	}
	return false
}

func (u *User) IsEmployeeOfClinic(clinicID uuid.UUID) bool {
	for _, employment := range u.Employments {
		if employment != nil && employment.ClinicID == clinicID {
			return true
		}
	}
	return false
}

func (u *User) IsEmployeeOfDepartment(departmentID uuid.UUID) bool {
	for _, employment := range u.Employments {
		if employment != nil && employment.DepartmentID == departmentID {
			return true
		}
	}
	return false
}

func (u *User) AssignEmployment(organizationID, clinicID, departmentID uuid.UUID, position *string) (uuid.UUID, error) {
	if u.IsEmployeeOfOrganization(organizationID) {
		return uuid.Nil, errors.ErrEmploymentAlreadyExistsInOrganization
	}

	assignedAt := time.Now().UTC()

	employment, err := NewEmployment(u.ID, organizationID, clinicID, departmentID, position, assignedAt)
	if err != nil {
		return uuid.Nil, err
	}

	u.Employments = append(u.Employments, employment)
	u.recordEvent(UserEmployedEvent{
		UserID:         u.ID,
		EmploymentID:   employment.ID,
		OrganizationID: organizationID,
		ClinicID:       clinicID,
		DepartmentID:   departmentID,
		Position:       position,
		AssignedAt:     assignedAt,
	})

	return employment.ID, nil
}

func (u *User) Dismiss(employmentID uuid.UUID) error {
	index := u.findEmploymentIndex(employmentID)
	if index < 0 {
		return errors.ErrEmploymentNotFound
	}

	u.Employments = append(u.Employments[:index], u.Employments[index+1:]...)
	u.recordEvent(UserDismissedEvent{
		UserID:       u.ID,
		EmploymentID: employmentID,
	})

	return nil
}
