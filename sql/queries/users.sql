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
