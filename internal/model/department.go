package model

import (
	"github.com/google/uuid"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/pkg/utils"
)

// Department represents a department within a clinic.
type Department struct {
	Entity
	ID          uuid.UUID
	ClinicID    uuid.UUID
	Name        string
	Description *string
}

// NewDepartment creates a new Department with validated data.
func NewDepartment(id, clinicID uuid.UUID, name string, description *string) (*Department, error) {
	d := &Department{
		ID:          id,
		ClinicID:    clinicID,
		Name:        name,
		Description: utils.PtrClone(description),
	}

	if err := validateDepartment(d); err != nil {
		return nil, err
	}

	d.recordEvent(DepartmentCreatedEvent{
		ID:          d.ID,
		ClinicID:    d.ClinicID,
		Name:        d.Name,
		Description: d.Description,
	})

	return d, nil
}

// RestoreDepartment reconstructs a Department from persistent storage.
func RestoreDepartment(id, clinicID uuid.UUID, name string, description *string) (*Department, error) {
	d := &Department{
		ID:          id,
		ClinicID:    clinicID,
		Name:        name,
		Description: utils.PtrClone(description),
	}

	if err := validateDepartment(d); err != nil {
		return nil, err
	}

	return d, nil
}

// UpdateName changes the department name.
func (d *Department) UpdateName(name string) error {
	if d.Name == name {
		return nil
	}

	if err := validateDepartmentName(name); err != nil {
		return errs.NewInvalidDepartmentError(errs.DepartmentFieldName, err)
	}

	d.Name = name
	d.recordEvent(DepartmentNameUpdatedEvent{
		ID:   d.ID,
		Name: name,
	})
	return nil
}

// SetDescription sets the department description.
func (d *Department) SetDescription(description string) error {
	if d.Description != nil && *d.Description == description {
		return nil
	}

	if err := validateDepartmentDescription(description); err != nil {
		return errs.NewInvalidDepartmentError(errs.DepartmentFieldDescription, err)
	}

	d.Description = &description
	d.recordEvent(DepartmentDescriptionUpdatedEvent{
		ID:          d.ID,
		Description: d.Description,
	})
	return nil
}

// RemoveDescription removes the department description.
func (d *Department) RemoveDescription() {
	if d.Description == nil {
		return
	}

	d.Description = nil
	d.recordEvent(DepartmentDescriptionUpdatedEvent{
		ID:          d.ID,
		Description: nil,
	})
}
