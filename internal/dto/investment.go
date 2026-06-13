package dto

type Investment struct {
	ID             string  `json:"id" format:"uuid"`
	InvestorID     string  `json:"investor_id" format:"uuid"`
	FundID         string  `json:"fund_id" format:"uuid"`
	AmountUSD      float64 `json:"amount_usd"`
	InvestmentDate string  `json:"investment_date" format:"date"`
}

type CreateInvestmentRequest struct {
	InvestorID     string  `json:"investor_id" required:"true" format:"uuid"`
	AmountUSD      float64 `json:"amount_usd" required:"true" minimum:"1" example:"75000000.00"`
	InvestmentDate string  `json:"investment_date" required:"true" format:"date" example:"2024-09-22"`
}
