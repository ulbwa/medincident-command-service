package user

import (
	"context"
	"errors"
)

//go:generate go-enum -f=$GOFILE --marshal

// IdentityGender represents the gender of a human identity.
// ENUM(
// Unspecified="GENDER_UNSPECIFIED"
// Female="GENDER_FEMALE"
// Male="GENDER_MALE"
// Diverse="GENDER_DIVERSE"
// )
type IdentityGender string

type IdentityEmail struct {
	Address    string
	IsVerified bool
}

type IdentityHuman struct {
	GivenName         string
	FamilyName        string
	NickName          *string
	DisplayName       string
	Gender            IdentityGender
	PreferredLanguage *string
}

type Identity struct {
	ID       int64
	Human    *IdentityHuman
	Email    IdentityEmail
	IsActive bool
}

var ErrIdentityNotFound = errors.New("identity not found")

type IdentityProvider interface {
	Get(ctx context.Context, id int64) (*Identity, error)
	UpdateHuman(ctx context.Context, id int64, human *IdentityHuman) (*Identity, error)
}
