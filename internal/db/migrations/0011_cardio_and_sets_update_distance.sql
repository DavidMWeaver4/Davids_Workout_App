-- +goose Up
ALTER TABLE cardio_and_sets
ALTER COLUMN distance TYPE DOUBLE PRECISION
USING distance::DOUBLE PRECISION;

-- +goose Down
ALTER TABLE cardio_and_sets
ALTER COLUMN distance TYPE DECIMAL(7,2)
USING ROUND(distance::numeric, 2);
