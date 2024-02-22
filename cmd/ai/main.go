package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "modernc.org/sqlite"

	api "github.com/KirillMironov/ai/internal/api/ai"
	llmapi "github.com/KirillMironov/ai/internal/api/llm"
	"github.com/KirillMironov/ai/internal/logger"
	"github.com/KirillMironov/ai/internal/model"
	"github.com/KirillMironov/ai/internal/server"
	"github.com/KirillMironov/ai/internal/service"
	"github.com/KirillMironov/ai/internal/storage"
	"github.com/KirillMironov/ai/internal/token"
	"github.com/KirillMironov/ai/migrations"
)

type config struct {
	Port int `envconfig:"PORT" default:"8080"`

	SQLite struct {
		Path string `envconfig:"SQLITE_PATH" required:"true"`
	}

	JWT struct {
		Secret   string        `envconfig:"JWT_SECRET" required:"true"`
		TokenTTL time.Duration `envconfig:"JWT_TOKEN_TTL" default:"24h"`
	}

	LLM struct {
		Address  string `envconfig:"LLM_ADDRESS" required:"true"`
		Insecure bool   `envconfig:"LLM_INSECURE" default:"true"`
	}
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// logger
	slog.SetDefault(logger.New())

	// config
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}

	// db
	if err := os.MkdirAll(filepath.Dir(cfg.SQLite.Path), os.ModePerm); err != nil {
		return err
	}

	dataSourceURI := fmt.Sprintf(`file:%s?_pragma=foreign_keys(1)&_time_format=sqlite`, cfg.SQLite.Path)

	db, err := sql.Open("sqlite", dataSourceURI)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return err
	}

	// migrations
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrations.Embed)
	if err != nil {
		return err
	}

	if _, err = provider.Up(ctx); err != nil {
		return err
	}

	// storage
	usersStorage := storage.NewUsers(db)
	conversationsStorage := storage.NewConversations(db)
	messagesStorage := storage.NewMessages(db)

	// llm client
	var opts []grpc.DialOption
	if cfg.LLM.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(cfg.LLM.Address, opts...)
	if err != nil {
		return err
	}

	llmClient := llmapi.NewLLMClient(conn)

	// services
	tokenManager := token.NewManager[model.TokenPayload]([]byte(cfg.JWT.Secret), cfg.JWT.TokenTTL)
	authenticator := service.NewAuthenticator(usersStorage, tokenManager)
	conversations := service.NewConversations(authenticator, conversationsStorage, messagesStorage, llmClient)

	// grpc server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	authenticatorServer := server.NewAuthenticator(authenticator)
	conversationsServer := server.NewConversations(conversations)
	api.RegisterAuthenticatorServer(grpcServer, authenticatorServer)
	api.RegisterConversationsServer(grpcServer, conversationsServer)

	go func() {
		slog.Info("starting grpc server", slog.String("address", listener.Addr().String()))
		if err = grpcServer.Serve(listener); err != nil {
			slog.Error(err.Error())
			cancel()
		}
	}()

	// graceful shutdown
	<-ctx.Done()
	slog.Info("shutting down grpc server")
	grpcServer.GracefulStop()

	return nil
}
