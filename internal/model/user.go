package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}
