package api

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// CollectionParams holds parameters for data collection
type CollectionParams struct {
	Query    string
	Location string
	Limit    int
	Offset   int
	Filters  map[string]interface{}
}

// RawRecord represents a raw data record from an API
type RawRecord struct {
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	Data        map[string]interface{} `json:"data"`
	CollectedAt time.Time              `json:"collected_at"`
}

// DataSource interface defines the contract for data collection sources
type DataSource interface {
	Collect(ctx context.Context, params CollectionParams) ([]RawRecord, error)
	GetName() string
	GetRateLimit() rate.Limit
	Validate() error
}

// BaseDataSource provides common functionality for data sources
type BaseDataSource struct {
	Name      string
	APIClient *APIClient
	RateLimit rate.Limit
}

// GetName returns the data source name
func (ds *BaseDataSource) GetName() string {
	return ds.Name
}

// GetRateLimit returns the rate limit for this data source
func (ds *BaseDataSource) GetRateLimit() rate.Limit {
	return ds.RateLimit
}