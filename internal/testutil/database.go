package testutil

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/ibbycs/private-markets-api/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type noopLogger struct{}

func (noopLogger) Printf(format string, v ...interface{}) {}

func NewTestDatabase(t *testing.T, logger *slog.Logger) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithLogger(noopLogger{}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	databaseUrl, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	pool, err := database.NewPostgresPool(ctx, databaseUrl)

	if err != nil {
		t.Fatal(err)
	}

	err = database.Migrate(logger, pool)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if pool != nil {
			pool.Close()
		}

		container.Container.Terminate(ctx)
	})

	return pool
}
