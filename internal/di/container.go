package di

import (
	"github.com/samber/do/v2"

	"github.com/ulbwa/medincident-command-service/internal/config"
)

// NewContainer assembles the full dependency graph and returns the DI injector.
func NewContainer(cfg *config.Config) do.Injector {
	injector := do.New()

	// Config
	do.ProvideValue(injector, cfg)

	provideZerolog(injector)
	provideDatabase(injector)

	return injector
}
