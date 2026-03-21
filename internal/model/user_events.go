package model

import (
	"time"

	"github.com/google/uuid"
)

// UserDomainEvent structs for CQRS routing decoupling the core domain from proto dtos.

type UserCreatedEvent struct {
	ID         uuid.UUID
	IdentityID string
	Name       UserName
}

type UserNameUpdatedEvent struct {
	ID   uuid.UUID
	Name UserName
}

type UserCustomNameUpdatedEvent struct {
	ID         uuid.UUID
	CustomName *UserName // nil if cleared
}

type UserGrantedAdminRoleEvent struct {
	ID        uuid.UUID
	GrantedAt time.Time
	GrantedBy uuid.UUID
}

type UserRevokedAdminRoleEvent struct {
	ID        uuid.UUID
	RevokedAt time.Time
	RevokedBy uuid.UUID
}

type UserEmployedEvent struct {
	UserID         uuid.UUID
	EmploymentID   uuid.UUID
	OrganizationID uuid.UUID
	ClinicID       uuid.UUID
	DepartmentID   uuid.UUID
	Position       *string
	AssignedAt     time.Time
}

type UserDismissedEvent struct {
	UserID       uuid.UUID
	EmploymentID uuid.UUID
}
