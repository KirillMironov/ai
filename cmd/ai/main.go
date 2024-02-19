package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"

	"github.com/KirillMironov/ai/api"
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

	// authenticator
	tokenManager := token.NewManager[model.TokenPayload]([]byte(cfg.JWT.Secret), cfg.JWT.TokenTTL)
	authenticator := service.NewAuthenticator(usersStorage, tokenManager)

	// grpc server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	authenticatorServer := server.NewAuthenticator(authenticator)
	api.RegisterAuthenticatorServer(grpcServer, authenticatorServer)

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
