package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/KirillMironov/ai/api"
	"github.com/KirillMironov/ai/internal/logger"
	"github.com/KirillMironov/ai/internal/server"
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

	// logger
	slog.SetDefault(logger.New())

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

	grpcServer := grpc.NewServer()
	llmServer := server.NewLLM(llamaLLM)
	api.RegisterLLMServer(grpcServer, llmServer)

	go func() {
		log.Printf("starting grpc server at %s", listener.Addr())
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	<-ctx.Done()
	log.Printf("shutting down grpc server")
	grpcServer.GracefulStop()
}
