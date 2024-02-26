package model

import "time"

type Conversation struct {
	ID        string
	UserID    string
	Title     string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}
