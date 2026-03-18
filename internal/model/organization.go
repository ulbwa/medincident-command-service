package model

import (
	"github.com/google/uuid"
)

// Organization represents a legal organization (aggregate root).
type Organization struct {
	Entity
	ID           uuid.UUID
	Name         string
	Description  *string
	LegalAddress Address
}

// NewOrganization creates a new Organization with validated data.
func NewOrganization(id uuid.UUID, name string, description *string, legalAddress Address) (*Organization, error) {
	o := &Organization{
		ID:           id,
		Name:         name,
		Description:  description,
		LegalAddress: legalAddress,
	}

	if err := validateOrganization(o); err != nil {
		return nil, err
	}

	o.recordEvent(OrganizationCreatedEvent{
		ID:           o.ID,
		Name:         o.Name,
		Description:  o.Description,
		LegalAddress: o.LegalAddress,
	})

	return o, nil
}

// RestoreOrganization reconstructs an Organization from persistent storage.
func RestoreOrganization(id uuid.UUID, name string, description *string, legalAddress Address) (*Organization, error) {
	o := &Organization{
		ID:           id,
		Name:         name,
		Description:  description,
		LegalAddress: legalAddress,
	}

	if err := validateOrganization(o); err != nil {
		return nil, err
	}

	return o, nil
}

// UpdateName changes the organization name.
func (o *Organization) UpdateName(name string) error {
	if o.Name == name {
		return nil
	}

	if err := validateOrganizationName(name); err != nil {
		return err
	}

	o.Name = name
	o.recordEvent(OrganizationNameUpdatedEvent{
		ID:   o.ID,
		Name: name,
	})
	return nil
}

// UpdateDescription changes the organization description.
func (o *Organization) UpdateDescription(description string) error {
	if o.Description != nil && *o.Description == description {
		return nil
	}

	if err := validateOrganizationDescription(description); err != nil {
		return err
	}

	o.Description = &description
	o.recordEvent(OrganizationDescriptionUpdatedEvent{
		ID:          o.ID,
		Description: o.Description,
	})
	return nil
}

// RemoveDescription removes the organization description.
func (o *Organization) RemoveDescription() {
	if o.Description == nil {
		return
	}

	o.Description = nil
	o.recordEvent(OrganizationDescriptionUpdatedEvent{
		ID:          o.ID,
		Description: nil,
	})
}

// UpdateLegalAddress changes the organization legal address.
func (o *Organization) UpdateLegalAddress(address Address) error {
	if o.LegalAddress.Equals(address) {
		return nil
	}

	o.LegalAddress = address
	o.recordEvent(OrganizationLegalAddressUpdatedEvent{
		ID:           o.ID,
		LegalAddress: address,
	})
	return nil
}
