package tabuamare

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL to be %s, got %s", defaultBaseURL, client.baseURL)
	}
	if client.httpClient.Timeout != defaultTimeout {
		t.Errorf("expected timeout to be %s, got %s", defaultTimeout, client.httpClient.Timeout)
	}
}

func TestWithBaseURL(t *testing.T) {
	customURL := "https://custom.api.com"
	client := NewClient(WithBaseURL(customURL))
	if client.baseURL != customURL {
		t.Errorf("expected baseURL to be %s, got %s", customURL, client.baseURL)
	}
}

func TestWithTimeout(t *testing.T) {
	customTimeout := 60 * time.Second
	client := NewClient(WithTimeout(customTimeout))
	if client.httpClient.Timeout != customTimeout {
		t.Errorf("expected timeout to be %s, got %s", customTimeout, client.httpClient.Timeout)
	}
}

func TestDoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [], "total": 0}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	body, err := client.doRequest(context.Background(), "GET", "/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(body) == 0 {
		t.Fatal("expected non-empty body")
	}
}

func TestDoRequest_RateLimitExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.doRequest(context.Background(), "GET", "/test")
	if err != ErrRateLimitExceeded {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}
}

func TestDoRequest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": {"code": 404, "msg": "not found"}}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.doRequest(context.Background(), "GET", "/test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsAPIError(err) {
		t.Errorf("expected APIError, got %T", err)
	}
}
