package dto

import "time"

type Fund struct {
	ID            string    `json:"id" format:"uuid"`
	Name          string    `json:"name"`
	VintageYear   int       `json:"vintage_year"`
	TargetSizeUSD float64   `json:"target_size_usd"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at" format:"date-time"`
}

type ListFundsRequest struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"20"`
	Offset int `query:"offset" minimum:"0" default:"0"`
}

type CreateFundRequest struct {
	Name          string  `json:"name" required:"true" minLength:"1" example:"Titanbay Growth Fund II"`
	VintageYear   int     `json:"vintage_year" required:"true" example:"2025"`
	TargetSizeUSD float64 `json:"target_size_usd" required:"true" minimum:"0" example:"500000000.00"`
	Status        string  `json:"status" required:"true" enum:"fundraising,investing,closed" example:"fundraising"`
}

type UpdateFundRequest struct {
	ID            string  `json:"id" required:"true" format:"uuid"`
	Name          string  `json:"name" required:"true" minLength:"1" example:"Titanbay Growth Fund I"`
	VintageYear   int     `json:"vintage_year" required:"true" example:"2024"`
	TargetSizeUSD float64 `json:"target_size_usd" required:"true" minimum:"0" example:"300000000.00"`
	Status        string  `json:"status" required:"true" enum:"fundraising,investing,closed" example:"investing"`
}
