-- +goose Up
CREATE TABLE cardio_and_sets(
    id UUID PRIMARY KEY,
    workout_exercises_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    set_number INT NOT NULL DEFAULT 1,
    distance DECIMAL(7,2),
    is_kilometers BOOLEAN NOT NULL DEFAULT true,
    duration_seconds INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE cardio_and_sets;
