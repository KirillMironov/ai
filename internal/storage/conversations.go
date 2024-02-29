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
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := queries.New(tx)

	if err = qtx.SaveConversation(ctx, queries.SaveConversationParams{
		ID:        conversation.ID,
		UserID:    conversation.UserID,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	}); err != nil {
		return err
	}

	for _, message := range conversation.Messages {
		if err = qtx.SaveMessage(ctx, queries.SaveMessageParams{
			ID:             message.ID,
			ConversationID: conversation.ID,
			Role:           int64(message.Role),
			Content:        message.Content,
			CreatedAt:      message.CreatedAt,
			UpdatedAt:      message.UpdatedAt,
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
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
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Conversation{}, false, err
	}
	defer tx.Rollback()

	qtx := queries.New(tx)

	dataConversation, err := qtx.GetConversationByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return model.Conversation{}, false, err
	}

	conversation = model.Conversation{
		ID:        dataConversation.ID,
		UserID:    dataConversation.UserID,
		Title:     dataConversation.Title,
		CreatedAt: dataConversation.CreatedAt,
		UpdatedAt: dataConversation.UpdatedAt,
	}

	dataMessages, err := qtx.GetMessagesByConversationID(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.Conversation{}, false, err
	}

	for _, message := range dataMessages {
		conversation.Messages = append(conversation.Messages, model.Message{
			ID:        message.ID,
			Role:      model.Role(message.Role),
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		})
	}

	return conversation, true, nil
}
