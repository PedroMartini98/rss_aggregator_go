// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFollow = `-- name: CreateFollow :one
INSERT INTO feed_follows(created_at,user_id,feed_id,feed_url)
SELECT $1,$2,$3,f.url
FROM feeds f
WHERE f.id = $3
RETURNING created_at, feed_id, user_id, feed_url
`

type CreateFollowParams struct {
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func (q *Queries) CreateFollow(ctx context.Context, arg CreateFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFollow, arg.CreatedAt, arg.UserID, arg.FeedID)
	var i FeedFollow
	err := row.Scan(
		&i.CreatedAt,
		&i.FeedID,
		&i.UserID,
		&i.FeedUrl,
	)
	return i, err
}

const deleteFollow = `-- name: DeleteFollow :one
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2
RETURNING created_at, feed_id, user_id, feed_url
`

type DeleteFollowParams struct {
	FeedID uuid.UUID `json:"feed_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteFollow(ctx context.Context, arg DeleteFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, deleteFollow, arg.FeedID, arg.UserID)
	var i FeedFollow
	err := row.Scan(
		&i.CreatedAt,
		&i.FeedID,
		&i.UserID,
		&i.FeedUrl,
	)
	return i, err
}
