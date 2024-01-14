package llama

import (
	"context"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"sync"

	"github.com/KirillMironov/ai/internal/httputil"
	"github.com/KirillMironov/ai/llm"
)

const (
	defaultServerHost  = "0.0.0.0"
	defaultServerPort  = 8080
	defaultContextSize = 2048
	defaultNumSlots    = 1
)

var _ llm.LLM = &Llama{}

type Llama struct {
	executablePath string
	modelPath      string
	cancel         context.CancelFunc
	once           sync.Once

	serverHost  string // llama server host
	serverPort  int    // llama server port
	contextSize int    // prompt context size
	numSlots    int    // number of slot to process requests
	numThreads  int    // number of threads to use during generation
	mmap        bool   // memory-map the model to load only necessary parts of it as needed
}

func New(executablePath, modelPath string, options ...Option) *Llama {
	llama := &Llama{
		executablePath: executablePath,
		modelPath:      modelPath,
		serverHost:     defaultServerHost,
		serverPort:     defaultServerPort,
		contextSize:    defaultContextSize,
		numSlots:       defaultNumSlots,
		numThreads:     max(runtime.NumCPU(), 4),
	}
	for _, option := range options {
		option(llama)
	}
	return llama
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

		args := []string{
			"--model", l.modelPath,
			"--host", l.serverHost,
			"--port", strconv.Itoa(l.serverPort),
			"--ctx-size", strconv.Itoa(l.contextSize),
			"--threads", strconv.Itoa(l.numThreads),
			"--parallel", strconv.Itoa(l.numSlots),
		}

		if !l.mmap {
			args = append(args, "--no-mmap")
		} else {
			args = append(args, "--mlock")
		}

		cmd := exec.CommandContext(ctx, l.executablePath, args...)
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
	return "http://" + net.JoinHostPort("localhost", strconv.Itoa(l.serverPort)) + path
}

type (
	completionRequest struct {
		Prompt string `json:"prompt"`
	}

	completionResponse struct {
		Content string `json:"content"`
	}
)
