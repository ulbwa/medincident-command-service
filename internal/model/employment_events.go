package model

import (
	"time"

	"github.com/google/uuid"
)

type EmploymentDeputyAssignedEvent struct {
	EmploymentID uuid.UUID
	DeputyID     uuid.UUID
}

type EmploymentDeputyRemovedEvent struct {
	EmploymentID uuid.UUID
}

type EmploymentVacationGrantedEvent struct {
	EmploymentID     uuid.UUID
	VacationStartsAt time.Time
	VacationEndsAt   *time.Time
}

type EmploymentVacationScheduledEvent struct {
	EmploymentID     uuid.UUID
	VacationStartsAt time.Time
	VacationEndsAt   *time.Time
}

type EmploymentVacationEndedEvent struct {
	EmploymentID uuid.UUID
}
