package model

import (
	"github.com/google/uuid"

	"github.com/ulbwa/medincident-command-service/pkg/utils"
)

// Clinic represents a medical clinic within an organization.
type Clinic struct {
	Entity
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Description     *string
	PhysicalAddress Address
}

// NewClinic creates a new Clinic with validated data.
func NewClinic(id, organizationID uuid.UUID, name string, description *string, physicalAddress Address) (*Clinic, error) {
	c := &Clinic{
		ID:              id,
		OrganizationID:  organizationID,
		Name:            name,
		Description:     utils.PtrClone(description),
		PhysicalAddress: physicalAddress.copy(),
	}

	if err := validateClinic(c); err != nil {
		return nil, err
	}

	c.recordEvent(ClinicCreatedEvent{
		ID:              c.ID,
		OrganizationID:  c.OrganizationID,
		Name:            c.Name,
		Description:     c.Description,
		PhysicalAddress: c.PhysicalAddress,
	})

	return c, nil
}

// RestoreClinic reconstructs a Clinic from persistent storage.
func RestoreClinic(id, organizationID uuid.UUID, name string, description *string, physicalAddress Address) (*Clinic, error) {
	c := &Clinic{
		ID:              id,
		OrganizationID:  organizationID,
		Name:            name,
		Description:     utils.PtrClone(description),
		PhysicalAddress: physicalAddress.copy(),
	}

	if err := validateClinic(c); err != nil {
		return nil, err
	}

	return c, nil
}

// UpdateName changes the clinic name.
func (c *Clinic) UpdateName(name string) error {
	if c.Name == name {
		return nil
	}

	if err := validateClinicName(name); err != nil {
		return err
	}

	c.Name = name
	c.recordEvent(ClinicNameUpdatedEvent{
		ID:   c.ID,
		Name: name,
	})
	return nil
}

// UpdateDescription changes the clinic description.
func (c *Clinic) UpdateDescription(description string) error {
	if c.Description != nil && *c.Description == description {
		return nil
	}

	if err := validateClinicDescription(description); err != nil {
		return err
	}

	c.Description = &description
	c.recordEvent(ClinicDescriptionUpdatedEvent{
		ID:          c.ID,
		Description: c.Description,
	})
	return nil
}

// RemoveDescription removes the clinic description.
func (c *Clinic) RemoveDescription() {
	if c.Description == nil {
		return
	}

	c.Description = nil
	c.recordEvent(ClinicDescriptionUpdatedEvent{
		ID:          c.ID,
		Description: nil,
	})
}

// UpdatePhysicalAddress changes the clinic physical address.
func (c *Clinic) UpdatePhysicalAddress(address Address) error {
	if c.PhysicalAddress.Equals(address) {
		return nil
	}

	if err := validateAddress(address); err != nil {
		return err
	}

	c.PhysicalAddress = address.copy()
	c.recordEvent(ClinicPhysicalAddressUpdatedEvent{
		ID:              c.ID,
		PhysicalAddress: c.PhysicalAddress,
	})
	return nil
}
