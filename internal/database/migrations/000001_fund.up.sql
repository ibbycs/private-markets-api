CREATE TYPE fund_status AS ENUM ('fundraising', 'investing', 'closed');

CREATE TABLE fund (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    vintage_year INTEGER NOT NULL,
    target_size_usd DECIMAL(15,2) NOT NULL,
    status fund_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
