package workers

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// Worker processes jobs from the job queue
type Worker struct {
	id       int
	jobQueue <-chan Job
	logger   *logrus.Logger
	metrics  *PoolMetrics
}

// NewWorker creates a new worker
func NewWorker(id int, jobQueue <-chan Job, logger *logrus.Logger, metrics *PoolMetrics) *Worker {
	return &Worker{
		id:       id,
		jobQueue: jobQueue,
		logger:   logger,
		metrics:  metrics,
	}
}

// Start begins processing jobs
func (w *Worker) Start(ctx context.Context) {
	w.logger.WithField("worker_id", w.id).Info("Worker started")
	atomic.AddInt32(&w.metrics.ActiveWorkers, 1)

	defer func() {
		atomic.AddInt32(&w.metrics.ActiveWorkers, -1)
		w.logger.WithField("worker_id", w.id).Info("Worker stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-w.jobQueue:
			if !ok {
				return // Queue closed
			}
			w.processJob(ctx, job)
		}
	}
}
