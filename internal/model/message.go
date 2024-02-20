package model

type Message struct {
	ID      string
	Role    Role
	Content string
}

type SendMessageRequest struct {
	ConversationID string
	Role           Role
	Content        string
}
