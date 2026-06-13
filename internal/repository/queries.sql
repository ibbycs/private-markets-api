-- name: ListFunds :many
SELECT * FROM fund ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetFundByID :one
SELECT * FROM fund WHERE id = $1;

-- name: CreateFund :one
INSERT INTO fund (id, name, vintage_year, target_size_usd, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateFund :one
UPDATE fund
SET name = $2, vintage_year = $3, target_size_usd = $4, status = $5
WHERE id = $1
RETURNING *;

-- name: ListInvestors :many
SELECT * FROM investor ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetInvestorByID :one
SELECT * FROM investor WHERE id = $1;

-- name: GetInvestorByEmail :one
SELECT * FROM investor WHERE email = $1;

-- name: CreateInvestor :one
INSERT INTO investor (id, name, investor_type, email, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListInvestmentsByFundID :many
SELECT * FROM investment WHERE fund_id = $1;

-- name: CreateInvestment :one
INSERT INTO investment (id, investor_id, fund_id, amount_usd, investment_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
