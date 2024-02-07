package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Post[T, U any](ctx context.Context, url string, statusCode int, payload T) (resp U, err error) {
	body, err := post(ctx, url, statusCode, payload)
	if err != nil {
		return resp, err
	}
	defer body.Close()

	return resp, json.NewDecoder(body).Decode(&resp)
}

func PostBody[T any](ctx context.Context, url string, statusCode int, payload T) (body io.ReadCloser, err error) {
	return post(ctx, url, statusCode, payload)
}

func post[T any](ctx context.Context, url string, statusCode int, payload T) (body io.ReadCloser, err error) {
	buf := new(bytes.Buffer)
	if err = json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != statusCode {
		httpResp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	return httpResp.Body, nil
}
