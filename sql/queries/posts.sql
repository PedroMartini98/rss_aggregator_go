-- name: CreatePost :one

INSERT INTO posts (
id,
created_at,
updated_at,
published_at,
url,
title,
description,
feed_id
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- name: GetPostsForUsers :many
SELECT posts.* FROM posts
JOIN feed_follows on posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
