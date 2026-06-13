package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ibbycs/private-markets-api/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type InvestmentService struct {
	queries     *repository.Queries
	fundSvc     *FundService
	investorSvc *InvestorService
}

func NewInvestmentService(queries *repository.Queries, fundSvc *FundService, investorSvc *InvestorService) *InvestmentService {
	return &InvestmentService{
		queries:     queries,
		fundSvc:     fundSvc,
		investorSvc: investorSvc,
	}
}

func (s *InvestmentService) ListByFundID(ctx context.Context, fundID string) ([]repository.Investment, error) {
	if _, err := s.fundSvc.GetByID(ctx, fundID); err != nil {
		return nil, err
	}

	uid, err := parseUUID(fundID)
	if err != nil {
		return nil, fmt.Errorf("invalid fund id: %s", fundID)
	}

	return s.queries.ListInvestmentsByFundID(ctx, uid)
}

func (s *InvestmentService) Create(ctx context.Context, fundID, investorID string, amountUSD float64, investmentDate string) (repository.Investment, error) {
	if _, err := s.fundSvc.GetByID(ctx, fundID); err != nil {
		return repository.Investment{}, fmt.Errorf("fund not found: %s", fundID)
	}
	if _, err := s.investorSvc.GetByID(ctx, investorID); err != nil {
		return repository.Investment{}, fmt.Errorf("investor not found: %s", investorID)
	}
	if amountUSD <= 0 {
		return repository.Investment{}, fmt.Errorf("amount_usd must be positive")
	}
	if investmentDate == "" {
		return repository.Investment{}, fmt.Errorf("investment_date is required")
	}

	fundUID, _ := parseUUID(fundID)
	investorUID, _ := parseUUID(investorID)
	newUID, _ := parseUUID(uuid.Must(uuid.NewV7()).String())

	numeric, err := parseNumeric(amountUSD)
	if err != nil {
		return repository.Investment{}, err
	}

	var date pgtype.Date
	if err := date.Scan(investmentDate); err != nil {
		return repository.Investment{}, fmt.Errorf("invalid investment_date: %s", investmentDate)
	}

	params := repository.CreateInvestmentParams{
		ID:             newUID,
		InvestorID:     investorUID,
		FundID:         fundUID,
		AmountUsd:      numeric,
		InvestmentDate: date,
	}

	return s.queries.CreateInvestment(ctx, params)
}
