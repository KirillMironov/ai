package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/KirillMironov/ai/internal/model"
)

type (
	usersStorage interface {
		SaveUser(ctx context.Context, user model.User) error
		GetUserByUsername(ctx context.Context, username string) (user model.User, exists bool, err error)
	}

	tokenManager interface {
		GenerateToken(payload model.TokenPayload) (token string, err error)
		ParseToken(token string) (payload model.TokenPayload, err error)
	}
)

type Authenticator struct {
	usersStorage usersStorage
	tokenManager tokenManager
}

func NewAuthenticator(usersStorage usersStorage, tokenManager tokenManager) Authenticator {
	return Authenticator{
		usersStorage: usersStorage,
		tokenManager: tokenManager,
	}
}

func (a Authenticator) SignUp(ctx context.Context, username, password string) (token string, err error) {
	if err = validateUsernamePassword(username, password); err != nil {
		return "", err
	}

	_, exists, err := a.usersStorage.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("get user by username: %w", err)
	}
	if exists {
		return "", fmt.Errorf("user with username '%s' already exists", username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := model.User{
		ID:             uuid.NewString(),
		Username:       username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}

	if err = a.usersStorage.SaveUser(ctx, user); err != nil {
		return "", fmt.Errorf("save user: %w", err)
	}

	return a.generateToken(user)
}

func (a Authenticator) SignIn(ctx context.Context, username, password string) (token string, err error) {
	if err = validateUsernamePassword(username, password); err != nil {
		return "", err
	}

	user, exists, err := a.usersStorage.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("get user by username: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("user with username '%s' does not exist", username)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("compare hash and password: %w", err)
	}

	return a.generateToken(user)
}

func (a Authenticator) Authenticate(token string) (model.TokenPayload, error) {
	if token == "" {
		return model.TokenPayload{}, errors.New("empty token")
	}

	return a.tokenManager.ParseToken(token)
}

func (a Authenticator) generateToken(user model.User) (token string, err error) {
	tokenPayload := model.TokenPayload{
		UserID:   user.ID,
		Username: user.Username,
	}

	token, err = a.tokenManager.GenerateToken(tokenPayload)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func validateUsernamePassword(username, password string) error {
	switch {
	case username == "":
		return errors.New("empty username")
	case password == "":
		return errors.New("empty password")
	default:
		return nil
	}
}
