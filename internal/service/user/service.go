package user

import (
	"errors"

	"github.com/ulbwa/medincident-command-service/internal/common/outbox"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
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
