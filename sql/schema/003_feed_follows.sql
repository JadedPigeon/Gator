-- +goose up
CREATE TABLE feed_follows (
    id uuid primary key default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    feed_id uuid not null references feeds(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    unique (feed_id, user_id)
);

-- +goose down
DROP TABLE feed_follows;