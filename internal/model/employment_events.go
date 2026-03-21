package model

import (
	"time"

	"github.com/google/uuid"
)

type EmploymentDeputyAssignedEvent struct {
	EmploymentID uuid.UUID
	DeputyID     uuid.UUID
}

func (EmploymentDeputyAssignedEvent) EventType() string { return "employment.deputy_assigned" }

type EmploymentDeputyRemovedEvent struct {
	EmploymentID uuid.UUID
}

func (EmploymentDeputyRemovedEvent) EventType() string { return "employment.deputy_removed" }

type EmploymentVacationGrantedEvent struct {
	EmploymentID     uuid.UUID
	VacationStartsAt time.Time
	VacationEndsAt   *time.Time
}

func (EmploymentVacationGrantedEvent) EventType() string { return "employment.vacation_granted" }

type EmploymentVacationScheduledEvent struct {
	EmploymentID     uuid.UUID
	VacationStartsAt time.Time
	VacationEndsAt   *time.Time
}

func (EmploymentVacationScheduledEvent) EventType() string { return "employment.vacation_scheduled" }

type EmploymentVacationEndedEvent struct {
	EmploymentID uuid.UUID
}

func (EmploymentVacationEndedEvent) EventType() string { return "employment.vacation_ended" }
