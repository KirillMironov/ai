package token

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
)

const (
	tokenSecret = "secret"
	tokenTTL    = time.Hour
)

type Payload struct {
	ID       int
	Username string
}

func TestManager(t *testing.T) {
	manager := NewManager[Payload]([]byte(tokenSecret), tokenTTL)
	payload := Payload{
		ID:       10,
		Username: "123",
	}

	token, err := manager.GenerateToken(payload)
	if err != nil {
		t.Fatalf("GenerateToken(payload) = %v, want %v", err, nil)
	}

	parsedPayload, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken(token) = %v, want %v", err, nil)
	}

	if !cmp.Equal(parsedPayload, payload) {
		t.Errorf("got %v, want %v", parsedPayload, payload)
	}
}

func TestManager_GenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		payload Payload
		wantErr bool
	}{
		{
			name:   "valid token",
			secret: tokenSecret,
			payload: Payload{
				ID:       10,
				Username: "123",
			},
			wantErr: false,
		},
		{
			name:   "empty secret",
			secret: "",
			payload: Payload{
				ID:       10,
				Username: "123",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager := NewManager[Payload]([]byte(tc.secret), tokenTTL)

			token, err := manager.GenerateToken(tc.payload)
			if (err != nil) != tc.wantErr {
				t.Errorf("GenerateToken(payload) = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				payload := mustParseToken(t, token, tc.secret)

				if !cmp.Equal(payload, tc.payload) {
					t.Errorf("got %v, want %v", payload, tc.payload)
				}
			}
		})
	}
}

func TestManager_ParseToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		secret      string
		wantPayload Payload
		wantErr     bool
	}{
		{
			name:        "valid token",
			token:       mustGenerateToken(t, Payload{ID: 10, Username: "123"}, tokenSecret, tokenTTL),
			secret:      tokenSecret,
			wantPayload: Payload{ID: 10, Username: "123"},
			wantErr:     false,
		},
		{
			name:        "invalid token format",
			token:       "123",
			secret:      tokenSecret,
			wantPayload: Payload{},
			wantErr:     true,
		},
		{
			name:        "invalid token signature",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			secret:      tokenSecret,
			wantPayload: Payload{},
			wantErr:     true,
		},
		{
			name:        "expired token",
			token:       mustGenerateToken(t, Payload{ID: 10, Username: "123"}, tokenSecret, -time.Hour),
			secret:      tokenSecret,
			wantPayload: Payload{},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager := NewManager[Payload]([]byte(tc.secret), tokenTTL)

			payload, err := manager.ParseToken(tc.token)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseToken(token) = %v, wantErr %v", err, tc.wantErr)
			}

			if !cmp.Equal(payload, tc.wantPayload) {
				t.Errorf("got %v, want %v", payload, tc.wantPayload)
			}
		})
	}
}

func mustParseToken(t *testing.T, token, secret string) (payload Payload) {
	t.Helper()

	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Fatalf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatalf("invalid token")
	}

	parsedClaims, ok := parsedToken.Claims.(*claims)
	if !ok {
		t.Fatalf("unexpected claims type: %T", parsedToken.Claims)
	}

	if err = json.Unmarshal(parsedClaims.Data, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	return payload
}

func mustGenerateToken(t *testing.T, payload Payload, secret string, ttl time.Duration) string {
	t.Helper()

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal token payload: %v", err)
	}

	now := time.Now()

	tokenClaims := claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return token
}
