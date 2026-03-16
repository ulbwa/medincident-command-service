package model

import "time"

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
