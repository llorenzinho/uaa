-- name: ListKeys :many
SELECT * FROM jwk_keys WHERE expires_at > NOW();

-- name: GetActiveJwk :one
SELECT * FROM jwk_keys WHERE expires_at > NOW() AND is_active = 1;

-- name: GetJwksKey :one
SELECT * FROM jwk_keys
WHERE kid = $1
AND expres_at > NOW();

-- name: CreateNewRs256Key :one
INSERT INTO jwk_keys (kid, private_key_pem, public_key_pem, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ActivateKey :many
UPDATE jwk_keys
SET is_active = (kid = $1)
RETURNING *;

-- name: DeleteExpiredKey :many
DELETE FROM jwk_keys WHERE expires_at <= NOW()
RETURNING *;