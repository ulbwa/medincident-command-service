package model

import (
	"time"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/pkg/utils"
)

// EmploymentDeputy holds the user ID of the designated deputy.
// Deputy delegation is non-transitive: only the direct deputy receives delegated authority.
type EmploymentDeputy struct {
	ID uuid.UUID
}

func (d EmploymentDeputy) copy() *EmploymentDeputy {
	cloned := d
	return &cloned
}

func NewEmploymentDeputy(deputyID uuid.UUID) (EmploymentDeputy, error) {
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

func NewEmploymentVacation(startsAt time.Time, endsAt *time.Time) (EmploymentVacation, error) {
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
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	ClinicID       uuid.UUID
	DepartmentID   uuid.UUID
	Position       *string
	AssignedAt     time.Time
	Deputy         *EmploymentDeputy
	Vacation       *EmploymentVacation
}

const EmploymentVacationMaxScheduleAheadMonths = 6

func (e *Employment) copy() *Employment {
	copied := &Employment{
		ID:             e.ID,
		UserID:         e.UserID,
		OrganizationID: e.OrganizationID,
		ClinicID:       e.ClinicID,
		DepartmentID:   e.DepartmentID,
		Position:       utils.PtrClone(e.Position),
		AssignedAt:     e.AssignedAt,
	}
	if e.Deputy != nil {
		copied.Deputy = e.Deputy.copy()
	}
	if e.Vacation != nil {
		copied.Vacation = e.Vacation.copy()
	}
	return copied
}

func NewEmployment(userID, organizationID, clinicID, departmentID uuid.UUID, position *string, assignedAt time.Time) (*Employment, error) {
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

func RestoreEmployment(id, userID, organizationID, clinicID, departmentID uuid.UUID, position *string, assignedAt time.Time, deputy *EmploymentDeputy, vacation *EmploymentVacation) (*Employment, error) {
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

func (e *Employment) AssignDeputy(deputyID uuid.UUID) error {
	if deputyID == e.UserID {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldDeputy, errs.ErrUserCannotBeOwnDeputy)
	}

	if e.Deputy != nil && e.Deputy.ID == deputyID {
		return nil
	}

	deputy, err := NewEmploymentDeputy(deputyID)
	if err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldDeputy, err)
	}
	e.Deputy = &deputy

	e.recordEvent(EmploymentDeputyAssignedEvent{
		EmploymentID: e.ID,
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
	})
}

func (e *Employment) GrantVacation(endsAt *time.Time) error {
	now := time.Now().UTC()
	return e.scheduleVacation(now, now, endsAt)
}

func (e *Employment) ScheduleVacation(startsAt time.Time, endsAt *time.Time) error {
	now := time.Now().UTC()

	if startsAt.IsZero() {
		return errs.NewInvalidEmploymentError(
			errs.EmploymentFieldVacation,
			errs.NewInvalidEmploymentVacationError(errs.EmploymentVacationFieldStartsAt, errs.NewValueRequiredError()),
		)
	}

	if err := validateTimestampNotBefore(startsAt, now); err != nil {
		return errs.NewInvalidEmploymentError(
			errs.EmploymentFieldVacation,
			errs.NewInvalidEmploymentVacationError(errs.EmploymentVacationFieldStartsAt, err),
		)
	}

	return e.scheduleVacation(now, startsAt, endsAt)
}

func (e *Employment) scheduleVacation(now, startsAt time.Time, endsAt *time.Time) error {
	if e.Vacation != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldVacation, errs.ErrEmploymentVacationAlreadyExists)
	}

	maximum := now.AddDate(0, EmploymentVacationMaxScheduleAheadMonths, 0)
	if err := validateTimestampNotAfter(startsAt, maximum); err != nil {
		return errs.NewInvalidEmploymentError(
			errs.EmploymentFieldVacation,
			errs.NewInvalidEmploymentVacationError(errs.EmploymentVacationFieldStartsAt, err),
		)
	}

	vacation, err := NewEmploymentVacation(startsAt, endsAt)
	if err != nil {
		return errs.NewInvalidEmploymentError(errs.EmploymentFieldVacation, err)
	}

	e.Vacation = &vacation

	if e.HasScheduledVacation() {
		e.recordEvent(EmploymentVacationScheduledEvent{
			EmploymentID:     e.ID,
			VacationStartsAt: e.Vacation.StartsAt,
			VacationEndsAt:   e.Vacation.EndsAt,
		})
		return nil
	}

	e.recordEvent(EmploymentVacationGrantedEvent{
		EmploymentID:     e.ID,
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
	})
}
