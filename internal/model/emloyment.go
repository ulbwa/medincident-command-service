package model

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

type Employment struct {
	UserID        int64
	DepartmentID  uuid.UUID
	Position      *string
	AssignedAt    time.Time
	DeputyID      *int64
	IsOnVacation  bool
	VacationEndAt *time.Time
}

func validateEmploymentPosition(position *string) error {
	if position == nil {
		return nil
	}
	positionLen := utf8.RuneCountInString(*position)
	if positionLen < 1 {
		return fmt.Errorf("%w: too short (min 1)", errs.ErrInvalidEmploymentPosition)
	}
	if positionLen > 255 {
		return fmt.Errorf("%w: too long (max 255)", errs.ErrInvalidEmploymentPosition)
	}
	return nil
}

func validateDeputyID(userID int64, deputyID *int64) error {
	if deputyID != nil && *deputyID == userID {
		return fmt.Errorf("%w: user cannot be their own deputy", errs.ErrInvalidEmploymentDeputy)
	}
	return nil
}

func validateEmployment(e Employment) error {
	if err := validateUserID(e.UserID); err != nil {
		return err
	}
	if err := validateDepartmentID(e.DepartmentID); err != nil {
		return err
	}
	if err := validateEmploymentPosition(e.Position); err != nil {
		return err
	}
	if err := validateDeputyID(e.UserID, e.DeputyID); err != nil {
		return err
	}
	if !e.IsOnVacation && e.VacationEndAt != nil {
		return fmt.Errorf("%w: end date must be nil if employee is not on vacation", errs.ErrInvalidEmploymentVacation)
	}
	return nil
}

func NewEmployment(userID int64, departmentID uuid.UUID, position *string, assignedAt time.Time) (*Employment, error) {
	e := &Employment{
		UserID:       userID,
		DepartmentID: departmentID,
		Position:     position,
		AssignedAt:   assignedAt,
	}
	if err := validateEmployment(*e); err != nil {
		return nil, err
	}
	return e, nil
}

func RestoreEmployment(userID int64, departmentID uuid.UUID, position *string, assignedAt time.Time, deputyID *int64, isOnVacation bool, vacationEndAt *time.Time) (*Employment, error) {
	e := &Employment{
		UserID:        userID,
		DepartmentID:  departmentID,
		Position:      position,
		AssignedAt:    assignedAt,
		DeputyID:      deputyID,
		IsOnVacation:  isOnVacation,
		VacationEndAt: vacationEndAt,
	}
	if err := validateEmployment(*e); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Employment) GrantVacation(endAt *time.Time) {
	e.IsOnVacation = true
	e.VacationEndAt = endAt
}

func (e *Employment) EndVacation() {
	e.IsOnVacation = false
	e.VacationEndAt = nil
}

func (e *Employment) AssignDeputy(deputyID int64) error {
	if err := validateDeputyID(e.UserID, &deputyID); err != nil {
		return err
	}
	e.DeputyID = &deputyID
	return nil
}

func (e *Employment) RemoveDeputy() {
	e.DeputyID = nil
}
