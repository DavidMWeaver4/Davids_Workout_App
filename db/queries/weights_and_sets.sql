-- name: CreateWeightsAndSets :one
INSERT INTO weights_and_sets(id, workout_exercises_id, weight, is_kilograms, set_number, reps_target, reps_actual, duration_seconds, rest_time_seconds, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11
)
RETURNING *;

-- name: GetWeightAndSetFromID :one
SELECT * FROM weights_and_sets
WHERE id = $1;

-- name: GetAllSetsFromSession :many
SELECT * FROM weights_and_sets
WHERE workout_exercises_id = $1;

-- name: GetSetVolume :one
SELECT reps_actual * weight FROM weights_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetTotalVolumeFromAllSets :one
SELECT SUM(reps_actual * weight) FROM weights_and_sets
WHERE workout_exercises_id = $1;

-- name: GetTotalDuration :one
SELECT duration_seconds + rest_time_seconds FROM weights_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetTotalDurationForAllSets :one
SELECT SUM(duration_seconds + rest_time_seconds) FROM weights_and_sets
WHERE workout_exercises_id = $1;

-- name: DeleteWeightAndSets :exec
DELETE FROM weights_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: UpdateWeightsAndSets :one
UPDATE weights_and_sets
SET
    weight = $1,
    is_kilograms = $2,
    set_number = $3,
    reps_target = $4,
    reps_actual = $5,
    duration_seconds = $6,
    rest_time_seconds = $7,
    updated_at = NOW()
WHERE id = $8 AND workout_exercises_id = $9
RETURNING *;
