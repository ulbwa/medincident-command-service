package zitadel

import (
	"context"
	"errors"
	"fmt"

	zitauth "github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	zitadelSDK "github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"github.com/ulbwa/medincident-command-service/internal/common/authorization"
)

var errInvalidOrInactiveToken = errors.New("invalid or inactive token")

// AccessTokenIntrospector implements the authorization.AccessTokenIntrospector port using the
// Zitadel OAuth2 token introspection endpoint authenticated via client credentials.
type AccessTokenIntrospector struct {
	authorizer *zitauth.Authorizer[*oauth.IntrospectionContext]
}

// NewAccessTokenIntrospector creates a AccessTokenIntrospector that contacts the Zitadel
// introspection endpoint at the given domain.
// clientID/clientSecret are the resource-server credentials registered in Zitadel.
// Optional zitadelSDK.Option values (e.g. zitadel.WithInsecure) can be passed for
// non-production environments.
func NewAccessTokenIntrospector(
	ctx context.Context,
	domain, clientID, clientSecret string,
	opts ...zitadelSDK.Option,
) (*AccessTokenIntrospector, error) {
	z := zitadelSDK.New(domain, opts...)

	authorizer, err := zitauth.New(
		ctx,
		z,
		oauth.WithIntrospection[*oauth.IntrospectionContext](
			oauth.ClientIDSecretIntrospectionAuthentication(clientID, clientSecret),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zitadel authorizer: %w", err)
	}

	return &AccessTokenIntrospector{authorizer: authorizer}, nil
}

// Introspect validates the bearer token via Zitadel and returns the caller's claims.
//
// Returns errInvalidOrInactiveToken when the token is invalid, expired, or inactive.
// Returns a wrapped infrastructure error when Zitadel is unavailable.
func (i *AccessTokenIntrospector) Introspect(ctx context.Context, bearerToken string) (*authorization.AccessClaims, error) {
	authCtx, err := i.authorizer.CheckAuthorization(ctx, "Bearer "+bearerToken)
	if err != nil {
		if errors.As(err, new(*zitauth.ServiceUnavailableErr)) {
			return nil, fmt.Errorf("identity provider unavailable: %w", err)
		}
		return nil, errInvalidOrInactiveToken
	}

	claims := &authorization.AccessClaims{
		ClientID:       authCtx.ClientID,
		Subject:        authCtx.Subject,
		Username:       authCtx.Username,
		Issuer:         authCtx.Issuer,
		TokenType:      authCtx.TokenType,
		TokenID:        authCtx.JWTID,
		OrganizationID: authCtx.OrganizationID(),
		Audience:       []string(authCtx.Audience),
		Scopes:         []string(authCtx.Scope),
		ExpiresAt:      authCtx.Expiration.AsTime(),
		IssuedAt:       authCtx.IssuedAt.AsTime(),
		NotBefore:      authCtx.NotBefore.AsTime(),
		AuthTime:       authCtx.AuthTime.AsTime(),
	}

	if !claims.IsValid() {
		return nil, errInvalidOrInactiveToken
	}

	return claims, nil
}

var _ authorization.AccessTokenIntrospector = (*AccessTokenIntrospector)(nil)
