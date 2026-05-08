-- +goose Up
ALTER TABLE exercises ADD COLUMN description TEXT;

-- +goose Down
ALTER TABLE exercises DROP COLUMN description;
