package models

import "time"

type Investor struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	InvestorType string    `json:"investor_type"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}
