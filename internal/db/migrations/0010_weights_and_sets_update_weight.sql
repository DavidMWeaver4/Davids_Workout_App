-- +goose Up
ALTER TABLE weights_and_sets
ALTER COLUMN weight TYPE float8;

-- +goose Down
ALTER TABLE weights_and_sets
ALTER COLUMN weight TYPE string;
