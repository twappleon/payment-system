package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OrderClient struct {
	baseURL string
	http    *http.Client
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{baseURL: baseURL, http: http.DefaultClient}
}

func (c *OrderClient) QueryPayment(ctx context.Context, platformOrderNo string, out interface{}) error {
	endpoint := c.baseURL + "/internal/v1/payment/query?platform_order_no=" + url.QueryEscape(platformOrderNo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
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

