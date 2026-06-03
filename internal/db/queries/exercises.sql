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

-- name: GetExercisesByMuscles :many
SELECT *
FROM exercises
WHERE target_muscles && $1::text[] AND (user_id IS NULL OR user_id = $2);

-- name: GetSameMuscleExercises :many
SELECT * FROM exercises
WHERE $1 = ANY(target_muscles) AND (user_id IS NULL OR user_id = $2);

-- name: GetSameEquipmentExercises :many
SELECT * FROM exercises
WHERE equipment = $1 AND (user_id IS NULL OR user_id = $2);

-- name: GetSameDifficultyExercises :many
SELECT * FROM exercises
WHERE difficulty_level = $1 AND (user_id IS NULL OR user_id = $2);

-- name: DeleteExerciseByID :exec
DELETE FROM exercises
WHERE id = $1 AND user_id = $2;

-- name: ListAvailableExercises :many
SELECT * FROM exercises
WHERE user_id IS NULL OR user_id = $1
ORDER BY exercise_name;

-- name: UpdateExercise :one
UPDATE exercises
SET exercise_name = $1,
    target_muscles = $2,
    equipment = $3,
    difficulty_level = $4,
    description = $5,
    updated_at = NOW()
WHERE id = $6 AND user_id = $7
RETURNING *;

-- name: SearchExercises :many
SELECT * FROM exercises
WHERE (user_id IS NULL OR user_id = $1)
AND ($2::text IS NULL OR difficulty_level = $2)
AND ($3::text IS NULL OR equipment = $3)
AND ($4::text[] IS NULL OR target_muscles && $4::text[])
ORDER BY exercise_name;
