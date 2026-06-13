//go:build integration

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFundForTest(t *testing.T, h *testHarness) string {
	t.Helper()
	resp := h.api.Post("/funds", map[string]interface{}{
		"name":            "Seed Fund",
		"vintage_year":    2024,
		"target_size_usd": 100000000.00,
		"status":          "fundraising",
	})
	require.Equal(t, http.StatusCreated, resp.Code)
	var fund dto.Fund
	json.NewDecoder(resp.Body).Decode(&fund)
	return fund.ID
}

func createInvestorForTest(t *testing.T, h *testHarness) string {
	t.Helper()
	resp := h.api.Post("/investors", map[string]interface{}{
		"name":          "Test Investor",
		"investor_type": "individual",
		"email":         "test@example.com",
	})
	require.Equal(t, http.StatusCreated, resp.Code)
	var inv dto.Investor
	json.NewDecoder(resp.Body).Decode(&inv)
	return inv.ID
}

func TestInvestmentLifecycle(t *testing.T) {
	h := newTestHarness(t)

	fundID := createFundForTest(t, h)
	investorID := createInvestorForTest(t, h)

	// Create
	createResp := h.api.Post(fmt.Sprintf("/funds/%s/investments", fundID), map[string]interface{}{
		"investor_id":     investorID,
		"amount_usd":      50000000.00,
		"investment_date": "2024-03-15",
	})
	require.Equal(t, 201, createResp.Code)

	var created dto.Investment
	err := json.NewDecoder(createResp.Body).Decode(&created)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, investorID, created.InvestorID)
	assert.Equal(t, fundID, created.FundID)
	assert.Equal(t, 50000000.00, created.AmountUSD)
	assert.Equal(t, "2024-03-15", created.InvestmentDate)

	// List by fund ID
	listResp := h.api.Get(fmt.Sprintf("/funds/%s/investments", fundID))
	require.Equal(t, http.StatusOK, listResp.Code)

	var investments []dto.Investment
	err = json.NewDecoder(listResp.Body).Decode(&investments)
	require.NoError(t, err)
	require.Len(t, investments, 1)
	assert.Equal(t, created.ID, investments[0].ID)
}

func TestCreateInvestment_InvalidInput(t *testing.T) {
	h := newTestHarness(t)

	fundID := createFundForTest(t, h)
	investorID := createInvestorForTest(t, h)

	t.Run("nonexistent fund", func(t *testing.T) {
		resp := h.api.Post("/funds/00000000-0000-0000-0000-000000000000/investments", map[string]interface{}{
			"investor_id":     investorID,
			"amount_usd":      100.0,
			"investment_date": "2024-01-01",
		})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("nonexistent investor", func(t *testing.T) {
		resp := h.api.Post(fmt.Sprintf("/funds/%s/investments", fundID), map[string]interface{}{
			"investor_id":     "00000000-0000-0000-0000-000000000000",
			"amount_usd":      100.0,
			"investment_date": "2024-01-01",
		})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("negative amount", func(t *testing.T) {
		resp := h.api.Post(fmt.Sprintf("/funds/%s/investments", fundID), map[string]interface{}{
			"investor_id":     investorID,
			"amount_usd":      -100.0,
			"investment_date": "2024-01-01",
		})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
}

func TestListInvestments_NotFound(t *testing.T) {
	h := newTestHarness(t)

	resp := h.api.Get("/funds/00000000-0000-0000-0000-000000000000/investments")
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
