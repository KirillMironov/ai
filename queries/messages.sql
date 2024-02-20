-- name: SaveMessage :exec
INSERT INTO messages (id, conversation_id, role, content) VALUES (?, ?, ?, ?);

-- name: GetMessagesByConversationID :many
SELECT id, conversation_id, role, content FROM messages WHERE conversation_id = ? ORDER BY id;
