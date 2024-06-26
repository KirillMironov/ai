// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package queries

import (
	"context"
	"time"
)

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, hashed_password, created_at FROM users WHERE username = ?
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
	)
	return i, err
}

const saveUser = `-- name: SaveUser :exec
INSERT INTO users (id, username, hashed_password, created_at) VALUES (?, ?, ?, ?)
`

type SaveUserParams struct {
	ID             string
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}

func (q *Queries) SaveUser(ctx context.Context, arg SaveUserParams) error {
	_, err := q.db.ExecContext(ctx, saveUser,
		arg.ID,
		arg.Username,
		arg.HashedPassword,
		arg.CreatedAt,
	)
	return err
}
