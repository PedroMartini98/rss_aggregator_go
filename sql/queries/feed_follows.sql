-- name: CreateFollow :one
INSERT INTO feed_follows(created_at,user_id,feed_id,feed_url)
SELECT $1,$2,$3,f.url
FROM feeds f
WHERE f.id = $3
RETURNING *;

-- name: DeleteFollow :one
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2
RETURNING *;

