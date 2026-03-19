package authorization

import (
	"context"
	"time"
)

// AccessClaims holds validated access-token claims used for application authorization.
type AccessClaims struct {
	// ClientID is the client_id of the application the token was issued to.
	ClientID string
	// Subject is the sub claim — unique identifier of the resource owner (user or service account).
	Subject string
	// Username is ZITADEL's login name of the user. Consists of username@primarydomain.
	Username string
	// Issuer is the iss claim — issuer of the token.
	Issuer string
	// TokenType is the token_type claim. Value is always Bearer.
	TokenType string
	// TokenID is the jti claim — unique identifier of the token.
	TokenID string
	// OrganizationID is the ZITADEL-specific resource-owner organization identifier.
	OrganizationID string
	// Audience is the aud claim — intended recipients of the token.
	Audience []string
	// Scopes is the scope claim — space-delimited list of scopes granted to the token.
	Scopes []string
	// ExpiresAt is the exp claim — time the token expires.
	ExpiresAt time.Time
	// IssuedAt is the iat claim — time the token was issued at.
	IssuedAt time.Time
	// NotBefore is the nbf claim — time before which the token must not be used.
	NotBefore time.Time
	// AuthTime is the auth_time claim — time when the user authentication occurred.
	AuthTime time.Time
}

// HasScope reports whether the token contains the specified OAuth scope.
func (c *AccessClaims) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// IsValid reports whether token claims pass application-level validity checks at current time.
func (c *AccessClaims) IsValid() bool {
	if c == nil {
		return false
	}

	now := time.Now()
	return c.NotBefore.Before(now) && now.Before(c.ExpiresAt)
}

// AccessTokenIntrospector validates bearer access tokens against identity provider
// and returns claims required for authorization decisions.
type AccessTokenIntrospector interface {
	Introspect(ctx context.Context, bearerToken string) (*AccessClaims, error)
}
