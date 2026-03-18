package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minUserNameLength = 1
	maxUserNameLength = 100
)

func validateUserName(name UserName) error {
	if strings.TrimSpace(name.GivenName) != name.GivenName {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidGivenName)
	}
	givenNameLen := utf8.RuneCountInString(name.GivenName)
	if givenNameLen < minUserNameLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidGivenName, minUserNameLength)
	}
	if givenNameLen > maxUserNameLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidGivenName, maxUserNameLength)
	}

	if strings.TrimSpace(name.FamilyName) != name.FamilyName {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidFamilyName)
	}
	familyNameLen := utf8.RuneCountInString(name.FamilyName)
	if familyNameLen < minUserNameLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidFamilyName, minUserNameLength)
	}
	if familyNameLen > maxUserNameLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidFamilyName, maxUserNameLength)
	}

	if name.MiddleName != nil {
		if strings.TrimSpace(*name.MiddleName) != *name.MiddleName {
			return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidMiddleName)
		}
		middleNameLen := utf8.RuneCountInString(*name.MiddleName)
		if middleNameLen < minUserNameLength {
			return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidMiddleName, minUserNameLength)
		}
		if middleNameLen > maxUserNameLength {
			return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidMiddleName, maxUserNameLength)
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

func validateUser(u *User) error {
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
	if u.AdminRole != nil {
		if u.AdminRole.GrantedAt.IsZero() {
			return fmt.Errorf("%w: must not be zero for admin user", errs.ErrInvalidAdminSince)
		}
		if err := validateUserID(u.AdminRole.GrantedBy); err != nil {
			return fmt.Errorf("invalid admin role granter id: %w", err)
		}
	}
	for _, employment := range u.Employments {
		if employment == nil {
			continue
		}
		if err := validateEmployment(employment); err != nil {
			return err
		}
		if employment.UserID != u.ID {
			return fmt.Errorf("invalid employment user id: expected %d, got %d", u.ID, employment.UserID)
		}
		if employment.Deputy != nil && employment.Deputy.ID == u.ID {
			return fmt.Errorf("%w: user cannot be their own deputy", errs.ErrInvalidEmploymentDeputy)
		}
	}
	return nil
}
