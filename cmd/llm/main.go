package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/KirillMironov/ai/api"
	"github.com/KirillMironov/ai/internal/llm"
	"github.com/KirillMironov/ai/llm/llama"
)

type config struct {
	Port int `envconfig:"PORT" default:"8081"`

	Llama struct {
		ExecutablePath string `envconfig:"LLAMA_EXECUTABLE_PATH" required:"true"`
		ModelPath      string `envconfig:"LLAMA_MODEL_PATH" required:"true"`
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// config
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	// llm
	llamaLLM := llama.New(cfg.Llama.ExecutablePath, cfg.Llama.ModelPath)

	if err := llamaLLM.Start(ctx); err != nil {
		log.Fatalf("failed to start llama llm: %v", err)
	}
	defer llamaLLM.Close(ctx)

	// grpc server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	llmServer := llm.NewServer(llamaLLM)
	api.RegisterLLMServer(server, llmServer)

	go func() {
		log.Printf("starting grpc server at %s", listener.Addr())
		if err = server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	<-ctx.Done()
	log.Printf("shutting down grpc server")
	server.GracefulStop()
}
