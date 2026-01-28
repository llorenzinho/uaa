-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2, email_verified = FALSE, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: VerifyUserEmail :one
UPDATE users
SET email_verified = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPasswordHash :one
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;