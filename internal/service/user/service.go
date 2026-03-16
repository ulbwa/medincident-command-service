package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/common/outbox"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
	"github.com/ulbwa/medincident-command-service/internal/model"
	"github.com/ulbwa/medincident-command-service/pkg/utils"
)

type Service struct {
	txFactory        persistence.TransactionFactory
	eventDispatcher  outbox.EventDispatcher
	identityProvider IdentityProvider
	repo             Repository
}

func NewService(txFactory persistence.TransactionFactory, eventDispatcher outbox.EventDispatcher, identityProvider IdentityProvider, repo Repository) (*Service, error) {
	if txFactory == nil {
		return nil, errors.New("transaction factory is required")
	}
	if eventDispatcher == nil {
		return nil, errors.New("event dispatcher is required")
	}
	if identityProvider == nil {
		return nil, errors.New("identity provider is required")
	}
	if repo == nil {
		return nil, errors.New("repository is required")
	}
	return &Service{
		txFactory:        txFactory,
		eventDispatcher:  eventDispatcher,
		identityProvider: identityProvider,
		repo:             repo,
	}, nil
}

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

		user, err := model.NewUser(req.ID, userName)
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

	s.dispatchBackgroundIdentitySync(ctx, createdUser, "create user")

	return &CreateUserResponse{
		User: createdUser,
	}, nil
}

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

		if err := user.OverrideName(customName); err != nil {
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

func (s *Service) dispatchBackgroundIdentitySync(ctx context.Context, user *model.User, operation string) {
	bgCtx := context.WithoutCancel(ctx)
	go func() {
		syncCtx, cancel := context.WithTimeout(bgCtx, 2*time.Minute)
		defer cancel()

		if err := s.syncHumanIdentity(syncCtx, user); err != nil {
			zerolog.Ctx(syncCtx).Error().
				Err(err).
				Int64("user_id", user.ID).
				Str("operation", operation).
				Msg("background profile sync to identity service failed")
			return
		}

		zerolog.Ctx(syncCtx).Debug().
			Int64("user_id", user.ID).
			Str("operation", operation).
			Msg("successfully synced user profile to identity service in background")
	}()
}

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
