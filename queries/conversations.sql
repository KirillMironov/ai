-- name: SaveConversation :exec
INSERT INTO conversations (id, user_id, title, created_at, updated_at) VALUES (?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    user_id    = EXCLUDED.user_id,
    title      = EXCLUDED.title,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at;

-- name: GetConversationsByUserID :many
SELECT id, user_id, title, created_at, updated_at
FROM conversations
WHERE user_id = ?
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: GetConversationByID :one
SELECT id, user_id, title, created_at, updated_at FROM conversations WHERE id = ?;

-- name: DeleteConversationByID :exec
DELETE FROM conversations WHERE id = ?;
