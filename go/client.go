package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://tabuamare.api.br/api/v2"
	defaultTimeout = 30 * time.Second
)

// Client é o cliente HTTP para a API Tábua de Marés (v2)
type Client struct {
	baseURL    string
	apiKey     string
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

// WithAPIKey configura a api_key usada na autenticação (header Authorization: Bearer)
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
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

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("X-Api-Key", c.apiKey)
	}

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
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, &RateLimitError{RetryAfter: retryAfter}
	}

	if resp.StatusCode >= 400 {
		var envelope struct {
			Error *APIError `json:"error"`
		}
		if err := json.Unmarshal(body, &envelope); err == nil && envelope.Error != nil && envelope.Error.Message != "" {
			apiErr := envelope.Error
			apiErr.Status = resp.StatusCode
			return nil, apiErr
		}
		return nil, &APIError{
			Status:  resp.StatusCode,
			Code:    resp.StatusCode,
			Message: string(body),
		}
	}

	return body, nil
}

// parseRetryAfter converte o header Retry-After (segundos) em time.Duration
func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}

	seconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || seconds < 0 {
		return 0
	}

	return time.Duration(seconds) * time.Second
}
