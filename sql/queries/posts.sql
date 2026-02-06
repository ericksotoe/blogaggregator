-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name
FROM posts
INNER JOIN feed_follows 
    ON posts.feed_id = feed_follows.feed_id
INNER JOIN feeds
    ON posts.feed_id = feeds.id

WHERE $1 = feed_follows.user_id
ORDER BY published_at DESC
LIMIT $2;
