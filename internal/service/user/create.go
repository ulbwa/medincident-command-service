package user

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"

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
		return nil, errors.New("request is required")
	}

	identity, err := s.identityProvider.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, ErrIdentityNotFound) {
			return nil, errors.New("identity with the given ID does not exist")
		}
		return nil, err
	}

	if identity.Human == nil {
		return nil, errors.New("identity profile must be a human to create a user")
	}

	userName, err := model.NewUserName(req.GivenName, req.FamilyName, req.MiddleName)
	if err != nil {
		return nil, errors.New("invalid user name: " + err.Error())
	}

	var createdUser *model.User

	txCtx, tx, err := s.txFactory.Begin(ctx)
	if err != nil {
		return nil, errors.New("failed to begin transaction: " + err.Error())
	}

	err = persistence.WithinTransaction(txCtx, tx, func() error {
		userExists, err := s.repo.ExistsByID(txCtx, req.ID)
		if err != nil {
			return err
		}
		if userExists {
			return errors.New("user with the given ID already exists")
		}

		user, err := model.NewUser(req.ID, *userName)
		if err != nil {
			return errors.New("invalid user: " + err.Error())
		}

		if err := s.repo.Save(txCtx, user); err != nil {
			return errors.New("failed to save user: " + err.Error())
		}

		if err := s.eventDispatcher.Dispatch(txCtx, tx, user); err != nil {
			return errors.New("failed to dispatch domain events: " + err.Error())
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
