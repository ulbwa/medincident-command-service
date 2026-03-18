package model

import "github.com/google/uuid"

type ClinicCreatedEvent struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Description     *string
	PhysicalAddress Address
}

type ClinicNameUpdatedEvent struct {
	ID   uuid.UUID
	Name string
}

type ClinicDescriptionUpdatedEvent struct {
	ID          uuid.UUID
	Description *string
}

type ClinicPhysicalAddressUpdatedEvent struct {
	ID              uuid.UUID
	PhysicalAddress Address
}
