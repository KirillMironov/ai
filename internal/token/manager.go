package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	Data []byte `json:"data,omitempty"`
	jwt.RegisteredClaims
}

type Manager[T any] struct {
	secret   []byte
	tokenTTL time.Duration
}

func NewManager[T any](secret []byte, tokenTTL time.Duration) Manager[T] {
	return Manager[T]{
		secret:   secret,
		tokenTTL: tokenTTL,
	}
}

func (m Manager[T]) GenerateToken(payload T) (token string, err error) {
	if len(m.secret) == 0 {
		return "", errors.New("empty secret")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal token payload: %w", err)
	}

	now := time.Now()

	tokenClaims := claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString(m.secret)
}

func (m Manager[T]) ParseToken(token string) (payload T, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return payload, err
	}

	if !parsedToken.Valid {
		return payload, errors.New("invalid token")
	}

	tokenClaims, ok := parsedToken.Claims.(*claims)
	if !ok {
		return payload, fmt.Errorf("unexpected claims type: %T", tokenClaims)
	}

	if err = json.Unmarshal(tokenClaims.Data, &payload); err != nil {
		return payload, fmt.Errorf("unmarshal token payload: %w", err)
	}

	return payload, nil
}
