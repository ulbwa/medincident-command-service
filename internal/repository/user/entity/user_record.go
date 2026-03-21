// Package entity defines storage record types and converters for the user context.
// Nothing from this package should be exposed to the service or handler layers.
package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ulbwa/medincident-command-service/internal/model"
)

// UserRecord maps directly to a row in the users table.
type UserRecord struct {
	ID                 uuid.UUID  `db:"id"`
	IdentityID         string     `db:"identity_id"`
	GivenName          string     `db:"given_name"`
	FamilyName         string     `db:"family_name"`
	MiddleName         *string    `db:"middle_name"`
	CustomGivenName    *string    `db:"custom_given_name"`
	CustomFamilyName   *string    `db:"custom_family_name"`
	CustomMiddleName   *string    `db:"custom_middle_name"`
	AdminRoleGrantedAt *time.Time `db:"admin_role_granted_at"`
	AdminRoleGrantedBy *uuid.UUID `db:"admin_role_granted_by"`
}

// EmploymentRecord maps directly to a row in the employments table.
type EmploymentRecord struct {
	ID               uuid.UUID  `db:"id"`
	UserID           uuid.UUID  `db:"user_id"`
	OrganizationID   uuid.UUID  `db:"organization_id"`
	ClinicID         uuid.UUID  `db:"clinic_id"`
	DepartmentID     uuid.UUID  `db:"department_id"`
	Position         *string    `db:"position"`
	AssignedAt       time.Time  `db:"assigned_at"`
	DeputyUserID     *uuid.UUID `db:"deputy_user_id"`
	VacationStartsAt *time.Time `db:"vacation_starts_at"`
	VacationEndsAt   *time.Time `db:"vacation_ends_at"`
}

// FromUser converts a User aggregate to its storage record.
func FromUser(u *model.User) UserRecord {
	rec := UserRecord{
		ID:         u.ID,
		IdentityID: u.IdentityID,
		GivenName:  u.Name.GivenName,
		FamilyName: u.Name.FamilyName,
		MiddleName: u.Name.MiddleName,
	}

	if u.CustomName != nil {
		rec.CustomGivenName = &u.CustomName.GivenName
		rec.CustomFamilyName = &u.CustomName.FamilyName
		rec.CustomMiddleName = u.CustomName.MiddleName
	}

	if u.AdminRole != nil {
		at := u.AdminRole.GrantedAt
		by := u.AdminRole.GrantedBy
		rec.AdminRoleGrantedAt = &at
		rec.AdminRoleGrantedBy = &by
	}

	return rec
}

// ToUser reconstructs a User aggregate from its storage records.
func ToUser(rec UserRecord, empRecs []EmploymentRecord) (*model.User, error) {
	name, err := model.NewUserName(rec.GivenName, rec.FamilyName, rec.MiddleName)
	if err != nil {
		return nil, fmt.Errorf("restore user name: %w", err)
	}

	var customName *model.UserName
	if rec.CustomGivenName != nil && rec.CustomFamilyName != nil {
		cn, cnErr := model.NewUserName(*rec.CustomGivenName, *rec.CustomFamilyName, rec.CustomMiddleName)
		if cnErr != nil {
			return nil, fmt.Errorf("restore custom name: %w", cnErr)
		}

		customName = &cn
	}

	var adminRole *model.AdminRole
	if rec.AdminRoleGrantedAt != nil && rec.AdminRoleGrantedBy != nil {
		adminRole = &model.AdminRole{
			GrantedAt: *rec.AdminRoleGrantedAt,
			GrantedBy: *rec.AdminRoleGrantedBy,
		}
	}

	employments := make([]*model.Employment, 0, len(empRecs))
	for _, empRec := range empRecs {
		emp, empErr := toEmployment(empRec)
		if empErr != nil {
			return nil, fmt.Errorf("restore employment %s: %w", empRec.ID, empErr)
		}

		employments = append(employments, emp)
	}

	return model.RestoreUser(rec.ID, rec.IdentityID, name, customName, adminRole, employments)
}

func toEmployment(rec EmploymentRecord) (*model.Employment, error) {
	var deputy *model.EmploymentDeputy
	if rec.DeputyUserID != nil {
		d, err := model.NewEmploymentDeputy(*rec.DeputyUserID)
		if err != nil {
			return nil, fmt.Errorf("restore deputy: %w", err)
		}

		deputy = &d
	}

	var vacation *model.EmploymentVacation
	if rec.VacationStartsAt != nil {
		vacation = &model.EmploymentVacation{
			StartsAt: *rec.VacationStartsAt,
			EndsAt:   rec.VacationEndsAt,
		}
	}

	return model.RestoreEmployment(
		rec.ID,
		rec.UserID,
		rec.OrganizationID,
		rec.ClinicID,
		rec.DepartmentID,
		rec.Position,
		rec.AssignedAt,
		deputy,
		vacation,
	)
}

// FromEmployments converts a slice of Employment aggregates to storage records.
func FromEmployments(employments []*model.Employment) []EmploymentRecord {
	recs := make([]EmploymentRecord, 0, len(employments))
	for _, emp := range employments {
		if emp == nil {
			continue
		}

		recs = append(recs, fromEmployment(emp))
	}

	return recs
}

func fromEmployment(e *model.Employment) EmploymentRecord {
	rec := EmploymentRecord{
		ID:             e.ID,
		UserID:         e.UserID,
		OrganizationID: e.OrganizationID,
		ClinicID:       e.ClinicID,
		DepartmentID:   e.DepartmentID,
		Position:       e.Position,
		AssignedAt:     e.AssignedAt,
	}

	if e.Deputy != nil {
		rec.DeputyUserID = &e.Deputy.ID
	}

	if e.Vacation != nil {
		rec.VacationStartsAt = &e.Vacation.StartsAt
		rec.VacationEndsAt = e.Vacation.EndsAt
	}

	return rec
}
