package model

type SendMessageRequest struct {
	Token          string
	ConversationID string
	Content        string
}
