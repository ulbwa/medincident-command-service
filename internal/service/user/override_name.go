package user

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

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
		if user == nil {
			return errs.ErrUserNotFound
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

	bgCtx := context.WithoutCancel(ctx)
	go func() {
		syncCtx, cancel := context.WithTimeout(bgCtx, 2*time.Minute)
		defer cancel()

		if err := s.syncHumanIdentity(syncCtx, updatedUser); err != nil {
			zerolog.Ctx(syncCtx).Error().Err(err).Int64("user_id", updatedUser.ID).Msg("background profile sync to identity service failed after name override")
			return
		}

		zerolog.Ctx(syncCtx).Debug().Int64("user_id", updatedUser.ID).Msg("successfully synced user profile to identity service in background after name override")
	}()

	return nil
}
