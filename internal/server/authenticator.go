package server

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/ai/api/ai"
)

type authenticatorService interface {
	SignUp(ctx context.Context, username, password string) (token string, err error)
	SignIn(ctx context.Context, username, password string) (token string, err error)
}

type Authenticator struct {
	service authenticatorService
	api.UnimplementedAuthenticatorServer
}

func NewAuthenticator(service authenticatorService) Authenticator {
	return Authenticator{service: service}
}

func (a Authenticator) SignUp(ctx context.Context, request *api.SignUpRequest) (*api.SignUpResponse, error) {
	token, err := a.service.SignUp(ctx, request.Username, request.Password)
	if err != nil {
		slog.Error("failed to call service.SignUp", err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &api.SignUpResponse{Token: token}, nil
}

func (a Authenticator) SignIn(ctx context.Context, request *api.SignInRequest) (*api.SignInResponse, error) {
	token, err := a.service.SignIn(ctx, request.Username, request.Password)
	if err != nil {
		slog.Error("failed to call service.SignIn", err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &api.SignInResponse{Token: token}, nil
}
