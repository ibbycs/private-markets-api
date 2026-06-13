package models

type Investment struct {
	ID             string  `json:"id"`
	InvestorID     string  `json:"investor_id"`
	FundID         string  `json:"fund_id"`
	AmountUSD      float64 `json:"amount_usd"`
	InvestmentDate string  `json:"investment_date"`
}
