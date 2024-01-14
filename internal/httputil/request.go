package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func Post[T, U any](ctx context.Context, url string, req T) (resp U, err error) {
	body := new(bytes.Buffer)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return resp, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return resp, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return resp, err
	}
	defer httpResp.Body.Close()

	return resp, json.NewDecoder(httpResp.Body).Decode(&resp)
}
