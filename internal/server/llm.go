package server

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/KirillMironov/ai/internal/api/llm"
	"github.com/KirillMironov/ai/llm"
)

var (
	errEmptyPrompt = status.Error(codes.InvalidArgument, "empty prompt")
	errNoMessages  = status.Error(codes.InvalidArgument, "no messages")
)

type LLM struct {
	llm llm.LLM
	api.UnimplementedLLMServer
}

func NewLLM(llm llm.LLM) LLM {
	return LLM{llm: llm}
}

func (l LLM) Completion(ctx context.Context, request *api.CompletionRequest) (*api.CompletionResponse, error) {
	prompt := request.GetPrompt()
	if prompt == "" {
		return nil, errEmptyPrompt
	}

	req := llm.CompletionRequest{Prompt: prompt}

	resp, err := l.llm.Completion(ctx, req)
	if err != nil {
		slog.Error("failed to call llm.Completion", err)
		return nil, err
	}

	return &api.CompletionResponse{Content: resp.Content}, nil
}

func (l LLM) CompletionStream(request *api.CompletionStreamRequest, stream api.LLM_CompletionStreamServer) error {
	prompt := request.GetPrompt()
	if prompt == "" {
		return errEmptyPrompt
	}

	req := llm.CompletionRequest{Prompt: prompt}

	onChunk := func(response llm.CompletionResponse) error {
		return stream.Send(&api.CompletionStreamResponse{Content: response.Content})
	}

	if err := l.llm.CompletionStream(stream.Context(), req, onChunk); err != nil {
		slog.Error("failed to call llm.CompletionStream", err)
		return err
	}

	return nil
}

func (l LLM) ChatCompletion(ctx context.Context, request *api.ChatCompletionRequest) (*api.ChatCompletionResponse, error) {
	apiMessages := request.GetMessages()
	if len(apiMessages) == 0 {
		return nil, errNoMessages
	}

	messages, err := messagesFromAPI(apiMessages)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := l.llm.ChatCompletion(ctx, llm.ChatCompletionRequest{Messages: messages})
	if err != nil {
		slog.Error("failed to call llm.ChatCompletion", err)
		return nil, err
	}

	return &api.ChatCompletionResponse{Message: messageToAPI(resp.Message)}, nil
}

func (l LLM) ChatCompletionStream(request *api.ChatCompletionStreamRequest, stream api.LLM_ChatCompletionStreamServer) error {
	apiMessages := request.GetMessages()
	if len(apiMessages) == 0 {
		return errNoMessages
	}

	messages, err := messagesFromAPI(apiMessages)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	req := llm.ChatCompletionRequest{Messages: messages}

	onChunk := func(response llm.ChatCompletionResponse) error {
		return stream.Send(&api.ChatCompletionStreamResponse{Message: messageToAPI(response.Message)})
	}

	if err = l.llm.ChatCompletionStream(stream.Context(), req, onChunk); err != nil {
		slog.Error("failed to call llm.ChatCompletionStream", err)
		return err
	}

	return nil
}

func messagesFromAPI(apiMessages []*api.Message) ([]llm.Message, error) {
	messages := make([]llm.Message, 0, len(apiMessages))
	for _, message := range apiMessages {
		var role llm.Role
		switch apiRole := message.GetRole(); apiRole {
		case api.Role_ROLE_LLM:
			role = llm.RoleLLM
		case api.Role_ROLE_USER:
			role = llm.RoleUser
		default:
			return nil, fmt.Errorf("unexpected message role: '%v'", apiRole)
		}
		messages = append(messages, llm.Message{
			Role:    role,
			Content: message.GetContent(),
		})
	}
	return messages, nil
}

func messageToAPI(message llm.Message) *api.Message {
	var apiRole api.Role
	switch message.Role {
	case llm.RoleLLM:
		apiRole = api.Role_ROLE_LLM
	case llm.RoleUser:
		apiRole = api.Role_ROLE_USER
	default:
		apiRole = api.Role_ROLE_UNSPECIFIED
	}
	return &api.Message{
		Role:    apiRole,
		Content: message.Content,
	}
}
