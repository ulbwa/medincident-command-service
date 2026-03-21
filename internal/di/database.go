package di

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver registration.
	"github.com/samber/do/v2"

	"github.com/ulbwa/medincident-command-service/internal/common/outbox"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
	"github.com/ulbwa/medincident-command-service/internal/config"
	pginfra "github.com/ulbwa/medincident-command-service/internal/repository/postgres"
	userrepo "github.com/ulbwa/medincident-command-service/internal/repository/user"
	serviceuser "github.com/ulbwa/medincident-command-service/internal/service/user"
)

func provideDatabase(injector do.Injector) {
	// Database pool
	do.Provide(injector, func(i do.Injector) (*sqlx.DB, error) {
		cfg := do.MustInvoke[*config.Config](i)
		db, err := sqlx.Connect("postgres", cfg.Database.DSN)
		if err != nil {
			return nil, fmt.Errorf("connect to postgres: %w", err)
		}

		return db, nil
	})

	// Transaction factory
	do.Provide(injector, func(i do.Injector) (persistence.TransactionFactory, error) {
		db := do.MustInvoke[*sqlx.DB](i)
		return pginfra.NewTxFactory(db), nil
	})

	// Outbox event dispatcher
	do.Provide(injector, func(_ do.Injector) (outbox.EventDispatcher, error) {
		return pginfra.NewOutboxDispatcher(), nil
	})

	// User write repository
	do.Provide(injector, func(i do.Injector) (serviceuser.Repository, error) {
		db := do.MustInvoke[*sqlx.DB](i)
		return userrepo.NewRepository(db), nil
	})
}
