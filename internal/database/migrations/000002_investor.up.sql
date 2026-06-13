CREATE TYPE investor_type AS ENUM ('individual', 'institution', 'family office');

CREATE TABLE investor (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    investor_type investor_type NOT NULL,
    email VARCHAR(320) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
