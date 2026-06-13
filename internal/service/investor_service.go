package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ibbycs/private-markets-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrEmailTaken = errors.New("email already in use")

func isValidInvestorType(s string) bool {
	switch repository.InvestorType(s) {
	case repository.InvestorTypeIndividual, repository.InvestorTypeInstitution, repository.InvestorTypeFamilyoffice:
		return true
	}
	return false
}

type InvestorService struct {
	queries *repository.Queries
}

func NewInvestorService(queries *repository.Queries) *InvestorService {
	return &InvestorService{queries: queries}
}

func (s *InvestorService) List(ctx context.Context, limit, offset int32) ([]repository.Investor, error) {
	return s.queries.ListInvestors(ctx, repository.ListInvestorsParams{Limit: limit, Offset: offset})
}

func (s *InvestorService) GetByID(ctx context.Context, id string) (repository.Investor, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return repository.Investor{}, fmt.Errorf("invalid investor id: %s", id)
	}

	investor, err := s.queries.GetInvestorByID(ctx, uid)
	if err != nil {
		return repository.Investor{}, fmt.Errorf("investor not found: %s", id)
	}
	return investor, nil
}

func (s *InvestorService) Create(ctx context.Context, name, investorType, email string) (repository.Investor, error) {
	if name == "" {
		return repository.Investor{}, fmt.Errorf("name is required")
	}
	if !isValidInvestorType(investorType) {
		return repository.Investor{}, fmt.Errorf("invalid investor_type: %s (must be Individual, Institution, or Family Office)", investorType)
	}
	if !strings.Contains(email, "@") {
		return repository.Investor{}, fmt.Errorf("invalid email: %s", email)
	}

	existing, err := s.queries.GetInvestorByEmail(ctx, email)
	if err == nil {
		return existing, ErrEmailTaken
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return repository.Investor{}, fmt.Errorf("failed to check email: %w", err)
	}

	uid, _ := parseUUID(uuid.Must(uuid.NewV7()).String())

	params := repository.CreateInvestorParams{
		ID:           uid,
		Name:         name,
		InvestorType: repository.InvestorType(investorType),
		Email:        email,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	}

	return s.queries.CreateInvestor(ctx, params)
}
