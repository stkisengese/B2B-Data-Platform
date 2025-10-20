package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stkisengese/B2B-Data-Platform/internal/api"
	"github.com/stkisengese/B2B-Data-Platform/internal/database"
	"github.com/stkisengese/B2B-Data-Platform/internal/workers"
)

// CollectionRequest represents a request for data collection
type CollectionRequest struct {
	Source     string
	Params     api.CollectionParams
	MaxRetries int
	Timeout    time.Duration
}

// CollectionJobImpl implements the Job interface for data collection
type CollectionJobImpl struct {
	workers.BaseJob
	Request       CollectionRequest
	SourceManager *api.SourceManager
	Storage       database.Storage
	Logger        *logrus.Logger
}

// Execute performs the actual data collection
func (cj *CollectionJobImpl) Execute(ctx context.Context) error {
	now := time.Now()
	cj.StartedAt = &now
	cj.Status = workers.StatusProcessing

	cj.Logger.WithFields(logrus.Fields{
		"job_id": cj.ID,
		"source": cj.Request.Source,
		"query":  cj.Request.Params.Query,
	}).Info("Starting data collection")

	// Create context with timeout
	collectCtx := ctx
	if cj.Request.Timeout > 0 {
		var cancel context.CancelFunc
		collectCtx, cancel = context.WithTimeout(ctx, cj.Request.Timeout)
		defer cancel()
	}

	// Collect data from source
	records, err := cj.SourceManager.CollectFromSource(collectCtx, cj.Request.Source, cj.Request.Params)
	if err != nil {
		return fmt.Errorf("data collection failed: %w", err)
	}

	// Store collected records
	if err := cj.storeRecords(ctx, records); err != nil {
		return fmt.Errorf("failed to store records: %w", err)
	}

	cj.Logger.WithFields(logrus.Fields{
		"job_id":        cj.ID,
		"source":        cj.Request.Source,
		"records_count": len(records),
		"duration":      time.Since(now),
	}).Info("Data collection completed")

	return nil
}

func (cj *CollectionJobImpl) storeRecords(ctx context.Context, records []api.RawRecord) error {
	// For now, we'll store raw records in a simple format
	// This will be enhanced in later issues with proper data processing

	for _, record := range records {
		if err := cj.Storage.StoreRawRecord(ctx, record); err != nil {
			cj.Logger.WithFields(logrus.Fields{
				"job_id":    cj.ID,
				"record_id": record.ID,
				"error":     err,
			}).Error("Failed to store record")
			return err
		}
	}

	return nil
}

// ShouldRetry determines retry logic for collection jobs
func (cj *CollectionJobImpl) ShouldRetry(err error) bool {
	// Don't retry on context cancellation or API key issues
	if err == context.Canceled || err == api.ErrAPIKeyMissing {
		return false
	}

	return cj.BaseJob.ShouldRetry(err)
}
