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

// processJob executes a single job and updates metrics
func (w *Worker) processJob(ctx context.Context, job Job) {
	start := time.Now()

	w.logger.WithFields(logrus.Fields{
		"worker_id": w.id,
		"job_id":    job.GetID(),
		"job_type":  job.GetType(),
	}).Debug("Processing job")

	err := job.Execute(ctx)
	duration := time.Since(start)

	// Update metrics
	w.metrics.mutex.Lock()
	if err != nil {
		w.metrics.JobsFailed++
		w.logger.WithFields(logrus.Fields{
			"worker_id": w.id,
			"job_id":    job.GetID(),
			"duration":  duration,
			"error":     err,
		}).Error("Job failed")
	} else {
		w.metrics.JobsCompleted++
		w.logger.WithFields(logrus.Fields{
			"worker_id": w.id,
			"job_id":    job.GetID(),
			"duration":  duration,
		}).Debug("Job completed successfully")
	}

	// Update average execution time
	if w.metrics.JobsCompleted > 0 {
		w.metrics.AverageExecTime = time.Duration(
			(int64(w.metrics.AverageExecTime)*w.metrics.JobsCompleted + int64(duration)) / (w.metrics.JobsCompleted + 1),
		)
	} else {
		w.metrics.AverageExecTime = duration
	}
	w.metrics.mutex.Unlock()

	// Handle job callbacks
	if err != nil {
		job.OnFailure(err)
	} else {
		job.OnSuccess()
	}
}
