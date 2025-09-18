package api

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_NewCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		Name:             "TestBreaker",
		MaxFailures:      3,
		ResetTimeout:     30 * time.Second,
		HalfOpenMaxCalls: 2,
	}

	cb := NewCircuitBreaker(config)

	if cb.Name != "TestBreaker" {
		t.Errorf("Expected name 'TestBreaker', got '%s'", cb.Name)
	}

	if cb.GetState() != Closed {
		t.Errorf("Expected initial state Closed, got %v", cb.GetState())
	}
}

func TestCircuitBreaker_ClosedToOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Name:             "Test",
		MaxFailures:      2,
		ResetTimeout:     30 * time.Second,
		HalfOpenMaxCalls: 1,
	})

	if !cb.CanExecute() {
		t.Error("Should be able to execute when closed")
	}

	// First failure
	err := cb.Execute(func() error {
		return errors.New("test error")
	})
	if err == nil {
		t.Error("Expected error from execute function")
	}
	if cb.GetState() != Closed {
		t.Error("Should still be closed after first failure")
	}

	// Second failure should open circuit
	cb.Execute(func() error {
		return errors.New("test error")
	})
	if cb.GetState() != Open {
		t.Error("Should be open after max failures")
	}

	if cb.CanExecute() {
		t.Error("Should not be able to execute when open")
	}
}

func TestCircuitBreaker_OpenToHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Name:             "Test",
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	})

	// Force circuit to open
	cb.Execute(func() error {
		return errors.New("test error")
	})

	if cb.GetState() != Open {
		t.Error("Circuit should be open")
	}

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	if !cb.CanExecute() {
		t.Error("Should be able to execute after reset timeout")
	}

	// This should transition to half-open
	cb.Execute(func() error {
		return nil // Success
	})

	if cb.GetState() != HalfOpen {
		t.Error("Should be in half-open state")
	}
}

func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Name:             "Test",
		MaxFailures:      1,
		ResetTimeout:     10 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	})

	// Force to open then half-open
	cb.Execute(func() error { return errors.New("error") })
	time.Sleep(15 * time.Millisecond)
	cb.Execute(func() error { return nil })

	if cb.GetState() != HalfOpen {
		t.Error("Should be half-open")
	}

	// Complete half-open calls successfully
	cb.Execute(func() error { return nil })

	if cb.GetState() != Closed {
		t.Error("Should be closed after successful half-open calls")
	}
}
