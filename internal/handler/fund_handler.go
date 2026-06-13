package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/ibbycs/private-markets-api/internal/mapper"
	"github.com/ibbycs/private-markets-api/internal/service"
)

type fundHandler struct {
	svc *service.FundService
}

func NewFundHandler(api huma.API, svc *service.FundService) {
	h := &fundHandler{svc: svc}

	huma.Register(api, huma.Operation{
		OperationID: "list-funds",
		Method:      http.MethodGet,
		Path:        "/funds",
		Summary:     "List all funds",
		Tags:        []string{"Funds"},
	}, h.List)

	huma.Register(api, huma.Operation{
		OperationID: "create-fund",
		Method:      http.MethodPost,
		Path:        "/funds",
		Summary:     "Create a new fund",
		Tags:        []string{"Funds"},
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID: "update-fund",
		Method:      http.MethodPut,
		Path:        "/funds",
		Summary:     "Update an existing fund",
		Tags:        []string{"Funds"},
	}, h.Update)

	huma.Register(api, huma.Operation{
		OperationID: "get-fund",
		Method:      http.MethodGet,
		Path:        "/funds/{id}",
		Summary:     "Get a specific fund",
		Tags:        []string{"Funds"},
	}, h.GetByID)
}

type listFundsOutput struct {
	Body []dto.Fund
}

func (h *fundHandler) List(ctx context.Context, input *dto.ListFundsRequest) (*listFundsOutput, error) {
	funds, err := h.svc.List(ctx, int32(input.Limit), int32(input.Offset))
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list funds", err)
	}
	return &listFundsOutput{Body: mapper.MapFundsToDTOs(funds)}, nil
}

type createFundInput struct {
	Body *dto.CreateFundRequest
}

type createFundOutput struct {
	Status int
	Body   *dto.Fund
}

func (h *fundHandler) Create(ctx context.Context, input *createFundInput) (*createFundOutput, error) {
	fund, err := h.svc.Create(ctx, input.Body.Name, input.Body.VintageYear, input.Body.TargetSizeUSD, input.Body.Status)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid input", err)
	}
	dto := mapper.MapFundToDTO(fund)
	return &createFundOutput{Status: 201, Body: &dto}, nil
}

type updateFundInput struct {
	Body *dto.UpdateFundRequest
}

type updateFundOutput struct {
	Body *dto.Fund
}

func (h *fundHandler) Update(ctx context.Context, input *updateFundInput) (*updateFundOutput, error) {
	fund, err := h.svc.Update(ctx, input.Body.ID, input.Body.Name, input.Body.VintageYear, input.Body.TargetSizeUSD, input.Body.Status)
	if err != nil {
		return nil, huma.Error404NotFound("Not found", err)
	}
	dto := mapper.MapFundToDTO(fund)
	return &updateFundOutput{Body: &dto}, nil
}

type getFundInput struct {
	ID string `path:"id" required:"true"`
}

type getFundOutput struct {
	Body *dto.Fund
}

func (h *fundHandler) GetByID(ctx context.Context, input *getFundInput) (*getFundOutput, error) {
	fund, err := h.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("Not found", err)
	}
	dto := mapper.MapFundToDTO(fund)
	return &getFundOutput{Body: &dto}, nil
}
