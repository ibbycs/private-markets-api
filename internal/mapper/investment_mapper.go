package mapper

import (
	"github.com/ibbycs/private-markets-api/internal/dto"
	"github.com/ibbycs/private-markets-api/internal/repository"
)

func MapInvestmentToDTO(m repository.Investment) dto.Investment {
	usd, _ := m.AmountUsd.Float64Value()
	return dto.Investment{
		ID:             m.ID.String(),
		InvestorID:     m.InvestorID.String(),
		FundID:         m.FundID.String(),
		AmountUSD:      usd.Float64,
		InvestmentDate: m.InvestmentDate.Time.Format("2006-01-02"),
	}
}

func MapInvestmentsToDTOs(investments []repository.Investment) []dto.Investment {
	dtos := make([]dto.Investment, 0, len(investments))
	for _, m := range investments {
		dtos = append(dtos, MapInvestmentToDTO(m))
	}
	return dtos
}
