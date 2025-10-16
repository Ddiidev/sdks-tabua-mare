package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://tabuamare.devtu.qzz.io/api/v1"
	defaultTimeout = 30 * time.Second
)

// Client é o cliente HTTP para a API Tide Table
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ClientOption é uma função que configura o Client
type ClientOption func(*Client)

// WithBaseURL configura uma URL base customizada
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithHTTPClient configura um http.Client customizado
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout configura o timeout das requisições
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient cria uma nova instância do cliente
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// doRequest executa uma requisição HTTP
func (c *Client) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &NetworkError{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimitExceeded
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = resp.StatusCode
			return nil, &apiErr
		}
		return nil, &APIError{
			Status:  resp.StatusCode,
			Code:    resp.StatusCode,
			Message: string(body),
		}
	}

	return body, nil
}
