package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Post[T any](ctx context.Context, url string, statusCode int, payload any) (resp T, err error) {
	body, err := post(ctx, url, statusCode, payload)
	if err != nil {
		return resp, err
	}
	defer body.Close()

	return resp, json.NewDecoder(body).Decode(&resp)
}

func PostBody(ctx context.Context, url string, statusCode int, payload any) (body io.ReadCloser, err error) {
	return post(ctx, url, statusCode, payload)
}

func post(ctx context.Context, url string, statusCode int, payload any) (body io.ReadCloser, err error) {
	buf := new(bytes.Buffer)
	if payload != nil {
		if err = json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, err
		}
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
