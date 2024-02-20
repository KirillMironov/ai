-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
    id              TEXT    NOT NULL,
    conversation_id TEXT    NOT NULL,
    role            INTEGER NOT NULL,
    content         TEXT    NOT NULL,
    CONSTRAINT pk_messages_id PRIMARY KEY (id),
    CONSTRAINT fk_messages_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversations (id)
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_conversation_id;

DROP TABLE IF EXISTS messages;
