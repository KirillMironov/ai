package model

import "fmt"

type Role uint8

const (
	RoleAssistant Role = iota + 1
	RoleUser
)

func (r Role) String() string {
	switch r {
	case RoleAssistant:
		return "assistant"
	case RoleUser:
		return "user"
	default:
		return fmt.Sprintf("unexpected role (%d)", r)
	}
}
