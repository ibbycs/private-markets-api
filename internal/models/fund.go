package models

import "time"

type Fund struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	VintageYear   int       `json:"vintage_year"`
	TargetSizeUSD float64   `json:"target_size_usd"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
