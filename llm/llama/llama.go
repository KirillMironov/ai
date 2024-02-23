package llama

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/KirillMironov/ai/internal/httputil"
	"github.com/KirillMironov/ai/llm"
)

const (
	defaultServerHost  = "0.0.0.0"
	defaultServerPort  = 8080
	defaultContextSize = 2048
	defaultNumSlots    = 1
	defaultNumThreads  = 4
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
		numThreads:     max(runtime.NumCPU(), defaultNumThreads),
	}
	for _, option := range options {
		option(llama)
	}
	return llama
}

func (l *Llama) Completion(ctx context.Context, request llm.CompletionRequest) (llm.CompletionResponse, error) {
	req := completionRequest{Prompt: request.Prompt}

	resp, err := httputil.Post[completionResponse](ctx, l.serverURL("/completion"), http.StatusOK, req)
	if err != nil {
		return llm.CompletionResponse{}, err
	}

	return llm.CompletionResponse{Content: resp.Content}, nil
}

func (l *Llama) CompletionStream(ctx context.Context, request llm.CompletionRequest, onChunk func(llm.CompletionResponse) error) error {
	req := completionRequest{Prompt: request.Prompt, Stream: true}

	body, err := httputil.PostBody(ctx, l.serverURL("/completion"), http.StatusOK, req)
	if err != nil {
		return err
	}
	defer body.Close()

	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		var chunk completionResponseChunk
		if err = json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &chunk); err != nil {
			return err
		}

		if err = onChunk(llm.CompletionResponse{Content: chunk.Content}); err != nil {
			return err
		}
	}

	return nil
}

func (l *Llama) ChatCompletion(ctx context.Context, request llm.ChatCompletionRequest) (llm.ChatCompletionResponse, error) {
	req := chatCompletionRequest{Messages: messagesToLlamaMessages(request.Messages)}

	resp, err := httputil.Post[chatCompletionResponse](ctx, l.serverURL("/v1/chat/completions"), http.StatusOK, req)
	if err != nil {
		return llm.ChatCompletionResponse{}, err
	}

	if len(resp.Choices) == 0 {
		return llm.ChatCompletionResponse{}, nil
	}

	return llm.ChatCompletionResponse{Message: messageFromLlamaMessage(resp.Choices[0].Message)}, nil
}

func (l *Llama) ChatCompletionStream(ctx context.Context, request llm.ChatCompletionRequest, onChunk func(llm.ChatCompletionResponse) error) error {
	req := chatCompletionRequest{Messages: messagesToLlamaMessages(request.Messages), Stream: true}

	body, err := httputil.PostBody(ctx, l.serverURL("/v1/chat/completions"), http.StatusOK, req)
	if err != nil {
		return err
	}
	defer body.Close()

	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		var chunk chatCompletionResponseChunk
		if err = json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &chunk); err != nil {
			return err
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		response := llm.ChatCompletionResponse{Message: messageFromLlamaMessage(chunk.Choices[0].Message)}

		if err = onChunk(response); err != nil {
			return err
		}
	}

	return nil
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

func messagesToLlamaMessages(messages []llm.Message) []message {
	llamaMessages := make([]message, 0, len(messages))
	for _, msg := range messages {
		llamaMessages = append(llamaMessages, message{
			Role:    roleToLlamaRole(msg.Role),
			Content: msg.Content,
		})
	}
	return llamaMessages
}

func messageFromLlamaMessage(message message) llm.Message {
	return llm.Message{
		Role:    roleFromLlamaRole(message.Role),
		Content: message.Content,
	}
}

func roleToLlamaRole(role llm.Role) role {
	switch role {
	case llm.RoleLLM:
		return roleAssistant
	case llm.RoleUser:
		return roleUser
	default:
		return 0
	}
}

func roleFromLlamaRole(role role) llm.Role {
	switch role {
	case roleAssistant:
		return llm.RoleLLM
	case roleUser:
		return llm.RoleUser
	default:
		return 0
	}
}
