-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    description,
    published_at,
    url,
    feed_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPost :one
SELECT *
FROM posts
WHERE id = $1;

-- name: GetPosts :many
SELECT *
FROM posts;

-- name: GetPostsByUser :many
SELECT posts.* 
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;