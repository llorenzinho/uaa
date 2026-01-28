-- name: ListKeys :many
SELECT * FROM jwk_keys WHERE expires_at > NOW();

-- name: GetJwksKey :one
SELECT * FROM jwk_keys
WHERE kid = $1
AND expres_at > NOW();