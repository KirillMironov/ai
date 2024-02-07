package llm

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/ai/api"
	"github.com/KirillMironov/ai/llm"
)

var errEmptyPrompt = status.Error(codes.InvalidArgument, "empty prompt")

type Server struct {
	llm llm.LLM
	api.UnimplementedLLMServer
}

func NewServer(llm llm.LLM) Server {
	return Server{llm: llm}
}

func (s Server) Completion(ctx context.Context, request *api.CompletionRequest) (*api.CompletionResponse, error) {
	if request.Prompt == "" {
		return nil, errEmptyPrompt
	}

	req := llm.CompletionRequest{Prompt: request.Prompt}

	resp, err := s.llm.Completion(ctx, req)
	if err != nil {
		slog.Error("failed to call llm.Completion", err)
		return nil, err
	}

	return &api.CompletionResponse{Content: resp.Content}, nil
}

func (s Server) CompletionStream(request *api.CompletionStreamRequest, stream api.LLM_CompletionStreamServer) error {
	if request.Prompt == "" {
		return errEmptyPrompt
	}

	req := llm.CompletionRequest{Prompt: request.Prompt}

	chunkProcessor := func(response llm.CompletionResponse) {
		_ = stream.Send(&api.CompletionStreamResponse{Content: response.Content})
	}

	if err := s.llm.CompletionStream(stream.Context(), req, chunkProcessor); err != nil {
		slog.Error("failed to call llm.CompletionStream", err)
		return err
	}

	return nil
}
