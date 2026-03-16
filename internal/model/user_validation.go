package model

import (
	"fmt"
	"unicode/utf8"

	errs "github.com/ulbwa/medincident-command-service/internal/errors"
)

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

// validateUserID checks that the user ID is a valid Snowflake ID from Zitadel
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
