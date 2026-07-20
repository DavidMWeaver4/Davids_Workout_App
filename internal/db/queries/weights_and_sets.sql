-- name: CreateWeightAndSet :one
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

-- name: GetAllSetsFromExercise :many
SELECT * FROM weights_and_sets
WHERE workout_exercises_id = $1;

-- name: GetSetVolume :one
SELECT COALESCE(reps_actual * weight, 0)::float8
FROM weights_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetTotalVolumeFromExerciseSets :one
SELECT COALESCE(SUM(reps_actual * weight), 0)::float8 FROM weights_and_sets
WHERE workout_exercises_id = $1;

-- name: GetTotalDuration :one
SELECT COALESCE(duration_seconds, 0) + COALESCE(rest_time_seconds, 0) FROM weights_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetTotalDurationForExercise :one
SELECT COALESCE(SUM(COALESCE(duration_seconds, 0) + COALESCE(rest_time_seconds, 0)), 0)::int4 FROM weights_and_sets
WHERE workout_exercises_id = $1;

-- name: DeleteWeightAndSet :exec
DELETE FROM weights_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: UpdateWeightAndSet :one
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

-- name: GetAllWeightSetsFromSession :many
SELECT weights_and_sets.*
FROM weights_and_sets
JOIN workout_exercises
ON weights_and_sets.workout_exercises_id = workout_exercises.id
WHERE workout_exercises.workout_session_id = $1
ORDER BY workout_exercises.order_index, weights_and_sets.set_number;

-- name: GetTotalSessionVolume :one
SELECT COALESCE(SUM(weights_and_sets.weight * weights_and_sets.reps_actual), 0)::float8
FROM weights_and_sets
JOIN workout_exercises
ON weights_and_sets.workout_exercises_id = workout_exercises.id
WHERE workout_exercises.workout_session_id = $1;

-- name: GetTotalSessionDuration :one
SELECT COALESCE(SUM(COALESCE(duration_seconds,0) + COALESCE(rest_time_seconds,0)), 0)::int4
FROM weights_and_sets
JOIN workout_exercises
ON weights_and_sets.workout_exercises_id = workout_exercises.id
WHERE workout_exercises.workout_session_id = $1;
