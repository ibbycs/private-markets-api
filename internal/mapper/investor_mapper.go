package mapper

import (
	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/ibbycs/private-markets-api/internal/repository"
)

func MapInvestorToDTO(m repository.Investor) dto.Investor {
	return dto.Investor{
		ID:           m.ID.String(),
		Name:         m.Name,
		InvestorType: string(m.InvestorType),
		Email:        m.Email,
		CreatedAt:    m.CreatedAt.Time,
	}
}

func MapInvestorsToDTOs(investors []repository.Investor) []dto.Investor {
	dtos := make([]dto.Investor, 0, len(investors))
	for _, m := range investors {
		dtos = append(dtos, MapInvestorToDTO(m))
	}
	return dtos
}
