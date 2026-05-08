-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(hashed_token, user_id, expires_at, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE hashed_token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET is_revoked = true, updated_at = NOW()
WHERE hashed_token = $1;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET is_revoked = true, updated_at = NOW()
WHERE user_id = $1;
