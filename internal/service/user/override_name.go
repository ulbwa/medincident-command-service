package user

import (
	"context"
	"fmt"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

type OverrideUserNameRequest struct {
	UserID     int64
	GivenName  string
	FamilyName string
	MiddleName *string
}

func (s *Service) OverrideName(ctx context.Context, req *OverrideUserNameRequest) error {
	if req == nil {
		return errs.ErrInvalidRequest
	}

	customName, err := model.NewUserName(req.GivenName, req.FamilyName, req.MiddleName)
	if err != nil {
		return err
	}

	var updatedUser *model.User

	txCtx, tx, err := s.txFactory.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = persistence.WithinTransaction(txCtx, tx, func() error {
		user, err := s.repo.GetByID(txCtx, req.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if err := user.OverrideName(*customName); err != nil {
			return err
		}

		if err := s.repo.Save(txCtx, user); err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}

		if err := s.eventDispatcher.Dispatch(txCtx, tx, user); err != nil {
			return fmt.Errorf("failed to dispatch domain events: %w", err)
		}

		updatedUser = user
		return nil
	})
	if err != nil {
		return err
	}

	s.dispatchBackgroundIdentitySync(ctx, updatedUser, "override name")

	return nil
}
