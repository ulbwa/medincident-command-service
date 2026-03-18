package model

import "github.com/google/uuid"

type OrganizationCreatedEvent struct {
	ID           uuid.UUID
	Name         string
	Description  *string
	LegalAddress Address
}

type OrganizationNameUpdatedEvent struct {
	ID   uuid.UUID
	Name string
}

type OrganizationDescriptionUpdatedEvent struct {
	ID          uuid.UUID
	Description *string
}

type OrganizationLegalAddressUpdatedEvent struct {
	ID           uuid.UUID
	LegalAddress Address
}
