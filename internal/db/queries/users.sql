-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetUserFromID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUsersEmails :many
SELECT email From users;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $1, updated_at = NOW()
WHERE id = $2;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1, updated_at = NOW()
WHERE id = $2;
