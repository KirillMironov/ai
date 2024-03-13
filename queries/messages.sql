-- name: SaveMessage :exec
INSERT INTO messages (id, conversation_id, role, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    conversation_id = EXCLUDED.conversation_id,
    role            = EXCLUDED.role,
    content         = EXCLUDED.content,
    created_at      = EXCLUDED.created_at,
    updated_at      = EXCLUDED.updated_at;

-- name: GetMessagesByConversationID :many
SELECT id, role, content, created_at, updated_at FROM messages WHERE conversation_id = ? ORDER BY created_at;

-- name: DeleteMessagesByConversationID :exec
DELETE FROM messages WHERE conversation_id = ?;
