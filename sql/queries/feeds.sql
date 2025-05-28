-- name: CreateFeed :one
INSERT INTO feeds(id,created_at,updated_at,name,url,user_id)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds where id = $1;

-- name: GetFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched ASC NULLS FIRST
LIMIT $1;

-- name: MarkFetched :one
UPDATE feeds
SET last_fetched = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

