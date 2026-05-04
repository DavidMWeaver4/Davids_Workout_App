-- +goose Up
CREATE TABLE exercises (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	exercise_name VARCHAR(100) NOT NULL,
	target_muscles TEXT[],
	equipment VARCHAR(100),
	difficulty_level VARCHAR(20) CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE exercises;
