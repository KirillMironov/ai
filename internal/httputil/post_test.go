package httputil

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		statusCode     int
		payload        any
		ctx            context.Context
		wantStatusCode int
		wantBody       string
		wantErr        bool
	}{
		{
			name:           "success",
			url:            "/test",
			statusCode:     http.StatusOK,
			payload:        map[string]string{"key": "value"},
			ctx:            context.Background(),
			wantStatusCode: http.StatusOK,
			wantBody:       `{"key":"value"}` + "\n",
			wantErr:        false,
		},
		{
			name:           "empty payload",
			url:            "/test",
			statusCode:     http.StatusOK,
			payload:        nil,
			ctx:            context.Background(),
			wantStatusCode: http.StatusOK,
			wantBody:       "",
			wantErr:        false,
		},
		{
			name:           "invalid url",
			url:            "without-leading-slash", // invalid url
			statusCode:     http.StatusOK,
			payload:        nil,
			ctx:            context.Background(),
			wantStatusCode: http.StatusOK,
			wantBody:       "",
			wantErr:        true,
		},
		{
			name:           "invalid payload",
			url:            "/test",
			statusCode:     http.StatusOK,
			payload:        make(chan int), // invalid payload
			ctx:            context.Background(),
			wantStatusCode: http.StatusOK,
			wantBody:       "",
			wantErr:        true,
		},
		{
			name:           "unexpected status code",
			url:            "/test",
			statusCode:     http.StatusInternalServerError, // unexpected status code
			payload:        nil,
			ctx:            context.Background(),
			wantStatusCode: http.StatusOK,
			wantBody:       "",
			wantErr:        true,
		},
		{
			name:           "canceled context",
			url:            "/test",
			statusCode:     http.StatusOK,
			payload:        nil,
			ctx:            canceledContext(),
			wantStatusCode: http.StatusOK,
			wantBody:       "",
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got, want := r.Method, http.MethodPost; got != want {
					t.Errorf("request method = %v, want %v", got, want)
				}

				if got, want := r.URL.String(), tc.url; got != want {
					t.Errorf("request URL = %q, want %q", got, want)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}

				if got, want := string(body), tc.wantBody; got != want {
					t.Errorf("request body = %q, want %q", got, want)
				}

				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()

			body, err := post(tc.ctx, server.URL+tc.url, tc.wantStatusCode, tc.payload)
			if (err != nil) != tc.wantErr {
				t.Errorf("post() = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil {
				body.Close()
			}
		})
	}
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
