package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CoreClient struct {
	baseURL string
	http    *http.Client
}

func NewCoreClient(baseURL string) *CoreClient {
	return &CoreClient{baseURL: baseURL, http: http.DefaultClient}
}

func (c *CoreClient) HandleChannelCallback(ctx context.Context, payload interface{}, out interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/v1/payment/callback", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("core service returned %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

