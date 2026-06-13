package mapper

import (
	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/ibbycs/private-markets-api/internal/repository"
)

func MapFundToDTO(m repository.Fund) dto.Fund {
	usd, _ := m.TargetSizeUsd.Float64Value()
	return dto.Fund{
		ID:            m.ID.String(),
		Name:          m.Name,
		VintageYear:   int(m.VintageYear),
		TargetSizeUSD: usd.Float64,
		Status:        string(m.Status),
		CreatedAt:     m.CreatedAt.Time,
	}
}

func MapFundsToDTOs(funds []repository.Fund) []dto.Fund {
	dtos := make([]dto.Fund, 0, len(funds))
	for _, m := range funds {
		dtos = append(dtos, MapFundToDTO(m))
	}
	return dtos
}
