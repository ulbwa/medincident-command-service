package model

import "github.com/google/uuid"

type ClinicCreatedEvent struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Description     *string
	PhysicalAddress Address
}

func (ClinicCreatedEvent) EventType() string { return "clinic.created" }

type ClinicNameUpdatedEvent struct {
	ID   uuid.UUID
	Name string
}

func (ClinicNameUpdatedEvent) EventType() string { return "clinic.name_updated" }

type ClinicDescriptionUpdatedEvent struct {
	ID          uuid.UUID
	Description *string
}

func (ClinicDescriptionUpdatedEvent) EventType() string { return "clinic.description_updated" }

type ClinicPhysicalAddressUpdatedEvent struct {
	ID              uuid.UUID
	PhysicalAddress Address
}

func (ClinicPhysicalAddressUpdatedEvent) EventType() string {
	return "clinic.physical_address_updated"
}
