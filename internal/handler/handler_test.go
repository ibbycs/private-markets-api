//go:build integration

package handler

import (
	"io"
	"log/slog"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/ibbycs/private-markets-api/internal/repository"
	"github.com/ibbycs/private-markets-api/internal/service"
	"github.com/ibbycs/private-markets-api/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

type testHarness struct {
	api  humatest.TestAPI
	pool *pgxpool.Pool
}

func newTestHarness(t *testing.T) *testHarness {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	api := testutil.NewTestAPI(t)
	pool := testutil.NewTestDatabase(t, logger)

	queries := repository.New(pool)

	fundSvc := service.NewFundService(queries)
	NewFundHandler(api, fundSvc)

	investorSvc := service.NewInvestorService(queries)
	NewInvestorHandler(api, investorSvc)

	investmentSvc := service.NewInvestmentService(queries, fundSvc, investorSvc)
	NewInvestmentHandler(api, investmentSvc)

	return &testHarness{api: api, pool: pool}
}
