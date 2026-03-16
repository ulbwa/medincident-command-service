package user

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"

	"github.com/ulbwa/medincident-command-service/internal/model"
	"github.com/ulbwa/medincident-command-service/pkg/utils"
)

func (s *Service) syncHumanIdentity(ctx context.Context, user *model.User) error {
	identity, err := s.identityProvider.Get(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get identity for user %d: %w", user.ID, err)
	}

	if identity.Human == nil {
		return fmt.Errorf("identity profile for user %d is not a human, cannot sync", user.ID)
	}

	preferredName := user.PreferredName()
	humanUpdate := &IdentityHuman{
		GivenName:         preferredName.GivenName,
		FamilyName:        preferredName.FamilyName,
		NickName:          utils.Ptr(preferredName.ShortName()),
		DisplayName:       preferredName.DisplayName(),
		Gender:            identity.Human.Gender,            // keep old
		PreferredLanguage: identity.Human.PreferredLanguage, // keep old
	}

	b := backoff.NewExponentialBackOff()

	operation := func() error {
		_, updateErr := s.identityProvider.UpdateHuman(ctx, user.ID, humanUpdate)
		if updateErr != nil {
			return fmt.Errorf("update API call failed: %w", updateErr)
		}
		return nil
	}

	err = backoff.Retry(operation, backoff.WithContext(b, ctx))
	if err != nil {
		return fmt.Errorf("failed to sync profile to identity service after multiple attempts: %w", err)
	}

	return nil
}
