package authorization

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessClaims_IsValid(t *testing.T) {
	t.Parallel()

	now := time.Now()

	base := AccessClaims{
		ClientID:  "custom-ui-client",
		Subject:   "user-123",
		Issuer:    "https://issuer.example.com",
		TokenType: "Bearer",
		Scopes:    []string{ScopeUserCreate},
		ExpiresAt: now.Add(10 * time.Minute),
		IssuedAt:  now.Add(-2 * time.Minute),
		NotBefore: now.Add(-1 * time.Minute),
		AuthTime:  now.Add(-3 * time.Minute),
	}

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		claims := base
		assert.True(t, claims.IsValid())
	})

	t.Run("NilClaims", func(t *testing.T) {
		t.Parallel()

		var claims *AccessClaims
		assert.False(t, claims.IsValid())
	})

	t.Run("Expired", func(t *testing.T) {
		t.Parallel()

		claims := base
		claims.ExpiresAt = now.Add(-time.Second)
		assert.False(t, claims.IsValid())
	})

	t.Run("NotYetValid", func(t *testing.T) {
		t.Parallel()

		claims := base
		claims.NotBefore = now.Add(time.Minute)
		assert.False(t, claims.IsValid())
	})

	t.Run("WrongTokenType", func(t *testing.T) {
		t.Parallel()

		claims := base
		claims.TokenType = "mac"
		assert.True(t, claims.IsValid())
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		t.Parallel()

		claims := base
		claims.ClientID = ""
		assert.True(t, claims.IsValid())
	})

	t.Run("MissingScopes", func(t *testing.T) {
		t.Parallel()

		claims := base
		claims.Scopes = nil
		assert.True(t, claims.IsValid())
	})

	t.Run("MissingExpiration", func(t *testing.T) {
		t.Parallel()

		claims := base
		claims.ExpiresAt = time.Time{}
		assert.False(t, claims.IsValid())
	})
}
