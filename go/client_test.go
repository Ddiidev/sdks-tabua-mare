package tabuamare

import (
	"context"
	"errors"
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

func TestWithAPIKey(t *testing.T) {
	client := NewClient(WithAPIKey("tm_live_test"))
	if client.apiKey != "tm_live_test" {
		t.Errorf("expected apiKey to be tm_live_test, got %s", client.apiKey)
	}
}

func TestDoRequest_AuthHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tm_live_test" {
			t.Errorf("expected Authorization header 'Bearer tm_live_test', got %s", got)
		}
		if got := r.Header.Get("X-Api-Key"); got != "tm_live_test" {
			t.Errorf("expected X-Api-Key header 'tm_live_test', got %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [], "total": 0}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithAPIKey("tm_live_test"))
	if _, err := client.doRequest(context.Background(), "GET", "/test"); err != nil {
		t.Fatalf("expected no error, got %v", err)
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
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.doRequest(context.Background(), "GET", "/test")
	if !errors.Is(err, ErrRateLimitExceeded) {
		t.Errorf("expected ErrRateLimitExceeded via errors.Is, got %v", err)
	}

	var rateErr *RateLimitError
	if !errors.As(err, &rateErr) {
		t.Fatalf("expected *RateLimitError, got %T", err)
	}
	if rateErr.RetryAfter != 60*time.Second {
		t.Errorf("expected RetryAfter 60s, got %s", rateErr.RetryAfter)
	}
}

func TestParseRetryAfter(t *testing.T) {
	testCases := []struct {
		value    string
		expected time.Duration
	}{
		{"", 0},
		{"60", 60 * time.Second},
		{" 3600 ", time.Hour},
		{"invalid", 0},
		{"-5", 0},
	}

	for _, tc := range testCases {
		if got := parseRetryAfter(tc.value); got != tc.expected {
			t.Errorf("parseRetryAfter(%q) = %s, expected %s", tc.value, got, tc.expected)
		}
	}
}

func TestDoRequest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": {"code": 404, "message": "not found"}}`))
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

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Message != "not found" {
		t.Errorf("expected message 'not found', got %s", apiErr.Message)
	}
}
