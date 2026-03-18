package model

import "github.com/google/uuid"

type DepartmentCreatedEvent struct {
	ID          uuid.UUID
	ClinicID    uuid.UUID
	Name        string
	Description *string
}

type DepartmentNameUpdatedEvent struct {
	ID   uuid.UUID
	Name string
}

type DepartmentDescriptionUpdatedEvent struct {
	ID          uuid.UUID
	Description *string
}
