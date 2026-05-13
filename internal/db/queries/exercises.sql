-- name: CreateExercises :one
INSERT INTO exercises(id, user_id, exercise_name, target_muscles, equipment, difficulty_level, description, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;

-- name: GetExerciseFromID :one
SELECT * FROM exercises
WHERE id = $1;

-- name: GetExerciseFromName :one
SELECT * FROM exercises
WHERE exercise_name = $1;

-- name: GetExercisesFromUser :many
SELECT * FROM exercises
WHERE user_id = $1;

-- name: GetSameMuscleExercises :many
SELECT * FROM exercises
WHERE $1 = ANY(target_muscles);

-- name: GetSameEquipmentExercises :many
SELECT * FROM exercises
WHERE equipment = $1;

-- name: GetSameDifficultyExercises :many
SELECT * FROM exercises
WHERE difficulty_level = $1;

-- name: DeleteExerciseByID :exec
DELETE FROM exercises
WHERE id = $1;
