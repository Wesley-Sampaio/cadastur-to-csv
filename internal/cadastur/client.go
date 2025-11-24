package cadastur

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
	"golang.org/x/net/html/charset"
)

// Client handles HTTP requests to the Cadastur API.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Client with a 30-second timeout.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs a GET request to the specified URL with context support.
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Respect charset declared in Content-Type and convert to UTF-8 when needed.
	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		// If unable to create reader for the declared charset, fall back to raw body.
		b, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return nil, err2
		}
		return b, nil
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Post performs a POST request to the specified URL with JSON payload and context support.
func (c *Client) Post(ctx context.Context, url string, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json;charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Respect charset declared in Content-Type and convert to UTF-8 when needed.
	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		b, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return nil, err2
		}
		return b, nil
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return body, nil
}
