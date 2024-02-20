package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KirillMironov/ai/internal/model"
	"github.com/KirillMironov/ai/internal/storage/queries"
)

type Messages struct {
	db *sql.DB
}

func NewMessages(db *sql.DB) Messages {
	return Messages{db: db}
}

func (m Messages) SaveMessage(ctx context.Context, conversationID string, message model.Message) error {
	return queries.New(m.db).SaveMessage(ctx, queries.SaveMessageParams{
		ID:             message.ID,
		ConversationID: conversationID,
		Role:           int64(message.Role),
		Content:        message.Content,
	})
}

func (m Messages) GetMessagesByConversationID(ctx context.Context, conversationID string) ([]model.Message, error) {
	dataMessages, err := queries.New(m.db).GetMessagesByConversationID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	messages := make([]model.Message, 0, len(dataMessages))

	for _, message := range dataMessages {
		messages = append(messages, model.Message{
			ID:      message.ID,
			Role:    model.Role(uint8(message.Role)),
			Content: message.Content,
		})
	}

	return messages, nil
}
