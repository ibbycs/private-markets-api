//go:build integration

package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFundLifecycle(t *testing.T) {
	h := newTestHarness(t)

	// Create
	createResp := h.api.Post("/funds", map[string]interface{}{
		"name":            "Test Growth Fund",
		"vintage_year":    2025,
		"target_size_usd": 250000000.00,
		"status":          "fundraising",
	})
	require.Equal(t, 201, createResp.Code)

	var created dto.Fund
	err := json.NewDecoder(createResp.Body).Decode(&created)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Test Growth Fund", created.Name)
	assert.Equal(t, 2025, created.VintageYear)
	assert.Equal(t, 250000000.00, created.TargetSizeUSD)
	assert.Equal(t, "fundraising", created.Status)
	assert.False(t, created.CreatedAt.IsZero())

	// Get by ID
	getResp := h.api.Get("/funds/" + created.ID)
	require.Equal(t, 200, getResp.Code)

	var fetched dto.Fund
	err = json.NewDecoder(getResp.Body).Decode(&fetched)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Name, fetched.Name)

	// List
	listResp := h.api.Get("/funds?limit=10&offset=0")
	require.Equal(t, 200, listResp.Code)

	var funds []dto.Fund
	err = json.NewDecoder(listResp.Body).Decode(&funds)
	require.NoError(t, err)
	require.Len(t, funds, 1)
	assert.Equal(t, created.ID, funds[0].ID)

	// Update
	updateResp := h.api.Put("/funds", map[string]interface{}{
		"id":              created.ID,
		"name":            "Test Growth Fund II",
		"vintage_year":    2025,
		"target_size_usd": 300000000.00,
		"status":          "investing",
	})
	require.Equal(t, 200, updateResp.Code)

	var updated dto.Fund
	err = json.NewDecoder(updateResp.Body).Decode(&updated)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Test Growth Fund II", updated.Name)
	assert.Equal(t, 300000000.00, updated.TargetSizeUSD)
	assert.Equal(t, "investing", updated.Status)
	assert.Equal(t, created.CreatedAt.Unix(), updated.CreatedAt.Unix())
}

func TestCreateFund_InvalidInput(t *testing.T) {
	h := newTestHarness(t)

	tests := []struct {
		name string
		body map[string]interface{}
	}{
		{"empty name", map[string]interface{}{"name": "", "vintage_year": 2025, "target_size_usd": 100.0, "status": "fundraising"}},
		{"invalid status", map[string]interface{}{"name": "Fund", "vintage_year": 2025, "target_size_usd": 100.0, "status": "invalid"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := h.api.Post("/funds", tt.body)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		})
	}
}

func TestGetFund_NotFound(t *testing.T) {
	h := newTestHarness(t)

	resp := h.api.Get("/funds/00000000-0000-0000-0000-000000000000")
	assert.Equal(t, 404, resp.Code)
}

func TestUpdateFund_NotFound(t *testing.T) {
	h := newTestHarness(t)

	resp := h.api.Put("/funds", map[string]interface{}{
		"id":              "00000000-0000-0000-0000-000000000000",
		"name":            "Ghost Fund",
		"vintage_year":    2025,
		"target_size_usd": 100.0,
		"status":          "fundraising",
	})
	assert.Equal(t, 404, resp.Code)
}
