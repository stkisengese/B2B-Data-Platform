package workers

import (
	"context"
	"time"
)

// JobType represents different types of collection jobs
type JobType string

const (
	CollectionJob JobType = "collection"
	ProcessingJob JobType = "processing"
	ExportJob     JobType = "export"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
	StatusRetrying   JobStatus = "retrying"
)

// Job interface defines the contract for all job types
type Job interface {
	GetID() string
	GetType() JobType
	Execute(ctx context.Context) error
	OnSuccess()
	OnFailure(error)
	GetRetryCount() int
	ShouldRetry(error) bool
}

// BaseJob provides common functionality for all job types
type BaseJob struct {
	ID          string
	Type        JobType
	Status      JobStatus
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	RetryCount  int
	MaxRetries  int
	Error       error
}

// GetID returns the job ID
func (bj *BaseJob) GetID() string {
	return bj.ID
}

// GetType returns the job type
func (bj *BaseJob) GetType() JobType {
	return bj.Type
}

// GetRetryCount returns the current retry count
func (bj *BaseJob) GetRetryCount() int {
	return bj.RetryCount
}

// ShouldRetry determines if a job should be retried based on error and retry count
func (bj *BaseJob) ShouldRetry(err error) bool {
	return bj.RetryCount < bj.MaxRetries
}

// OnSuccess is called when the job completes successfully
func (bj *BaseJob) OnSuccess() {
	now := time.Now()
	bj.Status = StatusCompleted
	bj.CompletedAt = &now
	bj.Error = nil
}

// OnFailure is called when the job fails
func (bj *BaseJob) OnFailure(err error) {
	bj.Error = err
	if bj.ShouldRetry(err) {
		bj.Status = StatusRetrying
		bj.RetryCount++
	} else {
		now := time.Now()
		bj.Status = StatusFailed
		bj.CompletedAt = &now
	}
}
