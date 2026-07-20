-- name: CreateCardioSet :one
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

-- name: GetCardioSetFromID :one
SELECT * FROM cardio_and_sets
WHERE id = $1;

-- name: GetAllCardioFromExercise :many
SELECT * FROM cardio_and_sets
WHERE workout_exercises_id = $1;

-- name: GetSetDistance :one
SELECT distance FROM cardio_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetAllSetsDistance :one
SELECT COALESCE(SUM(distance), 0)::float8
FROM cardio_and_sets
WHERE workout_exercises_id = $1;

-- name: GetCardioSetDuration :one
SELECT duration_seconds FROM cardio_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: GetAllCardioSetsDuration :one
SELECT COALESCE(SUM(duration_seconds),0)::bigint
FROM cardio_and_sets
WHERE workout_exercises_id = $1;

-- name: DeleteCardioSet :exec
DELETE FROM cardio_and_sets
WHERE id = $1 AND workout_exercises_id = $2;

-- name: UpdateCardioSet :one
UPDATE cardio_and_sets
SET
    set_number = $1,
    distance = $2,
    is_kilometers = $3,
    duration_seconds = $4,
    updated_at = NOW()
WHERE id = $5 AND workout_exercises_id = $6
RETURNING *;

-- name: GetAllCardioSetsFromSession :many
SELECT cardio_and_sets.*
FROM cardio_and_sets
JOIN workout_exercises
ON cardio_and_sets.workout_exercises_id = workout_exercises.id
WHERE workout_exercises.workout_session_id = $1
ORDER BY workout_exercises.order_index, cardio_and_sets.set_number;

-- name: GetTotalSessionCardioDuration :one
SELECT COALESCE(SUM(cardio_and_sets.duration_seconds), 0)::bigint
FROM cardio_and_sets
JOIN workout_exercises
ON cardio_and_sets.workout_exercises_id = workout_exercises.id
WHERE workout_exercises.workout_session_id = $1;

-- name: GetTotalSessionDistance :one
SELECT COALESCE(SUM(cardio_and_sets.distance), 0)::float8
FROM cardio_and_sets
JOIN workout_exercises
ON cardio_and_sets.workout_exercises_id = workout_exercises.id
WHERE workout_exercises.workout_session_id = $1;
