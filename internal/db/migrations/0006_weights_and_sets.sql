-- +goose Up
CREATE TABLE weights_and_sets(
    id UUID PRIMARY KEY,
    workout_exercises_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0,
    is_kilograms BOOLEAN NOT NULL DEFAULT true,
    set_number INT NOT NULL DEFAULT 1,
    reps_target INT NOT NULL DEFAULT 1,
    reps_actual INT NOT NULL DEFAULT 0,
    duration_seconds INT,
    rest_time_seconds INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE weights_and_sets;
