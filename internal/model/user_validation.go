package model

import (
	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minUserNameLength = 1
	maxUserNameLength = 100
)

func validateUserName(name UserName) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(name.GivenName); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldGivenName, err)
	}
	if err := validateStringNoConsecutiveSpaces(name.GivenName); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldGivenName, err)
	}
	if err := validateStringMinLength(name.GivenName, minUserNameLength); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldGivenName, err)
	}
	if err := validateStringMaxLength(name.GivenName, maxUserNameLength); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldGivenName, err)
	}

	if err := validateStringNoLeadingOrTrailingWhitespace(name.FamilyName); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldFamilyName, err)
	}
	if err := validateStringNoConsecutiveSpaces(name.FamilyName); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldFamilyName, err)
	}
	if err := validateStringMinLength(name.FamilyName, minUserNameLength); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldFamilyName, err)
	}
	if err := validateStringMaxLength(name.FamilyName, maxUserNameLength); err != nil {
		return errs.NewInvalidUserNameError(errs.UserNameFieldFamilyName, err)
	}

	if name.MiddleName != nil {
		if err := validateStringNoLeadingOrTrailingWhitespace(*name.MiddleName); err != nil {
			return errs.NewInvalidUserNameError(errs.UserNameFieldMiddleName, err)
		}
		if err := validateStringNoConsecutiveSpaces(*name.MiddleName); err != nil {
			return errs.NewInvalidUserNameError(errs.UserNameFieldMiddleName, err)
		}
		if err := validateStringMinLength(*name.MiddleName, minUserNameLength); err != nil {
			return errs.NewInvalidUserNameError(errs.UserNameFieldMiddleName, err)
		}
		if err := validateStringMaxLength(*name.MiddleName, maxUserNameLength); err != nil {
			return errs.NewInvalidUserNameError(errs.UserNameFieldMiddleName, err)
		}
	}

	return nil
}

// validateUserID checks that the user ID is a valid Snowflake ID from Zitadel.
func validateUserID(id int64) error {
	return validateSnowflakeID(id)
}

func validateAdminRole(adminRole AdminRole) error {
	if adminRole.GrantedAt.IsZero() {
		return errs.NewInvalidAdminRoleError(errs.AdminRoleFieldGrantedAt, errs.NewValueRequiredError())
	}
	if err := validateUserID(adminRole.GrantedBy); err != nil {
		return errs.NewInvalidAdminRoleError(errs.AdminRoleFieldGrantedBy, err)
	}
	return nil
}

func validateUser(u *User) error {
	if err := validateUserID(u.ID); err != nil {
		return errs.NewInvalidUserError(errs.UserFieldID, err)
	}
	if err := validateUserName(u.Name); err != nil {
		return errs.NewInvalidUserError(errs.UserFieldName, err)
	}
	if u.CustomName != nil {
		if err := validateUserName(*u.CustomName); err != nil {
			return errs.NewInvalidUserError(errs.UserFieldCustomName, err)
		}
	}
	if u.AdminRole != nil {
		if err := validateAdminRole(*u.AdminRole); err != nil {
			return errs.NewInvalidUserError(errs.UserFieldAdminRole, err)
		}
	}
	// validateUserEmployments already wraps errors in *InvalidUserError with item indexes preserved.
	if err := validateUserEmployments(u); err != nil {
		return err
	}
	return nil
}

func validateUserEmployments(u *User) error {
	organizationIndexes := make(map[uuid.UUID]int, len(u.Employments))

	for index, employment := range u.Employments {
		if employment == nil {
			return errs.NewInvalidUserError(
				errs.UserFieldEmployments,
				errs.NewInvalidCollectionItemError(index, errs.NewValueRequiredError()),
			)
		}

		// Validate the employment first so that structural errors (e.g. invalid
		// OrganizationID) are surfaced before the uniqueness check. Otherwise two
		// employments with uuid.Nil would be reported as a duplicate instead of
		// an invalid OrganizationID.
		if err := validateEmployment(employment); err != nil {
			return errs.NewInvalidUserError(
				errs.UserFieldEmployments,
				errs.NewInvalidCollectionItemError(index, err),
			)
		}

		if _, exists := organizationIndexes[employment.OrganizationID]; exists {
			return errs.NewInvalidUserError(
				errs.UserFieldEmployments,
				errs.NewInvalidCollectionItemError(
					index,
					errs.NewInvalidEmploymentError(
						errs.EmploymentFieldOrganizationID,
						errs.NewValueDuplicateError(employment.OrganizationID),
					),
				),
			)
		}
		organizationIndexes[employment.OrganizationID] = index

		if employment.UserID != u.ID {
			return errs.NewInvalidUserError(
				errs.UserFieldEmployments,
				errs.NewInvalidCollectionItemError(
					index,
					errs.NewInvalidEmploymentError(
						errs.EmploymentFieldUserID,
						errs.NewValueMismatchError(u.ID, employment.UserID),
					),
				),
			)
		}
	}
	return nil
}
