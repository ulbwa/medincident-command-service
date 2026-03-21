package user

import (
	"context"

	"github.com/google/uuid"
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

// IdentityUserMetadata stores application-level metadata attached to an identity.
// UserID is optional: nil means this identity has not yet been linked to a registered user
// (e.g. the identity exists in the IdP but registration is incomplete or the identity
// belongs to a service account).
type IdentityUserMetadata struct {
	UserID *uuid.UUID
}

// Identity represents an account in the external identity provider.
// ID is an opaque string — its format depends on the concrete IdP (e.g. Zitadel snowflake,
// OIDC subject, etc.) and must not be interpreted by domain code.
type Identity struct {
	ID       string
	Human    *IdentityHuman
	Email    IdentityEmail
	IsActive bool
	Metadata *IdentityUserMetadata
}

// IdentityProvider is the port for interacting with the external identity provider.
// All IDs passed to this interface are opaque strings originating from the IdP.
type IdentityProvider interface {
	// Get retrieves an identity by its IdP-issued identifier.
	Get(ctx context.Context, identityID string) (*Identity, error)
	// UpdateHuman updates the human profile fields of an identity.
	UpdateHuman(ctx context.Context, identityID string, human *IdentityHuman) (*Identity, error)
	// UpdateUserMetadata stores application metadata on the identity, in particular the
	// application-level UserID once a user record has been created. This is called after
	// the user creation transaction commits so that the IdP reflects the linkage.
	UpdateUserMetadata(ctx context.Context, identityID string, metadata *IdentityUserMetadata) error
}
