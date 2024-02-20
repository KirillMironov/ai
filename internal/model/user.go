package model

import "time"

type User struct {
	ID             string
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}
