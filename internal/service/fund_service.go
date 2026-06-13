package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ibbycs/private-markets-api/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

func isValidFundStatus(s string) bool {
	switch repository.FundStatus(s) {
	case repository.FundStatusFundraising, repository.FundStatusInvesting, repository.FundStatusClosed:
		return true
	}
	return false
}

type FundService struct {
	queries *repository.Queries
}

func NewFundService(queries *repository.Queries) *FundService {
	return &FundService{queries: queries}
}

func (s *FundService) List(ctx context.Context, limit, offset int32) ([]repository.Fund, error) {
	return s.queries.ListFunds(ctx, repository.ListFundsParams{Limit: limit, Offset: offset})
}

func (s *FundService) GetByID(ctx context.Context, id string) (repository.Fund, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return repository.Fund{}, fmt.Errorf("invalid fund id: %s", id)
	}

	fund, err := s.queries.GetFundByID(ctx, uid)
	if err != nil {
		return repository.Fund{}, fmt.Errorf("fund not found: %s", id)
	}
	return fund, nil
}

func (s *FundService) Create(ctx context.Context, name string, vintageYear int, targetSizeUSD float64, status string) (repository.Fund, error) {
	if name == "" {
		return repository.Fund{}, fmt.Errorf("name is required")
	}
	if !isValidFundStatus(status) {
		return repository.Fund{}, fmt.Errorf("invalid status: %s (must be Fundraising, Investing, or Closed)", status)
	}

	uid, _ := parseUUID(uuid.Must(uuid.NewV7()).String())
	numeric, err := parseNumeric(targetSizeUSD)
	if err != nil {
		return repository.Fund{}, err
	}

	params := repository.CreateFundParams{
		ID:            uid,
		Name:          name,
		VintageYear:   int32(vintageYear),
		TargetSizeUsd: numeric,
		Status:        repository.FundStatus(status),
		CreatedAt:     pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	}

	return s.queries.CreateFund(ctx, params)
}

func (s *FundService) Update(ctx context.Context, id, name string, vintageYear int, targetSizeUSD float64, status string) (repository.Fund, error) {
	if !isValidFundStatus(status) {
		return repository.Fund{}, fmt.Errorf("invalid status: %s (must be Fundraising, Investing, or Closed)", status)
	}

	uid, err := parseUUID(id)
	if err != nil {
		return repository.Fund{}, fmt.Errorf("invalid fund id: %s", id)
	}

	numeric, err := parseNumeric(targetSizeUSD)
	if err != nil {
		return repository.Fund{}, err
	}

	params := repository.UpdateFundParams{
		ID:            uid,
		Name:          name,
		VintageYear:   int32(vintageYear),
		TargetSizeUsd: numeric,
		Status:        repository.FundStatus(status),
	}

	fund, err := s.queries.UpdateFund(ctx, params)
	if err != nil {
		return repository.Fund{}, fmt.Errorf("fund not found: %s", id)
	}
	return fund, nil
}

func parseNumeric(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(strconv.FormatFloat(f, 'f', 2, 64))
	return n, err
}

func parseUUID(s string) (pgtype.UUID, error) {
	var uid pgtype.UUID
	if err := uid.Scan(s); err != nil {
		return pgtype.UUID{}, err
	}
	return uid, nil
}
