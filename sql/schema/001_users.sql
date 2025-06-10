-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id uuid primary key default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    name text unique not null
);

-- +goose Down
DROP TABLE users;
