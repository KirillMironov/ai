package llama

import (
	"context"
	"net"
	"os/exec"
	"sync"

	"github.com/KirillMironov/ai/internal/httputil"
	"github.com/KirillMironov/ai/llm"
)

var _ llm.LLM = &Llama{}

type Llama struct {
	executablePath string
	modelPath      string
	serverPort     string
	cancel         context.CancelFunc
	once           sync.Once
}

func New(executablePath, modelPath string) *Llama {
	return &Llama{
		executablePath: executablePath,
		modelPath:      modelPath,
		serverPort:     "8080",
	}
}

func (l *Llama) Completion(ctx context.Context, request llm.CompletionRequest) (llm.CompletionResponse, error) {
	req := completionRequest{Prompt: request.Prompt}

	resp, err := httputil.Post[completionRequest, completionResponse](ctx, l.serverURL("/completion"), req)
	if err != nil {
		return llm.CompletionResponse{}, err
	}

	return llm.CompletionResponse{Content: resp.Content}, nil
}

func (l *Llama) Start(ctx context.Context) error {
	var err error
	l.once.Do(func() {
		ctx, l.cancel = context.WithCancel(ctx)
		cmd := exec.CommandContext(ctx, l.executablePath, "-m", l.modelPath, "--host", "0.0.0.0", "--port", l.serverPort, "--ctx-size", "2048", "-t", "8", "--parallel", "1", "--mlock") // todo: opts with defaults
		err = cmd.Start()
	})
	return err
}

func (l *Llama) Close(_ context.Context) error {
	if l.cancel != nil {
		l.cancel()
	}
	return nil
}

func (l *Llama) serverURL(path string) string {
	return "http://" + net.JoinHostPort("localhost", l.serverPort) + path
}

type (
	completionRequest struct {
		Prompt string `json:"prompt"`
	}

	completionResponse struct {
		Content string `json:"content"`
	}
)
