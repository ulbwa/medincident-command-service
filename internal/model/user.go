package model

import (
	"fmt"
)

type UserName struct {
	GivenName  string
	FamilyName string
	MiddleName *string
}

func (u *UserName) DisplayName() string {
	if u.MiddleName != nil {
		return fmt.Sprintf("%s %s %s", u.FamilyName, u.GivenName, *u.MiddleName)
	}
	return fmt.Sprintf("%s %s", u.FamilyName, u.GivenName)
}

func (u *UserName) ShortName() string {
	givenInitial := []rune(u.GivenName)[0]
	if u.MiddleName != nil {
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

func NewUserName(givenName, familyName string, middleName *string) (*UserName, error) {
	name := UserName{
		GivenName:  givenName,
		FamilyName: familyName,
		MiddleName: middleName,
	}
	if err := validateUserName(name); err != nil {
		return nil, err
	}
	return &name, nil
}

type User struct {
	Entity
	ID         int64
	Name       UserName
	CustomName *UserName
}

func NewUser(id int64, name UserName) (*User, error) {
	u := &User{
		ID:   id,
		Name: name,
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

func RestoreUser(id int64, name UserName, customName *UserName) (*User, error) {
	u := &User{
		ID:         id,
		Name:       name,
		CustomName: customName,
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

func (u *User) ClearCustomName() {
	if u.CustomName == nil {
		return
	}

	u.CustomName = nil
	u.recordEvent(UserCustomNameUpdatedEvent{
		ID:         u.ID,
		CustomName: nil,
	})
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
