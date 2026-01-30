-- name: CreateAuthorizationCode :one
INSERT INTO authorization_codes
(code, user_id, client_id, redirect_uri, scope, code_challenge, expires_at, created_at)
VALUES
($1, $2, $3, $4, $5, $6, $7, NOW())
RETURNING *;

-- name: UseAuthorizationCode :exec
UPDATE authorization_codes
SET used = 1
WHERE code = $1;

-- name: DeleteExpiredAuthorizationCodes :many
DELETE FROM authorization_codes 
WHERE expires_at < NOW()
RETURNING *;