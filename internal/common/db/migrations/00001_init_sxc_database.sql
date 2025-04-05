-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS promocode (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    advantage_percent INTEGER NOT NULL,
    restrictions JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS weather (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    town VARCHAR(256) NOT NULL,
    temp DECIMAL NOT NULL,
    type VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);


-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS promocode;
DROP TABLE IF EXISTS weather;