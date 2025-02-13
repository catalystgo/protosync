package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	log "github.com/catalystgo/logger/cli"
)

type Client struct {
	ctx    context.Context
	client *http.Client
}

func NewClient(ctx context.Context) *Client {
	return &Client{ctx: ctx, client: &http.Client{}}
}

// Post sends a POST request to the given URL with the given body.
func (c *Client) Post(url string, body []byte, opts ...func(*http.Request)) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	log.Debugf("POST: %s", url)

	return c.do(req.WithContext(c.ctx))
}

func (c *Client) Get(url string, opts ...func(*http.Request)) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	log.Debugf("GET: %s", url)

	return c.do(req.WithContext(c.ctx))
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
