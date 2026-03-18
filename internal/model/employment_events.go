package model

import (
	"time"

	"github.com/google/uuid"
)

type EmploymentDeputyAssignedEvent struct {
	EmploymentID uuid.UUID
	UserID       int64
	DeputyID     int64
}

type EmploymentDeputyRemovedEvent struct {
	EmploymentID uuid.UUID
	UserID       int64
}

type EmploymentVacationGrantedEvent struct {
	EmploymentID     uuid.UUID
	UserID           int64
	VacationStartsAt time.Time
	VacationEndsAt   *time.Time
}

type EmploymentVacationScheduledEvent struct {
	EmploymentID     uuid.UUID
	UserID           int64
	VacationStartsAt time.Time
	VacationEndsAt   *time.Time
}

type EmploymentVacationEndedEvent struct {
	EmploymentID uuid.UUID
	UserID       int64
}
