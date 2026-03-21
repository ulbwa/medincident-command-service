package model

import "github.com/google/uuid"

type DepartmentCreatedEvent struct {
	ID          uuid.UUID
	ClinicID    uuid.UUID
	Name        string
	Description *string
}

func (DepartmentCreatedEvent) EventType() string { return "department.created" }

type DepartmentNameUpdatedEvent struct {
	ID   uuid.UUID
	Name string
}

func (DepartmentNameUpdatedEvent) EventType() string { return "department.name_updated" }

type DepartmentDescriptionUpdatedEvent struct {
	ID          uuid.UUID
	Description *string
}

func (DepartmentDescriptionUpdatedEvent) EventType() string { return "department.description_updated" }
