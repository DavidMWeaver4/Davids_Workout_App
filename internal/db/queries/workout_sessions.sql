-- name: CreateWorkoutSessions :one
INSERT INTO workout_sessions(id, user_id, workout_date, description, notes, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetWorkoutSessionByID :one
SELECT * FROM workout_sessions
WHERE id = $1;

-- name: GetAllWorkoutSessionsSorted :many
SELECT * FROM workout_sessions
WHERE user_id = $1
ORDER BY workout_date DESC;

-- name: GetLastWorkoutSession :one
SELECT * FROM workout_sessions
WHERE user_id = $1
ORDER BY workout_date DESC
LIMIT 1;

-- name: GetWorkoutSessionsCount :one
SELECT COUNT(*) FROM workout_sessions
WHERE user_id = $1;

-- name: DeleteWorkoutSession :exec
DELETE FROM workout_sessions
WHERE id = $1 AND user_id = $2;

-- name: UpdateWorkoutSession :exec
UPDATE workout_sessions
SET workout_date = $1, description = $2, notes = $3, updated_at = NOW()
WHERE id = $4;

-- name: GetLastNWorkoutSessions :many
SELECT * FROM workout_sessions
WHERE user_id = $1
ORDER BY workout_date DESC
LIMIT $2;
