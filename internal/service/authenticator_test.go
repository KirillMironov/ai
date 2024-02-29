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

var (
	errStorage      = errors.New("storage error")
	errTokenManager = errors.New("token manager error")
)

type authenticatorMocks struct {
	usersStorage usersStorage
	tokenManager tokenManager
}

func TestAuthenticator_SignUp(t *testing.T) {
	const (
		testPassword = "password"
		testToken    = "token"
		testUsername = "username"
	)

	tests := []struct {
		name      string
		username  string
		password  string
		wantErr   bool
		wantToken string
		mocks     authenticatorMocks
	}{
		{
			name:      "success",
			username:  testUsername,
			password:  testPassword,
			wantErr:   false,
			wantToken: testToken,
			mocks: authenticatorMocks{
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
				},
			},
		},
		{
			name:      "empty username",
			username:  "",
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: nil,
			},
		},
		{
			name:      "empty password",
			username:  testUsername,
			password:  "",
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: nil,
			},
		},
		{
			name:      "user already exists",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{}, true, nil
					},
				},
				tokenManager: nil,
			},
		},
		{
			name:      "get user by username error",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{}, false, errStorage
					},
				},
				tokenManager: nil,
			},
		},
		{
			name:      "save user error",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{}, false, nil
					},
					SaveUserFunc: func(context.Context, model.User) error {
						return errStorage
					},
				},
				tokenManager: nil,
			},
		},
		{
			name:      "generate token error",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
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
				},
			},
		},
		{
			name:      "generate hash from password error",
			username:  testUsername,
			password:  string(make([]byte, 100)), // bcrypt: password too long (> 72 bytes)
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{}, false, nil
					},
				},
				tokenManager: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := NewAuthenticator(tc.mocks.usersStorage, tc.mocks.tokenManager)

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
	const (
		testPassword = "password"
		testToken    = "token"
		testUserID   = "user_id"
		testUsername = "username"
	)

	tests := []struct {
		name      string
		username  string
		password  string
		wantErr   bool
		wantToken string
		mocks     authenticatorMocks
	}{
		{
			name:      "success",
			username:  testUsername,
			password:  testPassword,
			wantErr:   false,
			wantToken: testToken,
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{
							ID:             testUserID,
							Username:       testUsername,
							HashedPassword: hashFromPassword(t, testPassword),
							CreatedAt:      time.Now(),
						}, true, nil
					},
				},
				tokenManager: &mock.TokenManager{
					GenerateTokenFunc: func(model.TokenPayload) (string, error) {
						return testToken, nil
					},
				},
			},
		},
		{
			name:      "empty username",
			username:  "",
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: nil,
			},
		},
		{
			name:      "empty password",
			username:  testUsername,
			password:  "",
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: nil,
			},
		},
		{
			name:      "invalid password",
			username:  testUsername,
			password:  "invalid-password",
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{
							ID:             testUserID,
							Username:       testUsername,
							HashedPassword: hashFromPassword(t, testPassword),
							CreatedAt:      time.Now(),
						}, true, nil
					},
				},
				tokenManager: nil,
			},
		},
		{
			name:      "user does not exist",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{}, false, nil
					},
				},
				tokenManager: nil,
			},
		},
		{
			name:      "get user by username error",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
				usersStorage: &mock.UsersStorage{
					GetUserByUsernameFunc: func(context.Context, string) (model.User, bool, error) {
						return model.User{}, false, errStorage
					},
				},
				tokenManager: nil,
			},
		},
		{
			name:      "generate token error",
			username:  testUsername,
			password:  testPassword,
			wantErr:   true,
			wantToken: "",
			mocks: authenticatorMocks{
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
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := NewAuthenticator(tc.mocks.usersStorage, tc.mocks.tokenManager)

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
	const (
		testToken    = "token"
		testUserID   = "user_id"
		testUsername = "username"
	)

	tests := []struct {
		name             string
		token            string
		wantErr          bool
		wantTokenPayload model.TokenPayload
		mocks            authenticatorMocks
	}{
		{
			name:             "success",
			token:            testToken,
			wantErr:          false,
			wantTokenPayload: model.TokenPayload{UserID: testUserID, Username: testUsername},
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: &mock.TokenManager{
					ParseTokenFunc: func(string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
			},
		},
		{
			name:             "empty token",
			token:            "",
			wantErr:          true,
			wantTokenPayload: model.TokenPayload{},
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: nil,
			},
		},
		{
			name:             "generate token error",
			token:            "",
			wantErr:          true,
			wantTokenPayload: model.TokenPayload{},
			mocks: authenticatorMocks{
				usersStorage: nil,
				tokenManager: &mock.TokenManager{
					ParseTokenFunc: func(string) (model.TokenPayload, error) {
						return model.TokenPayload{}, errors.New("token manager error")
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := NewAuthenticator(tc.mocks.usersStorage, tc.mocks.tokenManager)

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
