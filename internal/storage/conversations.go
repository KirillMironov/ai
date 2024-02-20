package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KirillMironov/ai/internal/model"
	"github.com/KirillMironov/ai/internal/storage/queries"
)

type Conversations struct {
	db *sql.DB
}

func NewConversations(db *sql.DB) Conversations {
	return Conversations{db: db}
}

func (c Conversations) SaveConversation(ctx context.Context, conversation model.Conversation) error {
	return queries.New(c.db).SaveConversation(ctx, queries.SaveConversationParams{
		ID:        conversation.ID,
		UserID:    conversation.UserID,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	})
}

func (c Conversations) GetConversationsByUserID(ctx context.Context, userID string, offset, limit int) ([]model.Conversation, error) {
	dataConversations, err := queries.New(c.db).GetConversationsByUserID(ctx, queries.GetConversationsByUserIDParams{
		UserID: userID,
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	conversations := make([]model.Conversation, 0, len(dataConversations))

	for _, conversation := range dataConversations {
		conversations = append(conversations, model.Conversation{
			ID:        conversation.ID,
			UserID:    conversation.UserID,
			Title:     conversation.Title,
			CreatedAt: conversation.CreatedAt,
			UpdatedAt: conversation.UpdatedAt,
		})
	}

	return conversations, nil
}

func (c Conversations) GetConversationByID(ctx context.Context, id string) (conversation model.Conversation, exists bool, err error) {
	dataConversation, err := queries.New(c.db).GetConversationByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return conversation, false, err
	}

	conversation = model.Conversation{
		ID:        dataConversation.ID,
		UserID:    dataConversation.UserID,
		Title:     dataConversation.Title,
		CreatedAt: dataConversation.CreatedAt,
		UpdatedAt: dataConversation.UpdatedAt,
	}

	return conversation, true, nil
}
