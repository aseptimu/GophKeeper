package config

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
)

func MigrateDB(config *AppConfig) error {
	slog.Info("Running migration", "source path", config.MigrationsDir)

	migrationPath := config.MigrationsDir

	if config.MigrationsDir == "" {
		migrationPath = "./db/migrations" // Should run binary from source if MIGRATIONS_DIR env not specified
	}
	m, err := migrate.New(
		"file://"+migrationPath,
		config.DSN,
	)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("Nothing to migrate")
	}
	return nil
}
