// File: internal/api/circuit_breaker.go
package api

import (
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	Closed CircuitState = iota
	Open
	HalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	Name             string
	MaxFailures      int
	ResetTimeout     time.Duration
	HalfOpenMaxCalls int

	mutex         sync.RWMutex
	state         CircuitState
	failures      int
	lastFailTime  time.Time
	halfOpenCalls int
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Name             string
	MaxFailures      int
	ResetTimeout     time.Duration
	HalfOpenMaxCalls int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		Name:             config.Name,
		MaxFailures:      config.MaxFailures,
		ResetTimeout:     config.ResetTimeout,
		HalfOpenMaxCalls: config.HalfOpenMaxCalls,
		state:            Closed,
	}
}