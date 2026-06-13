package dto

import "time"

type Investor struct {
	ID           string    `json:"id" format:"uuid"`
	Name         string    `json:"name"`
	InvestorType string    `json:"investor_type"`
	Email        string    `json:"email" format:"email"`
	CreatedAt    time.Time `json:"created_at" format:"date-time"`
}

type ListInvestorsRequest struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"20"`
	Offset int `query:"offset" minimum:"0" default:"0"`
}

type CreateInvestorRequest struct {
	Name         string `json:"name" required:"true" minLength:"1" example:"CalPERS"`
	InvestorType string `json:"investor_type" required:"true" enum:"individual,institution,family office" example:"institution"`
	Email        string `json:"email" required:"true" format:"email" example:"privateequity@calpers.ca.gov"`
}
