-- +goose Up
CREATE TABLE IF NOT EXISTS conversations (
    id         TEXT     NOT NULL,
    user_id    TEXT     NOT NULL,
    title      TEXT     NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT pk_conversations_id PRIMARY KEY (id),
    CONSTRAINT fk_conversations_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE IF EXISTS conversations;
