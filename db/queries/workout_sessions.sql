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

-- name: GetWorkoutSessionsCount :one
SELECT COUNT(*) FROM workout_sessions
WHERE user_id = $1;

-- name: DeleteWorkoutSession :exec
DELETE FROM workout_sessions
WHERE id = $1 AND user_id = $2;
