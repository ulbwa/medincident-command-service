package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

type (
	CreateUserRequest struct {
		ID         int64
		GivenName  string
		FamilyName string
		MiddleName *string
	}

	CreateUserResponse struct {
		User *model.User
	}
)

func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	identity, err := s.identityProvider.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, errs.ErrIdentityNotFound) {
			return nil, errs.ErrIdentityNotFound
		}
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}

	if identity.Human == nil {
		return nil, errs.ErrIdentityNotHuman
	}

	userName, err := model.NewUserName(req.GivenName, req.FamilyName, req.MiddleName)
	if err != nil {
		return nil, err
	}

	var createdUser *model.User

	txCtx, tx, err := s.txFactory.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = persistence.WithinTransaction(txCtx, tx, func() error {
		userExists, err := s.repo.ExistsByID(txCtx, req.ID)
		if err != nil {
			return err
		}
		if userExists {
			return errs.ErrUserAlreadyExists
		}

		user, err := model.NewUser(req.ID, *userName)
		if err != nil {
			return err
		}

		if err := s.repo.Save(txCtx, user); err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}

		if err := s.eventDispatcher.Dispatch(txCtx, tx, user); err != nil {
			return fmt.Errorf("failed to dispatch domain events: %w", err)
		}

		createdUser = user
		return nil
	})
	if err != nil {
		return nil, err
	}

	bgCtx := context.WithoutCancel(ctx)
	go func() {
		syncCtx, cancel := context.WithTimeout(bgCtx, 2*time.Minute)
		defer cancel()

		if err := s.syncHumanIdentity(syncCtx, createdUser); err != nil {
			zerolog.Ctx(syncCtx).Error().Err(err).Int64("user_id", createdUser.ID).Msg("background profile sync to identity service failed")
			return
		}

		zerolog.Ctx(syncCtx).Debug().Int64("user_id", createdUser.ID).Msg("successfully synced user profile to identity service in background")
	}()

	return &CreateUserResponse{
		User: createdUser,
	}, nil
}
