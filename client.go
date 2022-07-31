package render

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	defaultAddress    = "https://api.render.com"
	defaultAPIVersion = "v1"

	envAddress = "RENDER_ADDR"
	envAPIKey  = "RENDER_API_KEY"
)

type Config struct {
	Address    string
	APIKey     string
	APIVersion string
}

type Client struct {
	cfg *Config
}

func DefaultConfig() *Config {
	cfg := &Config{
		Address:    defaultAddress,
		APIKey:     os.Getenv(envAPIKey),
		APIVersion: defaultAPIVersion,
	}

	if addr := os.Getenv(envAddress); addr != "" {
		cfg.Address = addr
	}

	return cfg
}

func NewClient(cfg *Config) (*Client, error) {
	return &Client{cfg: cfg}, nil
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	path = strings.Trim(path, "/")

	url := fmt.Sprintf("%s/%s/%s", c.cfg.Address, c.cfg.APIVersion, path)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	return req, nil
}

func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		err := ErrorFromResponse(resp)
		resp.Body.Close()
		return nil, err
	}

	return resp, err
}

func (c *Client) Services() *Services {
	return NewServices(c)
}
