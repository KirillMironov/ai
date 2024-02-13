package llm

import "context"

type LLM interface {
	Completion(ctx context.Context, request CompletionRequest) (CompletionResponse, error)
	CompletionStream(ctx context.Context, request CompletionRequest, onChunk func(CompletionResponse) error) error
	ChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, request ChatCompletionRequest, onChunk func(ChatCompletionResponse) error) error
	Start(ctx context.Context) error
	Close(ctx context.Context) error
}

type (
	CompletionRequest struct {
		Prompt string
	}

	CompletionResponse struct {
		Content string
	}

	ChatCompletionRequest struct {
		Messages []Message
	}

	ChatCompletionResponse struct {
		Message Message
	}

	Message struct {
		Role    Role
		Content string
	}
)

type Role uint8

const (
	RoleLLM Role = iota + 1
	RoleUser
)
