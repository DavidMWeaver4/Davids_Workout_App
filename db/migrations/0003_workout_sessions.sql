-- +goose Up
CREATE TABLE workout_sessions
 (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	workout_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	description TEXT,
	notes TEXT,
	created_at TIMESTAMPTZ NOT NULL  DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE workout_sessions;
