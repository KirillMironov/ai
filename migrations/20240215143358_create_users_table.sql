-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id              TEXT     NOT NULL,
    username        TEXT     NOT NULL,
    hashed_password TEXT     NOT NULL,
    created_at      DATETIME NOT NULL,
    CONSTRAINT pk_users_id PRIMARY KEY (id),
    CONSTRAINT uk_users_username UNIQUE (username)
);

-- +goose Down
DROP TABLE IF EXISTS users;
