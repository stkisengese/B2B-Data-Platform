package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// APIClient provides a configurable HTTP client with rate limiting and circuit breaking
type APIClient struct {
	HTTPClient     *http.Client
	RateLimiter    *rate.Limiter
	CircuitBreaker *CircuitBreaker
	Logger         *logrus.Logger
	APIName        string
	BaseURL        string
	APIKey         string
	RequestCount   int64
	mutex          sync.RWMutex
}

// ClientConfig holds configuration for API clients
type ClientConfig struct {
	APIName          string
	BaseURL          string
	APIKey           string
	RateLimit        rate.Limit
	RateBurst        int
	Timeout          time.Duration
	MaxRetries       int
	CircuitThreshold int
}

// NewAPIClient creates a new API client with rate limiting and circuit breaker
func NewAPIClient(config ClientConfig) *APIClient {
	return &APIClient{
		HTTPClient: &http.Client{
			Timeout: config.Timeout,
		},
		RateLimiter: rate.NewLimiter(config.RateLimit, config.RateBurst),
		CircuitBreaker: NewCircuitBreaker(CircuitBreakerConfig{
			Name:             config.APIName,
			MaxFailures:      config.CircuitThreshold,
			ResetTimeout:     30 * time.Second,
			HalfOpenMaxCalls: 3,
		}),
		Logger:  logrus.New(),
		APIName: config.APIName,
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
	}
}

// MakeRequest makes an HTTP request with rate limiting, retries, and circuit breaking
func (c *APIClient) MakeRequest(ctx context.Context, method, endpoint string, headers map[string]string) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Check circuit breaker
	if !c.CircuitBreaker.CanExecute() {
		return nil, fmt.Errorf("circuit breaker is open for %s", c.APIName)
	}

	url := c.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if configured
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute request with circuit breaker
	var resp *http.Response
	err = c.CircuitBreaker.Execute(func() error {
		c.mutex.Lock()
		c.RequestCount++
		c.mutex.Unlock()

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			c.Logger.WithFields(logrus.Fields{
				"api":      c.APIName,
				"method":   method,
				"endpoint": endpoint,
				"error":    err,
			}).Error("HTTP request failed")
			return err
		}

		if resp.StatusCode >= 400 {
			c.Logger.WithFields(logrus.Fields{
				"api":         c.APIName,
				"method":      method,
				"endpoint":    endpoint,
				"status_code": resp.StatusCode,
			}).Warn("HTTP request returned error status")
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}

		c.Logger.WithFields(logrus.Fields{
			"api":         c.APIName,
			"method":      method,
			"endpoint":    endpoint,
			"status_code": resp.StatusCode,
		}).Debug("HTTP request successful")

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetRequestCount returns the total number of requests made
func (c *APIClient) GetRequestCount() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.RequestCount
}
