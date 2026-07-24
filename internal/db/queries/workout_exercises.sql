-- name: CreateWorkoutExercises :one
INSERT INTO workout_exercises(id, workout_session_id, exercise_id, order_index, notes, created_at, updated_at)
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

-- name: GetWorkoutExerciseFromID :one
SELECT * FROM workout_exercises
WHERE id = $1;

-- name: GetWorkoutsInSession :many
SELECT * FROM workout_exercises
WHERE workout_session_id = $1
ORDER BY order_index;

-- name: GetNumberOfExercisesInWorkout :one
SELECT COUNT(*) FROM workout_exercises
WHERE workout_session_id = $1;

-- name: UpdateWorkoutExerciseOrder :one
UPDATE workout_exercises
SET order_index = $1, updated_at = NOW()
WHERE id = $2 AND workout_session_id = $3
RETURNING *;

-- name: DeleteWorkoutExercises :exec
DELETE FROM workout_exercises
WHERE id = $1 AND workout_session_id = $2;

-- name: UpdateWorkoutExercise :one
UPDATE workout_exercises
SET exercise_id = $1, order_index = $2, notes = $3, updated_at = NOW()
WHERE id = $4 AND workout_session_id = $5
RETURNING *;
