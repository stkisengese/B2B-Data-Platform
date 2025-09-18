package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestAPIClient_NewAPIClient(t *testing.T) {
	config := ClientConfig{
		APIName:          "TestAPI",
		BaseURL:          "https://api.example.com",
		APIKey:           "test-key",
		RateLimit:        rate.Limit(5),
		RateBurst:        10,
		Timeout:          30 * time.Second,
		MaxRetries:       3,
		CircuitThreshold: 5,
	}

	client := NewAPIClient(config)

	if client.APIName != "TestAPI" {
		t.Errorf("Expected APIName to be 'TestAPI', got '%s'", client.APIName)
	}

	if client.BaseURL != "https://api.example.com" {
		t.Errorf("Expected BaseURL to be 'https://api.example.com', got '%s'", client.BaseURL)
	}

	if client.APIKey != "test-key" {
		t.Errorf("Expected APIKey to be 'test-key', got '%s'", client.APIKey)
	}
}

func TestAPIClient_MakeRequest_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", auth)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	config := ClientConfig{
		APIName:          "TestAPI",
		BaseURL:          server.URL,
		APIKey:           "test-key",
		RateLimit:        rate.Limit(10),
		RateBurst:        20,
		Timeout:          5 * time.Second,
		CircuitThreshold: 3,
	}

	client := NewAPIClient(config)
	ctx := context.Background()

	resp, err := client.MakeRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if client.GetRequestCount() != 1 {
		t.Errorf("Expected request count 1, got %d", client.GetRequestCount())
	}
}

func TestAPIClient_MakeRequest_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := ClientConfig{
		APIName:          "TestAPI",
		BaseURL:          server.URL,
		RateLimit:        rate.Limit(1), // Very low rate limit
		RateBurst:        1,
		Timeout:          5 * time.Second,
		CircuitThreshold: 3,
	}

	client := NewAPIClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// First request should succeed
	_, err := client.MakeRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("First request should succeed, got error: %v", err)
	}

	// Second immediate request should be rate limited due to context timeout
	_, err = client.MakeRequest(ctx, "GET", "/test", nil)
	if err == nil {
		t.Fatal("Expected rate limit error, got nil")
	}
}

func TestAPIClient_MakeRequest_CircuitBreaker(t *testing.T) {
	failCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failCount++
		if failCount <= 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	config := ClientConfig{
		APIName:          "TestAPI",
		BaseURL:          server.URL,
		RateLimit:        rate.Limit(100),
		RateBurst:        100,
		Timeout:          5 * time.Second,
		CircuitThreshold: 3,
	}

	client := NewAPIClient(config)
	ctx := context.Background()

	// Make requests that should fail and trigger circuit breaker
	for i := 0; i < 3; i++ {
		_, err := client.MakeRequest(ctx, "GET", "/test", nil)
		if err == nil {
			t.Errorf("Request %d should have failed", i+1)
		}
	}

	// Circuit breaker should be open now
	_, err := client.MakeRequest(ctx, "GET", "/test", nil)
	if err == nil || err.Error() != "circuit breaker is open for TestAPI" {
		t.Errorf("Expected circuit breaker error, got: %v", err)
	}
}
