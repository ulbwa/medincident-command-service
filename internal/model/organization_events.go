package model

import "github.com/google/uuid"

type OrganizationCreatedEvent struct {
	ID           uuid.UUID
	Name         string
	Description  *string
	LegalAddress Address
}

func (OrganizationCreatedEvent) EventType() string { return "organization.created" }

type OrganizationNameUpdatedEvent struct {
	ID   uuid.UUID
	Name string
}

func (OrganizationNameUpdatedEvent) EventType() string { return "organization.name_updated" }

type OrganizationDescriptionUpdatedEvent struct {
	ID          uuid.UUID
	Description *string
}

func (OrganizationDescriptionUpdatedEvent) EventType() string {
	return "organization.description_updated"
}

type OrganizationLegalAddressUpdatedEvent struct {
	ID           uuid.UUID
	LegalAddress Address
}

func (OrganizationLegalAddressUpdatedEvent) EventType() string {
	return "organization.legal_address_updated"
}
