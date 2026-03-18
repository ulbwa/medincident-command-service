package model

import (
	"time"

	"github.com/google/uuid"
)

// UserDomainEvent structs for CQRS routing decoupling the core domain from proto dtos.

type UserCreatedEvent struct {
	ID   int64
	Name UserName
}

type UserNameUpdatedEvent struct {
	ID   int64
	Name UserName
}

type UserCustomNameUpdatedEvent struct {
	ID         int64
	CustomName *UserName // nil if cleared
}

type UserGrantedAdminRoleEvent struct {
	ID        int64
	GrantedAt time.Time
	GrantedBy int64
}

type UserRevokedAdminRoleEvent struct {
	ID        int64
	RevokedAt time.Time
	RevokedBy int64
}

type UserEmployedEvent struct {
	UserID         int64
	EmploymentID   uuid.UUID
	OrganizationID uuid.UUID
	ClinicID       uuid.UUID
	DepartmentID   uuid.UUID
	Position       *string
	AssignedAt     time.Time
}

type UserDismissedEvent struct {
	UserID       int64
	EmploymentID uuid.UUID
}
