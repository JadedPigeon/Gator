-- +goose up
CREATE TABLE posts (
    id uuid primary key default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    title text not null,
    url text unique not null,
    description text,
    published_at timestamp,
    feed_id uuid not null references feeds(id) on delete cascade
);

-- +goose down
DROP TABLE posts;