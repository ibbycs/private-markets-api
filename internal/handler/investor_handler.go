package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/ibbycs/private-markets-api/internal/mapper"
	"github.com/ibbycs/private-markets-api/internal/service"
)

type investorHandler struct {
	svc *service.InvestorService
}

func NewInvestorHandler(api huma.API, svc *service.InvestorService) {
	h := &investorHandler{svc: svc}

	huma.Register(api, huma.Operation{
		OperationID: "list-investors",
		Method:      http.MethodGet,
		Path:        "/investors",
		Summary:     "List all investors",
		Tags:        []string{"Investors"},
	}, h.List)

	huma.Register(api, huma.Operation{
		OperationID: "create-investor",
		Method:      http.MethodPost,
		Path:        "/investors",
		Summary:     "Create a new investor",
		Tags:        []string{"Investors"},
	}, h.Create)
}

type listInvestorsOutput struct {
	Body []dto.Investor
}

func (h *investorHandler) List(ctx context.Context, input *dto.ListInvestorsRequest) (*listInvestorsOutput, error) {
	investors, err := h.svc.List(ctx, int32(input.Limit), int32(input.Offset))
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list investors", err)
	}
	return &listInvestorsOutput{Body: mapper.MapInvestorsToDTOs(investors)}, nil
}

type createInvestorInput struct {
	Body *dto.CreateInvestorRequest
}

type createInvestorOutput struct {
	Status int
	Body   *dto.Investor
}

func (h *investorHandler) Create(ctx context.Context, input *createInvestorInput) (*createInvestorOutput, error) {
	investor, err := h.svc.Create(ctx, input.Body.Name, input.Body.InvestorType, input.Body.Email)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			return nil, huma.Error409Conflict("Email already in use", err)
		}
		return nil, huma.Error400BadRequest("Invalid input", err)
	}
	dto := mapper.MapInvestorToDTO(investor)
	return &createInvestorOutput{Status: 201, Body: &dto}, nil
}
