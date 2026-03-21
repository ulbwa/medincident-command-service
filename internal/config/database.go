package config

// DatabaseConfig holds the PostgreSQL connection parameters.
type DatabaseConfig struct {
	// DSN is a Go-driver-compatible PostgreSQL connection string.
	// Example: postgres://user:pass@localhost:5432/dbname?sslmode=disable
	// Populated from the DATABASE_URL environment variable via config.yaml.
	DSN string `yaml:"dsn" validate:"required"`
}
