package model

type SendMessageRequest struct {
	Token          string
	ConversationID string
	Role           Role
	Content        string
}
