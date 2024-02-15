package model

import "github.com/google/uuid"

type TokenPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
