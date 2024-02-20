-- name: SaveConversation :exec
INSERT INTO conversations (id, user_id, title, created_at, updated_at) VALUES (?, ?, ?, ?, ?);

-- name: GetConversationsByUserID :many
SELECT id, user_id, title, created_at, updated_at
FROM conversations
WHERE user_id = ?
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: GetConversationByID :one
SELECT id, user_id, title, created_at, updated_at FROM conversations WHERE id = ?;
