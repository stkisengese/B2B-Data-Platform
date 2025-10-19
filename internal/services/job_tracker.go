package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/stkisengese/B2B-Data-Platform/internal/workers"
)

// JobTracker manages job status and history
type JobTracker struct {
	jobs  map[string]*TrackedJob
	mutex sync.RWMutex
}

// TrackedJob wraps a job with additional tracking information
type TrackedJob struct {
	Job       workers.Job
	Status    workers.JobStatus
	Error     error
	UpdatedAt time.Time
}

// JobStatus represents the status of a tracked job
type JobStatus struct {
	ID          string            `json:"id"`
	Type        workers.JobType   `json:"type"`
	Status      workers.JobStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	RetryCount  int               `json:"retry_count"`
	Error       string            `json:"error,omitempty"`
	Duration    *time.Duration    `json:"duration,omitempty"`
}

// JobMetrics provides statistics about jobs
type JobMetrics struct {
	TotalJobs      int64     `json:"total_jobs"`
	PendingJobs    int64     `json:"pending_jobs"`
	ProcessingJobs int64     `json:"processing_jobs"`
	CompletedJobs  int64     `json:"completed_jobs"`
	FailedJobs     int64     `json:"failed_jobs"`
	RetryingJobs   int64     `json:"retrying_jobs"`
	LastJobAt      time.Time `json:"last_job_at"`
}

// CollectorMetrics combines pool and job metrics
type CollectorMetrics struct {
	PoolMetrics *workers.PoolMetrics `json:"pool_metrics"`
	JobMetrics  JobMetrics          `json:"job_metrics"`
}

// NewJobTracker creates a new job tracker
func NewJobTracker() *JobTracker {
	return &JobTracker{
		jobs: make(map[string]*TrackedJob),
	}
}

// TrackJob adds a job to the tracker
func (jt *JobTracker) TrackJob(job workers.Job) {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()

	jt.jobs[job.GetID()] = &TrackedJob{
		Job:       job,
		Status:    workers.StatusPending,
		UpdatedAt: time.Now(),
	}
}

// UpdateJob updates a job's status and error
func (jt *JobTracker) UpdateJob(jobID string, status workers.JobStatus, err error) {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()

	if tracked, exists := jt.jobs[jobID]; exists {
		tracked.Status = status
		tracked.Error = err
		tracked.UpdatedAt = time.Now()
	}
}
