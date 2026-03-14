package model

import (
	"fmt"
	"unicode/utf8"

	errs "github.com/ulbwa/medincident-command-service/internal/errors"
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

func validateUserName(name UserName) error {
	givenNameLen := utf8.RuneCountInString(name.GivenName)
	if givenNameLen < 1 {
		return fmt.Errorf("%w: too short (min 1)", errs.ErrInvalidGivenName)
	}
	if givenNameLen > 100 {
		return fmt.Errorf("%w: too long (max 100)", errs.ErrInvalidGivenName)
	}

	familyNameLen := utf8.RuneCountInString(name.FamilyName)
	if familyNameLen < 1 {
		return fmt.Errorf("%w: too short (min 1)", errs.ErrInvalidFamilyName)
	}
	if familyNameLen > 100 {
		return fmt.Errorf("%w: too long (max 100)", errs.ErrInvalidFamilyName)
	}

	if name.MiddleName != nil {
		middleNameLen := utf8.RuneCountInString(*name.MiddleName)
		if middleNameLen < 1 {
			return fmt.Errorf("%w: too short (min 1)", errs.ErrInvalidMiddleName)
		}
		if middleNameLen > 100 {
			return fmt.Errorf("%w: too long (max 100)", errs.ErrInvalidMiddleName)
		}
	}
	return nil
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
	ID         int64
	Name       UserName
	CustomName *UserName
}

// validateUserID check that the user ID is a valid Snowflake ID from Zitadel
func validateUserID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: must be greater than zero", errs.ErrInvalidUserID)
	}
	// Check that the timestamp component of the Snowflake ID is greater than zero
	if (id >> 22) <= 0 {
		return fmt.Errorf("%w: timestamp component must be greater than zero", errs.ErrInvalidUserID)
	}
	return nil
}

func validateUser(u User) error {
	if err := validateUserID(u.ID); err != nil {
		return err
	}
	if err := validateUserName(u.Name); err != nil {
		return err
	}
	if u.CustomName != nil {
		if err := validateUserName(*u.CustomName); err != nil {
			return err
		}
	}
	return nil
}

func NewUser(id int64, name UserName) (*User, error) {
	u := &User{
		ID:   id,
		Name: name,
	}
	if err := validateUser(*u); err != nil {
		return nil, err
	}
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
	if err := validateUserName(customName); err != nil {
		return err
	}
	u.CustomName = &customName
	return nil
}

func (u *User) ClearCustomName() {
	u.CustomName = nil
}

func (u *User) UpdateName(name UserName) error {
	if err := validateUserName(name); err != nil {
		return err
	}
	u.Name = name
	return nil
}
