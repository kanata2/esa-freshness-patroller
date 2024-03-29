package esa

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	v1BaseURL          = "https://api.esa.io/v1"
	MaxElementsPerPage = 100
)

type Client struct {
	baseURL    string
	apiKey     string
	team       string
	httpClient *http.Client
	debug      bool
}

type clientOption func(*Client)

func WithHTTPClient(hc *http.Client) clientOption {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func WithDebug() clientOption {
	return func(c *Client) {
		c.debug = true
	}
}

func NewClient(team, apiKey string, opts ...clientOption) *Client {
	c := &Client{
		baseURL:    v1BaseURL,
		apiKey:     apiKey,
		team:       team,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) newRequest(ctx context.Context, method, rpath string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+"/"+rpath, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", fmt.Sprintf("esa-freshness-patroller/v%s", version))
	if method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
