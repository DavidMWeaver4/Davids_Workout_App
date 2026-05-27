-- +goose Up
ALTER TABLE weights_and_sets
ALTER COLUMN weight TYPE DOUBLE PRECISION
USING weight::DOUBLE PRECISION;

-- +goose Down
ALTER TABLE weights_and_sets
ALTER COLUMN weight TYPE DECIMAL(5,2)
USING ROUND(weight::numeric, 2);
