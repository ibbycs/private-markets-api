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

func TestInvestorLifecycle(t *testing.T) {
	h := newTestHarness(t)

	// Create
	createResp := h.api.Post("/investors", map[string]interface{}{
		"name":          "CalPERS",
		"investor_type": "institution",
		"email":         "privateequity@calpers.ca.gov",
	})
	require.Equal(t, 201, createResp.Code)

	var created dto.Investor
	err := json.NewDecoder(createResp.Body).Decode(&created)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "CalPERS", created.Name)
	assert.Equal(t, "institution", created.InvestorType)
	assert.Equal(t, "privateequity@calpers.ca.gov", created.Email)
	assert.False(t, created.CreatedAt.IsZero())

	// List
	listResp := h.api.Get("/investors?limit=10&offset=0")
	require.Equal(t, http.StatusOK, listResp.Code)

	var investors []dto.Investor
	err = json.NewDecoder(listResp.Body).Decode(&investors)
	require.NoError(t, err)
	require.Len(t, investors, 1)
	assert.Equal(t, created.ID, investors[0].ID)
	assert.Equal(t, created.Email, investors[0].Email)
}

func TestCreateInvestor_DuplicateEmail(t *testing.T) {
	h := newTestHarness(t)

	body := map[string]interface{}{
		"name":          "CalPERS",
		"investor_type": "institution",
		"email":         "dup@test.com",
	}

	// First creation succeeds
	createResp := h.api.Post("/investors", body)
	require.Equal(t, 201, createResp.Code)

	// Second creation with same email fails
	dupResp := h.api.Post("/investors", body)
	require.Equal(t, http.StatusConflict, dupResp.Code)
}

func TestCreateInvestor_InvalidInput(t *testing.T) {
	h := newTestHarness(t)

	tests := []struct {
		name string
		body map[string]interface{}
	}{
		{"empty name", map[string]interface{}{"name": "", "investor_type": "institution", "email": "a@b.com"}},
		{"invalid type", map[string]interface{}{"name": "Test", "investor_type": "alien", "email": "a@b.com"}},
		{"missing email", map[string]interface{}{"name": "Test", "investor_type": "institution", "email": ""}},
		{"invalid email format", map[string]interface{}{"name": "Test", "investor_type": "institution", "email": "notanemail"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := h.api.Post("/investors", tt.body)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		})
	}
}
