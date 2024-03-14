package llama

import (
	"encoding/json"
	"fmt"
)

const unexpectedRole = "unexpected role: '%v'"

type (
	completionRequest struct {
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream,omitempty"`
	}

	completionResponse struct {
		Content string `json:"content"`
	}

	completionResponseChunk struct {
		Content string `json:"content"`
	}
)

type (
	chatCompletionRequest struct {
		Messages    []message `json:"messages"`
		Stream      bool      `json:"stream,omitempty"`
		CachePrompt bool      `json:"cache_prompt,omitempty"`
	}

	chatCompletionResponse struct {
		Choices []struct {
			Message message `json:"message"`
		} `json:"choices"`
	}

	chatCompletionResponseChunk struct {
		Choices []struct {
			Message message `json:"delta"`
		} `json:"choices"`
	}

	message struct {
		Role    role   `json:"role"`
		Content string `json:"content"`
	}
)

type role uint8

const (
	roleUser role = iota + 1
	roleAssistant
)

func (r *role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "user":
		*r = roleUser
	case "assistant":
		*r = roleAssistant
	default:
		return fmt.Errorf(unexpectedRole, s)
	}
	return nil
}

func (r role) MarshalJSON() ([]byte, error) {
	var s string
	switch r {
	case roleUser:
		s = "user"
	case roleAssistant:
		s = "assistant"
	default:
		return nil, fmt.Errorf(unexpectedRole, r)
	}
	return json.Marshal(s)
}
