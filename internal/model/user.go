package model

import (
	"fmt"
	"time"

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
	ID         int64
	Name       UserName
	CustomName *UserName
	AdminRole  *AdminRole
}

func NewUser(id int64, name UserName) (*User, error) {
	u := &User{
		ID:        id,
		Name:      name,
		AdminRole: nil,
	}
	if err := validateUser(*u); err != nil {
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
		ID:         id,
		Name:       name,
		CustomName: customName,
		AdminRole:  adminRole,
	}
	if err := validateUser(*u); err != nil {
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

func (u *User) ClearCustomName() error {
	if u.CustomName == nil {
		return errors.ErrCustomNameAlreadyEmpty
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
		return fmt.Errorf("%w: actor admin tenure is less than 72 hours", errors.ErrAdminRoleGrantForbidden)
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
