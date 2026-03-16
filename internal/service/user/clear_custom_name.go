package user

import (
	"context"
	"fmt"

	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

// ClearCustomName resets user's custom name matching business invariant logic
func (s *Service) ClearCustomName(ctx context.Context, userID int64) error {
	var syncedUser *model.User

	txCtx, tx, err := s.txFactory.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = persistence.WithinTransaction(txCtx, tx, func() error {
		u, err := s.repo.GetByID(txCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if err := u.ClearCustomName(); err != nil {
			return err
		}

		if err := s.repo.Save(txCtx, u); err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}

		if err := s.eventDispatcher.Dispatch(txCtx, tx, u); err != nil {
			return fmt.Errorf("failed to dispatch domain events: %w", err)
		}

		syncedUser = u
		return nil
	})
	if err != nil {
		return err
	}

	s.dispatchBackgroundIdentitySync(ctx, syncedUser, "clear custom name")

	return nil
}
