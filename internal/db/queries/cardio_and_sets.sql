-- name: CreateCardioAndSets :one
INSERT INTO cardio_and_sets(id, workout_exercises_id, set_number, distance, is_kilometers, duration_seconds, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetCardioAndSetFromID :one
SELECT * FROM cardio_and_sets
WHERE id = $1;

-- name: GetAllCardioFromSession :many
SELECT * FROM cardio_and_sets
WHERE workout_exercises_id = $1;

-- name: GetSetDistance :one
SELECT distance FROM cardio_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetAllSetsDistance :one
SELECT SUM(distance) FROM cardio_and_sets
WHERE workout_exercises_id = $1;

-- name: GetSetDuration :one
SELECT duration_seconds FROM cardio_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetAllSetsDuration :one
SELECT SUM(duration_seconds) FROM cardio_and_sets
WHERE workout_exercises_id = $1;

-- name: DeleteCardioAndSets :exec
DELETE FROM cardio_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: UpdateCardioAndSets :one
UPDATE cardio_and_sets
SET
    set_number = $1,
    distance = $2,
    is_kilometers = $3,
    duration_seconds = $4,
    updated_at = NOW()
WHERE id = $5 AND workout_exercises_id = $6
RETURNING *;
