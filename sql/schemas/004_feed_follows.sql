-- +goose Up
CREATE TABLE feed_follows(
    created_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE ,
    feed_url TEXT NOT NULL REFERENCES feeds(url) ON DELETE CASCADE,
    PRIMARY KEY (user_id,feed_id)
);


-- +goose Down
DROP TABLE feed_follows;
