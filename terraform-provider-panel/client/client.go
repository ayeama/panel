package client

import (
	"context"
	"io"
	"net/http"
)

type Client struct {
	Endpoint string
	Username string
	Password string

	Client *http.Client
}

func (c *Client) NewRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	url := c.Endpoint + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if c.Username != "" || c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
