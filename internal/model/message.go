package model

import "time"

type Message struct {
	ID        string
	Role      Role
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
