-- name: SaveMessage :exec
INSERT INTO messages (id, conversation_id, role, content) VALUES (?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    conversation_id = EXCLUDED.conversation_id,
    role            = EXCLUDED.role,
    content         = EXCLUDED.content;

-- name: GetMessagesByConversationID :many
SELECT id, role, content FROM messages WHERE conversation_id = ? ORDER BY id;
