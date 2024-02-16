package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KirillMironov/ai/internal/model"
	"github.com/KirillMironov/ai/internal/storage/queries"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) Users {
	return Users{db: db}
}

func (u Users) SaveUser(ctx context.Context, user model.User) error {
	return queries.New(u.db).SaveUser(ctx, queries.SaveUserParams{
		ID:             user.ID,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
	})
}

func (u Users) GetUserByUsername(ctx context.Context, username string) (user model.User, exists bool, err error) {
	dataUser, err := queries.New(u.db).GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return model.User{}, false, err
	}

	user = model.User{
		ID:             dataUser.ID,
		Username:       dataUser.Username,
		HashedPassword: dataUser.HashedPassword,
		CreatedAt:      dataUser.CreatedAt,
	}

	return user, true, nil
}
