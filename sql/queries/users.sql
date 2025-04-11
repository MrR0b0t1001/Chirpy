-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4, 
    $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE users.email = $1; 

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY chirps.created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE chirps.ID = $1; 

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5, 
    $6
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT *
FROM refresh_tokens 
WHERE refresh_tokens.token = $1;

-- name: RevokeUserToken :exec 
UPDATE refresh_tokens 
SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE refresh_tokens.token = $1;

-- name: UpdateUserCreds :exec 
UPDATE users 
SET email = $1, hashed_password = $2, updated_at = CURRENT_TIMESTAMP
WHERE users.id = $3;
