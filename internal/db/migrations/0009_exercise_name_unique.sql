-- +goose Up
ALTER TABLE exercises ADD CONSTRAINT exercises_name_unique UNIQUE (exercise_name);

-- +goose Down
ALTER TABLE exercises DROP CONSTRAINT exercises_name_unique;
