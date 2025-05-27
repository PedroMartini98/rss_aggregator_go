-- +goose Up

ALTER TABLE feeds ADD COLUMN last_fetched TIMESTAMP;

-- +goose Down

ALTER TABLE feeds DROPS COLUMN last_fetched;

