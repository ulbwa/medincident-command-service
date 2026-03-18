package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/pkg/utils"
)

type EmploymentDeputy struct {
	ID int64
}

func (d EmploymentDeputy) copy() *EmploymentDeputy {
	cloned := d
	return &cloned
}

func RestoreEmploymentDeputy(deputyID int64) (EmploymentDeputy, error) {
	deputy, err := newEmploymentDeputy(deputyID)
	if err != nil {
		return EmploymentDeputy{}, err
	}
	return deputy, nil
}

func newEmploymentDeputy(deputyID int64) (EmploymentDeputy, error) {
	if err := validateDeputyID(deputyID); err != nil {
		return EmploymentDeputy{}, err
	}
	return EmploymentDeputy{ID: deputyID}, nil
}

type EmploymentVacation struct {
	StartsAt time.Time
	EndsAt   *time.Time
}

func (v EmploymentVacation) copy() *EmploymentVacation {
	cloned := v
	if v.EndsAt != nil {
		endsAt := *v.EndsAt
		cloned.EndsAt = &endsAt
	}
	return &cloned
}

func RestoreEmploymentVacation(startsAt time.Time, endsAt *time.Time) (EmploymentVacation, error) {
	vacation, err := newEmploymentVacation(startsAt, endsAt)
	if err != nil {
		return EmploymentVacation{}, err
	}
	return vacation, nil
}

func newEmploymentVacation(startsAt time.Time, endsAt *time.Time) (EmploymentVacation, error) {
	v := EmploymentVacation{
		StartsAt: startsAt,
		EndsAt:   utils.PtrClone(endsAt),
	}
	if err := validateEmploymentVacation(&v); err != nil {
		return EmploymentVacation{}, err
	}
	return v, nil
}

func (v EmploymentVacation) IsActive() bool {
	at := time.Now().UTC()
	if at.Before(v.StartsAt) {
		return false
	}
	if v.EndsAt == nil {
		return true
	}
	return !at.After(*v.EndsAt)
}

func (v EmploymentVacation) IsScheduled() bool {
	at := time.Now().UTC()
	return at.Before(v.StartsAt)
}

type Employment struct {
	Entity
	ID             uuid.UUID
	UserID         int64
	OrganizationID uuid.UUID
	ClinicID       uuid.UUID
	DepartmentID   uuid.UUID
	Position       *string
	AssignedAt     time.Time
	Deputy         *EmploymentDeputy
	Vacation       *EmploymentVacation
}

func NewEmployment(userID int64, organizationID, clinicID, departmentID uuid.UUID, position *string, assignedAt time.Time) (*Employment, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	e := &Employment{
		ID:             id,
		UserID:         userID,
		OrganizationID: organizationID,
		ClinicID:       clinicID,
		DepartmentID:   departmentID,
		Position:       utils.PtrClone(position),
		AssignedAt:     assignedAt,
		Deputy:         nil,
		Vacation:       nil,
	}
	if err := validateEmployment(e); err != nil {
		return nil, err
	}

	return e, nil
}

func RestoreEmployment(id uuid.UUID, userID int64, organizationID, clinicID, departmentID uuid.UUID, position *string, assignedAt time.Time, deputy *EmploymentDeputy, vacation *EmploymentVacation) (*Employment, error) {
	var deputyCopy *EmploymentDeputy
	if deputy != nil {
		deputyCopy = deputy.copy()
	}

	var vacationCopy *EmploymentVacation
	if vacation != nil {
		vacationCopy = vacation.copy()
	}

	e := &Employment{
		ID:             id,
		UserID:         userID,
		OrganizationID: organizationID,
		ClinicID:       clinicID,
		DepartmentID:   departmentID,
		Position:       utils.PtrClone(position),
		AssignedAt:     assignedAt,
		Deputy:         deputyCopy,
		Vacation:       vacationCopy,
	}
	if err := validateEmployment(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (e Employment) HasDeputy() bool {
	return e.Deputy != nil
}

func (e Employment) IsOnVacation() bool {
	return e.Vacation != nil && e.Vacation.IsActive()
}

func (e Employment) HasScheduledVacation() bool {
	return e.Vacation != nil && e.Vacation.IsScheduled()
}

func (e *Employment) AssignDeputy(deputyID int64) error {
	if deputyID == e.UserID {
		return fmt.Errorf("%w: user cannot be their own deputy", errors.ErrInvalidEmploymentDeputy)
	}

	if e.Deputy != nil && e.Deputy.ID == deputyID {
		return nil
	}

	deputy, err := newEmploymentDeputy(deputyID)
	if err != nil {
		return err
	}
	e.Deputy = &deputy

	e.recordEvent(EmploymentDeputyAssignedEvent{
		EmploymentID: e.ID,
		UserID:       e.UserID,
		DeputyID:     deputyID,
	})

	return nil
}

func (e *Employment) RemoveDeputy() {
	if e.Deputy == nil {
		return
	}

	e.Deputy = nil
	e.recordEvent(EmploymentDeputyRemovedEvent{
		EmploymentID: e.ID,
		UserID:       e.UserID,
	})
}

func (e *Employment) GrantVacation(endAt *time.Time) error {
	return e.ScheduleVacation(time.Now().UTC(), endAt)
}

func (e *Employment) ScheduleVacation(startsAt time.Time, endsAt *time.Time) error {
	vacation, err := newEmploymentVacation(startsAt, endsAt)
	if err != nil {
		return err
	}

	if e.Vacation != nil && e.Vacation.StartsAt.Equal(vacation.StartsAt) {
		if e.Vacation.EndsAt == nil && vacation.EndsAt == nil {
			return nil
		}
		if e.Vacation.EndsAt != nil && vacation.EndsAt != nil && e.Vacation.EndsAt.Equal(*vacation.EndsAt) {
			return nil
		}
	}

	e.Vacation = &vacation

	if e.HasScheduledVacation() {
		e.recordEvent(EmploymentVacationScheduledEvent{
			EmploymentID:     e.ID,
			UserID:           e.UserID,
			VacationStartsAt: e.Vacation.StartsAt,
			VacationEndsAt:   e.Vacation.EndsAt,
		})
		return nil
	}

	e.recordEvent(EmploymentVacationGrantedEvent{
		EmploymentID:     e.ID,
		UserID:           e.UserID,
		VacationStartsAt: e.Vacation.StartsAt,
		VacationEndsAt:   e.Vacation.EndsAt,
	})

	return nil
}

func (e *Employment) EndVacation() {
	if e.Vacation == nil {
		return
	}

	e.Vacation = nil
	e.recordEvent(EmploymentVacationEndedEvent{
		EmploymentID: e.ID,
		UserID:       e.UserID,
	})
}
