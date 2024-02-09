package llm

import "context"

type LLM interface {
	Completion(ctx context.Context, request CompletionRequest) (CompletionResponse, error)
	CompletionStream(ctx context.Context, request CompletionRequest, onChunk func(CompletionResponse) error) error
	Start(ctx context.Context) error
	Close(ctx context.Context) error
}

type (
	CompletionRequest struct {
		Prompt  string
		Options Options
	}

	CompletionResponse struct {
		Content string
		Stats   Stats
	}

	// CompletionRequest struct {
	//	Messages []Message
	//	Options  Options
	// }

	// CompletionResponse struct {
	//	Message Message
	//	Stats   Stats
	// }

	Message struct {
		Role    Role
		Content string
	}

	Options map[string]string

	Stats map[string]string
)

type Role uint8

const (
	RoleLLM Role = iota + 1
	RoleUser
)
