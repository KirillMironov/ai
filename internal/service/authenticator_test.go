package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/KirillMironov/ai/internal/mock"
	"github.com/KirillMironov/ai/internal/model"
)

const (
	testUserID   = "user_id"
	testUsername = "username"
	testPassword = "password"
	testToken    = "token"
)

var (
	errStorage      = errors.New("storage error")
	errTokenManager = errors.New("token manager error")
)

func TestAuthenticator_SignUp(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		password     string
		usersStorage usersStorage
		tokenManager tokenManager
		wantErr      bool
		wantToken    string
	}{
		{
			name:     "success",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, nil
				},
				SaveUserFunc: func(context.Context, model.User) error {
					return nil
				},
			},
			tokenManager: &mock.TokenManager{
				GenerateTokenFunc: func(model.TokenPayload) (string, error) {
					return testToken, nil
				},
				ParseTokenFunc: nil,
			},
			wantErr:   false,
			wantToken: testToken,
		},
		{
			name:         "empty username",
			username:     "",
			password:     testPassword,
			usersStorage: nil,
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:         "empty password",
			username:     testUsername,
			password:     "",
			usersStorage: nil,
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "user already exists",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, true, nil
				},
				SaveUserFunc: nil,
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "get user by username error",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, errStorage
				},
				SaveUserFunc: nil,
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "save user error",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, nil
				},
				SaveUserFunc: func(context.Context, model.User) error {
					return errStorage
				},
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "generate token error",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, nil
				},
				SaveUserFunc: func(context.Context, model.User) error {
					return nil
				},
			},
			tokenManager: &mock.TokenManager{
				GenerateTokenFunc: func(model.TokenPayload) (string, error) {
					return "", errTokenManager
				},
				ParseTokenFunc: nil,
			},
			wantErr:   true,
			wantToken: "",
		},
		{
			name:     "generate hash from password error",
			username: testUsername,
			password: string(make([]byte, 100)), // bcrypt: password too long (> 72 bytes)
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, nil
				},
				SaveUserFunc: nil,
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := NewAuthenticator(tc.usersStorage, tc.tokenManager)

			token, err := authenticator.SignUp(context.Background(), tc.username, tc.password)
			if (err != nil) != tc.wantErr {
				t.Fatalf("authenticator.SignUp() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := token, tc.wantToken; got != want {
				t.Errorf("token = %q, want %q", got, want)
			}
		})
	}
}

func TestAuthenticator_SignIn(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		password     string
		usersStorage usersStorage
		tokenManager tokenManager
		wantErr      bool
		wantToken    string
	}{
		{
			name:     "success",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{
						ID:             testUserID,
						Username:       testUsername,
						HashedPassword: hashFromPassword(t, testPassword),
						CreatedAt:      time.Now(),
					}, true, nil
				},
				SaveUserFunc: nil,
			},
			tokenManager: &mock.TokenManager{
				GenerateTokenFunc: func(model.TokenPayload) (string, error) {
					return testToken, nil
				},
				ParseTokenFunc: nil,
			},
			wantErr:   false,
			wantToken: testToken,
		},
		{
			name:         "empty username",
			username:     "",
			password:     testPassword,
			usersStorage: nil,
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:         "empty password",
			username:     testUsername,
			password:     "",
			usersStorage: nil,
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "invalid password",
			username: testUsername,
			password: "invalid-password",
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{
						ID:             testUserID,
						Username:       testUsername,
						HashedPassword: hashFromPassword(t, testPassword),
						CreatedAt:      time.Now(),
					}, true, nil
				},
				SaveUserFunc: nil,
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "user does not exist",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, nil
				},
				SaveUserFunc: nil,
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "get user by username error",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{}, false, errStorage
				},
				SaveUserFunc: nil,
			},
			tokenManager: nil,
			wantErr:      true,
			wantToken:    "",
		},
		{
			name:     "generate token error",
			username: testUsername,
			password: testPassword,
			usersStorage: &mock.UsersStorage{
				GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
					return model.User{
						ID:             testUserID,
						Username:       testUsername,
						HashedPassword: hashFromPassword(t, testPassword),
						CreatedAt:      time.Now(),
					}, true, nil
				},
				SaveUserFunc: func(context.Context, model.User) error {
					return nil
				},
			},
			tokenManager: &mock.TokenManager{
				GenerateTokenFunc: func(model.TokenPayload) (string, error) {
					return "", errTokenManager
				},
				ParseTokenFunc: nil,
			},
			wantErr:   true,
			wantToken: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := NewAuthenticator(tc.usersStorage, tc.tokenManager)

			token, err := authenticator.SignIn(context.Background(), tc.username, tc.password)
			if (err != nil) != tc.wantErr {
				t.Fatalf("authenticator.SignIn() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := token, tc.wantToken; got != want {
				t.Errorf("token = %q, want %q", got, want)
			}
		})
	}
}

func TestAuthenticator_Authenticate(t *testing.T) {
	tests := []struct {
		name             string
		token            string
		usersStorage     usersStorage
		tokenManager     tokenManager
		wantErr          bool
		wantTokenPayload model.TokenPayload
	}{
		{
			name:         "success",
			token:        testToken,
			usersStorage: nil,
			tokenManager: &mock.TokenManager{
				GenerateTokenFunc: nil,
				ParseTokenFunc: func(string) (model.TokenPayload, error) {
					return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
				},
			},
			wantErr:          false,
			wantTokenPayload: model.TokenPayload{UserID: testUserID, Username: testUsername},
		},
		{
			name:             "empty token",
			token:            "",
			usersStorage:     nil,
			tokenManager:     nil,
			wantErr:          true,
			wantTokenPayload: model.TokenPayload{},
		},
		{
			name:         "generate token error",
			token:        "",
			usersStorage: nil,
			tokenManager: &mock.TokenManager{
				GenerateTokenFunc: nil,
				ParseTokenFunc: func(string) (model.TokenPayload, error) {
					return model.TokenPayload{}, errors.New("token manager error")
				},
			},
			wantErr:          true,
			wantTokenPayload: model.TokenPayload{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := NewAuthenticator(tc.usersStorage, tc.tokenManager)

			tokenPayload, err := authenticator.Authenticate(tc.token)
			if (err != nil) != tc.wantErr {
				t.Fatalf("authenticator.Authenticate() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := tokenPayload, tc.wantTokenPayload; got != want {
				t.Errorf("tokenPayload = %q, want %q", got, want)
			}
		})
	}
}

func hashFromPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword() = %v", err)
	}
	return string(hash)
}
