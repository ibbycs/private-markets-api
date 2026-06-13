package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/ibbycs/private-markets-api/internal/mapper"
	"github.com/ibbycs/private-markets-api/internal/service"
)

type investmentHandler struct {
	svc *service.InvestmentService
}

func NewInvestmentHandler(api huma.API, svc *service.InvestmentService) {
	h := &investmentHandler{svc: svc}

	huma.Register(api, huma.Operation{
		OperationID: "list-investments",
		Method:      http.MethodGet,
		Path:        "/funds/{fund_id}/investments",
		Summary:     "List all investments for a specific fund",
		Tags:        []string{"Investments"},
	}, h.ListByFundID)

	huma.Register(api, huma.Operation{
		OperationID: "create-investment",
		Method:      http.MethodPost,
		Path:        "/funds/{fund_id}/investments",
		Summary:     "Create a new investment to a fund",
		Tags:        []string{"Investments"},
	}, h.Create)
}

type listInvestmentsInput struct {
	FundID string `path:"fund_id" required:"true"`
}

type listInvestmentsOutput struct {
	Body []dto.Investment
}

func (h *investmentHandler) ListByFundID(ctx context.Context, input *listInvestmentsInput) (*listInvestmentsOutput, error) {
	investments, err := h.svc.ListByFundID(ctx, input.FundID)
	if err != nil {
		return nil, huma.Error404NotFound("Not found", err)
	}
	return &listInvestmentsOutput{Body: mapper.MapInvestmentsToDTOs(investments)}, nil
}

type createInvestmentInput struct {
	FundID string `path:"fund_id" required:"true"`
	Body   *dto.CreateInvestmentRequest
}

type createInvestmentOutput struct {
	Status int
	Body   *dto.Investment
}

func (h *investmentHandler) Create(ctx context.Context, input *createInvestmentInput) (*createInvestmentOutput, error) {
	investment, err := h.svc.Create(ctx, input.FundID, input.Body.InvestorID, input.Body.AmountUSD, input.Body.InvestmentDate)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid input", err)
	}
	dto := mapper.MapInvestmentToDTO(investment)
	return &createInvestmentOutput{Status: 201, Body: &dto}, nil
}
