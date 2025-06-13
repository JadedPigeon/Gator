-- name: CreatePost :one
INSERT INTO posts (title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPostsForUser :many
SELECT 
    p.id, 
    p.title, 
    p.url, 
    p.description, 
    p.published_at, 
    p.feed_id
FROM posts p
JOIN feeds f ON p.feed_id = f.id
WHERE f.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;
