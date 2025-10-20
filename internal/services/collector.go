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
	sourceManager *api.SourceManager  // manages registered data sources
	workerPool    *workers.WorkerPool // worker pool for processing collection jobs
	storage       database.Storage    // interface type for decoupling storage implementation
	logger        *logrus.Logger      // logger for logging events and errors
	jobTracker    *JobTracker         // tracks job statuses and metrics
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
func NewCollectorService(sourceManager *api.SourceManager, storage database.Storage, config CollectorConfig) *CollectorService {
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

// ScheduleCollection creates and queues a collection job for execution
func (cs *CollectorService) ScheduleCollection(request CollectionRequest) (string, error) {
	// Create collection job and inject dependencies(source manager, storage, logger)
	job := &CollectionJobImpl{
		BaseJob: workers.BaseJob{
			ID:         generateJobID(),
			Type:       workers.CollectionJob,
			Status:     workers.StatusPending,
			CreatedAt:  time.Now(),
			MaxRetries: request.MaxRetries,
		},
		Request:       request,
		SourceManager: cs.sourceManager,
		Storage:       cs.storage,
		Logger:        cs.logger,
	}

	// Track job status in job tracker for monitoring and metrics purposes
	cs.jobTracker.TrackJob(job)

	// Submit job to worker pool for processing
	if err := cs.workerPool.Submit(job); err != nil {
		cs.jobTracker.UpdateJob(job.GetID(), workers.StatusFailed, err)
		return "", fmt.Errorf("failed to submit collection job: %w", err)
	}

	// Log job scheduling
	cs.logger.WithFields(logrus.Fields{
		"job_id": job.GetID(),
		"source": request.Source,
		"query":  request.Params.Query,
	}).Info("Collection job scheduled")

	return job.GetID(), nil
}

// GetJobStatus returns the current status of a job by its ID
func (cs *CollectorService) GetJobStatus(jobID string) (*JobStatus, error) {
	return cs.jobTracker.GetJobStatus(jobID)
}

// GetMetrics returns collector service metrics
func (cs *CollectorService) GetMetrics() *CollectorMetrics {
	poolMetrics := cs.workerPool.GetMetrics()
	jobMetrics := cs.jobTracker.GetMetrics()

	return &CollectorMetrics{
		PoolMetrics: poolMetrics,
		JobMetrics:  jobMetrics,
	}
}
