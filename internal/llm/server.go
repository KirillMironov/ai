package llm

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/ai/api"
	"github.com/KirillMironov/ai/llm"
)

type Server struct {
	llm llm.LLM
	api.UnimplementedLLMServer
}

func NewServer(llm llm.LLM) Server {
	return Server{llm: llm}
}

func (s Server) Completion(ctx context.Context, request *api.CompletionRequest) (*api.CompletionResponse, error) {
	if request.Prompt == "" {
		return nil, status.Error(codes.InvalidArgument, "empty prompt")
	}

	req := llm.CompletionRequest{Prompt: request.Prompt}

	resp, err := s.llm.Completion(ctx, req)
	if err != nil {
		return nil, err
	}

	return &api.CompletionResponse{Content: resp.Content}, nil
}
