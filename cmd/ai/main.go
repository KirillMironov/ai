package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"

	"github.com/KirillMironov/ai/api"
	"github.com/KirillMironov/ai/internal/ai"
	"github.com/KirillMironov/ai/internal/model"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// config
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	// db
	dataSourceURI := fmt.Sprintf(`file:%s?_pragma=foreign_keys(1)&_time_format=sqlite`, cfg.SQLite.Path)

	db, err := sql.Open("sqlite", dataSourceURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// migrations
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrations.Embed)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = provider.Up(ctx); err != nil {
		log.Fatal(err)
	}

	// storage
	usersStorage := storage.NewUsers(db)

	// authenticator
	tokenManager := token.NewManager[model.TokenPayload]([]byte(cfg.JWT.Secret), cfg.JWT.TokenTTL)
	authenticator := service.NewAuthenticator(usersStorage, tokenManager)

	// grpc server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	authenticatorServer := ai.NewAuthenticatorServer(authenticator)
	api.RegisterAuthenticatorServer(server, authenticatorServer)

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
