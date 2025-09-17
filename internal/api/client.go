package api

import (
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
