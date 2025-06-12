-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllFeeds :many
SELECT 
    feeds.id,
    feeds.created_at,
    feeds.updated_at,
    feeds.name,
    feeds.url,
    users.name AS user_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id
ORDER BY feeds.created_at DESC;

-- name: GetFeedByUrl :one
SELECT 
    feeds.id,
    feeds.created_at,
    feeds.updated_at,
    feeds.name,
    feeds.url
FROM feeds
WHERE feeds.url = $1;