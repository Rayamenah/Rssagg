-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_key)
VALUES ($1, $2, $3, $4, encode(sha2(UUID(), 256), 'hex'))
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1