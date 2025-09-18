package api

import (
	"context"
	"testing"

	"golang.org/x/time/rate"
)

// MockDataSource for testing
type MockDataSource struct {
	name      string
	rateLimit rate.Limit
	shouldErr bool
	records   []RawRecord
}

func (m *MockDataSource) Collect(ctx context.Context, params CollectionParams) ([]RawRecord, error) {
	if m.shouldErr {
		return nil, ErrAPIKeyMissing
	}
	return m.records, nil
}

func (m *MockDataSource) GetName() string {
	return m.name
}

func (m *MockDataSource) GetRateLimit() rate.Limit {
	return m.rateLimit
}

func (m *MockDataSource) Validate() error {
	if m.shouldErr {
		return ErrAPIKeyMissing
	}
	return nil
}

func TestSourceManager_RegisterSource(t *testing.T) {
	manager := NewSourceManager()

	// Test successful registration
	mockSource := &MockDataSource{
		name:      "test-source",
		rateLimit: rate.Limit(5),
		shouldErr: false,
	}

	err := manager.RegisterSource(mockSource)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test registration with validation error
	invalidSource := &MockDataSource{
		name:      "invalid-source",
		shouldErr: true,
	}

	err = manager.RegisterSource(invalidSource)
	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestSourceManager_GetSource(t *testing.T) {
	manager := NewSourceManager()
	mockSource := &MockDataSource{
		name:      "test-source",
		rateLimit: rate.Limit(5),
	}

	manager.RegisterSource(mockSource)

	// Test getting existing source
	source, err := manager.GetSource("test-source")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if source.GetName() != "test-source" {
		t.Errorf("Expected source name 'test-source', got '%s'", source.GetName())
	}

	// Test getting non-existent source
	_, err = manager.GetSource("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent source, got nil")
	}
}

func TestSourceManager_CollectFromAllSources(t *testing.T) {
	manager := NewSourceManager()

	// Register successful source
	successSource := &MockDataSource{
		name: "success-source",
		records: []RawRecord{
			{ID: "1", Source: "success-source"},
			{ID: "2", Source: "success-source"},
		},
	}
	manager.RegisterSource(successSource)

	// Register failing source
	failSource := &MockDataSource{
		name:      "fail-source",
		shouldErr: true,
	}
	// Skip validation for test
	manager.sources["fail-source"] = failSource

	params := CollectionParams{
		Query: "test",
		Limit: 10,
	}

	results, errors := manager.CollectFromAllSources(context.Background(), params)

	// Check successful collection
	if len(results["success-source"]) != 2 {
		t.Errorf("Expected 2 records from success-source, got %d", len(results["success-source"]))
	}

	// Check failed collection
	if errors["fail-source"] == nil {
		t.Error("Expected error from fail-source, got nil")
	}
}
