-- name: SaveUser :exec
INSERT INTO users (id, username, hashed_password, created_at) VALUES (?, ?, ?, ?);

-- name: GetUserByUsername :one
SELECT id, username, hashed_password, created_at FROM users WHERE username = ?;
