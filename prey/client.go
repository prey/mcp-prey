package prey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"mcp-prey/internal"
)

type Client struct {
	BaseURL string
	Client  *http.Client
	APIKey  string
	Limiter *internal.MultiLimiter
}

func NewClient(cfg Config) *Client {
	transport := http.DefaultTransport
	client := &http.Client{Timeout: cfg.Timeout, Transport: transport}
	return &Client{
		BaseURL: cfg.URL,
		Client:  client,
		APIKey:  cfg.APIKey,
		Limiter: limiterFromConfig(cfg),
	}
}

func limiterFromConfig(cfg Config) *internal.MultiLimiter {
	if cfg.DisableRateLimit {
		return nil
	}
	return internal.NewPreyLimiter()
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("missing PREY_API_KEY")
	}
	if c.Limiter != nil {
		if err := c.Limiter.Wait(req.Context()); err != nil {
			return nil, err
		}
	}
	req.Header.Set("apikey", c.APIKey)
	return c.Client.Do(req)
}

func (c *Client) NewRequest(method, path string, q url.Values, body any) (*http.Request, error) {
	base := strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
	if len(q) > 0 {
		base += "?" + q.Encode()
	}
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, base, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) DoJSON(req *http.Request, out any) error {
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("prey api error: %s", resp.Status)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) DoRaw(req *http.Request) ([]byte, string, error) {
	resp, err := c.do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("prey api error: %s", resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return b, resp.Header.Get("Content-Type"), nil
}
