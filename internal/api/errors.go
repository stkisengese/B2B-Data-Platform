package api

import "errors"

var (
	ErrAPIKeyMissing = errors.New("API key is missing")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")
	ErrInvalidResponse = errors.New("invalid response from API")
	ErrUnauthorized = errors.New("unauthorized request")
)