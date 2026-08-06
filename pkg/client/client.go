package panel

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type service struct {
	client *Client
}

type Client struct {
	url  string
	http *http.Client

	Image    *ImageService
	Instance *InstanceService
}

func NewClient(url string) *Client {
	client := &Client{
		url:  url,
		http: &http.Client{},
	}

	client.Image = &ImageService{client}
	client.Instance = &InstanceService{client}

	return client
}

func (c *Client) newRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.url, url), body)
}
