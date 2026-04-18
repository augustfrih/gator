-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT DISTINCT posts.*, feeds.name AS feed_name FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN posts ON feed_follows.feed_id = posts.feed_id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE users.name = $1
ORDER BY posts.published_at DESC
LIMIT $2;
