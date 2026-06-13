package database

import (
	"embed"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(logger *slog.Logger, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logger.Error("Failed to create migration driver", "error", err)
		return err
	}

	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Error("Failed to load migration files", "error", err)
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		logger.Error("Failed to create migrator", "error", err)
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("Failed to run migrations", "error", err)
		return err
	}

	logger.Info("Database migrations applied successfully")
	return nil
}
