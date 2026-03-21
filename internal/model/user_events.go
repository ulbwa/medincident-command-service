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

func (UserCreatedEvent) EventType() string { return "user.created" }

type UserNameUpdatedEvent struct {
	ID   uuid.UUID
	Name UserName
}

func (UserNameUpdatedEvent) EventType() string { return "user.name_updated" }

type UserCustomNameUpdatedEvent struct {
	ID         uuid.UUID
	CustomName *UserName // nil if cleared
}

func (UserCustomNameUpdatedEvent) EventType() string { return "user.custom_name_updated" }

type UserGrantedAdminRoleEvent struct {
	ID        uuid.UUID
	GrantedAt time.Time
	GrantedBy uuid.UUID
}

func (UserGrantedAdminRoleEvent) EventType() string { return "user.admin_role_granted" }

type UserRevokedAdminRoleEvent struct {
	ID        uuid.UUID
	RevokedAt time.Time
	RevokedBy uuid.UUID
}

func (UserRevokedAdminRoleEvent) EventType() string { return "user.admin_role_revoked" }

type UserEmployedEvent struct {
	UserID         uuid.UUID
	EmploymentID   uuid.UUID
	OrganizationID uuid.UUID
	ClinicID       uuid.UUID
	DepartmentID   uuid.UUID
	Position       *string
	AssignedAt     time.Time
}

func (UserEmployedEvent) EventType() string { return "user.employed" }

type UserDismissedEvent struct {
	UserID       uuid.UUID
	EmploymentID uuid.UUID
}

func (UserDismissedEvent) EventType() string { return "user.dismissed" }
