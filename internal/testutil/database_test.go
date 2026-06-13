//go:build unit

package testutil_test

import (
	"testing"

	"github.com/ibbycs/private-markets-api/internal/logger"
	"github.com/ibbycs/private-markets-api/internal/testutil"
)

func TestDatabaseSetup(t *testing.T) {
	logger := logger.New()
	pool := testutil.NewTestDatabase(t, logger)

	if err := pool.Ping(t.Context()); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Log("successfully connected to test database")
}
