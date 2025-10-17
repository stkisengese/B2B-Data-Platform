package services

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stkisengese/B2B-Data-Platform/internal/api"
	"github.com/stkisengese/B2B-Data-Platform/internal/database"
	"github.com/stkisengese/B2B-Data-Platform/internal/workers"
)

// CollectorService orchestrates data collection from multiple sources
type CollectorService struct {
	sourceManager *api.SourceManager
	workerPool    *workers.WorkerPool
	storage       *database.Storage
	logger        *logrus.Logger
	jobTracker    *JobTracker
}

// CollectorConfig holds configuration for the collector service
type CollectorConfig struct {
	WorkerCount   int
	QueueSize     int
	RetryAttempts int
	RetryDelay    time.Duration
	Logger        *logrus.Logger
}

// NewCollectorService creates a new collector service
func NewCollectorService(sourceManager *api.SourceManager, storage *database.Storage, config CollectorConfig) *CollectorService {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	poolConfig := workers.PoolConfig{
		WorkerCount: config.WorkerCount,
		QueueSize:   config.QueueSize,
		Logger:      config.Logger,
	}

	return &CollectorService{
		sourceManager: sourceManager,
		workerPool:    workers.NewWorkerPool(poolConfig),
		storage:       storage,
		logger:        config.Logger,
		jobTracker:    NewJobTracker(),
	}
}

// Start initializes the collector service
func (cs *CollectorService) Start() error {
	cs.logger.Info("Starting collector service")
	return cs.workerPool.Start()
}

// Stop gracefully shuts down the collector service
func (cs *CollectorService) Stop(timeout time.Duration) error {
	cs.logger.Info("Stopping collector service")
	return cs.workerPool.Shutdown(timeout)
}

