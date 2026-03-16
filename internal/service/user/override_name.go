package user

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"

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
		return errors.New("request is required")
	}

	customName, err := model.NewUserName(req.GivenName, req.FamilyName, req.MiddleName)
	if err != nil {
		return errors.New("invalid custom name: " + err.Error())
	}

	var updatedUser *model.User

	txCtx, tx, err := s.txFactory.Begin(ctx)
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}

	err = persistence.WithinTransaction(txCtx, tx, func() error {
		user, err := s.repo.GetByID(txCtx, req.UserID)
		if err != nil {
			return errors.New("failed to get user: " + err.Error())
		}
		if user == nil {
			return errors.New("user not found")
		}

		if err := user.OverrideName(*customName); err != nil {
			return errors.New("failed to override name: " + err.Error())
		}

		if err := s.repo.Save(txCtx, user); err != nil {
			return errors.New("failed to save user: " + err.Error())
		}

		if err := s.eventDispatcher.Dispatch(txCtx, tx, user); err != nil {
			return errors.New("failed to dispatch domain events: " + err.Error())
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
