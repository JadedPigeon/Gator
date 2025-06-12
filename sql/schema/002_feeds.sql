-- +goose up
CREATE TABLE feeds (
    id uuid primary key default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    name text unique not null,
    url text not null unique,
    user_id uuid not null references users(id) on delete cascade
);

-- +goose down
DROP TABLE feeds;