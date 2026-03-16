package model

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
