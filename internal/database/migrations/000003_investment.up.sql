CREATE TABLE investment (
    id UUID PRIMARY KEY,
    investor_id UUID NOT NULL REFERENCES investor(id),
    fund_id UUID NOT NULL REFERENCES fund(id),
    amount_usd DECIMAL(15,2) NOT NULL CHECK (amount_usd > 0),
    investment_date DATE NOT NULL
);
