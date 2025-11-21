package workers

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// MockJob for testing
type MockJob struct {
	BaseJob
	ExecuteFunc   func(ctx context.Context) error
	OnFailureFunc func(err error)
	shouldExecute int32
}

func (mj *MockJob) Execute(ctx context.Context) error {
	atomic.AddInt32(&mj.shouldExecute, 1)
	if mj.ExecuteFunc != nil {
		return mj.ExecuteFunc(ctx)
	}
	return nil
}

func (mj *MockJob) OnFailure(err error) {
	if mj.OnFailureFunc != nil {
		mj.OnFailureFunc(err)
	} else {
		mj.BaseJob.OnFailure(err)
	}
}

func TestWorkerPool_NewWorkerPool(t *testing.T) {
	config := PoolConfig{
		WorkerCount: 3,
		QueueSize:   10,
		Logger:      logrus.New(),
	}

	pool := NewWorkerPool(config)

	if pool.workerCount != 3 {
		t.Errorf("Expected worker count 3, got %d", pool.workerCount)
	}

	if cap(pool.jobQueue) != 10 {
		t.Errorf("Expected queue size 10, got %d", cap(pool.jobQueue))
	}
}

func TestWorkerPool_StartAndShutdown(t *testing.T) {
	config := PoolConfig{
		WorkerCount: 2,
		QueueSize:   5,
		Logger:      logrus.New(),
	}

	pool := NewWorkerPool(config)

	// Test start
	err := pool.Start()
	if err != nil {
		t.Fatalf("Failed to start pool: %v", err)
	}

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	metrics := pool.GetMetrics()
	if metrics.ActiveWorkers != 2 {
		t.Errorf("Expected 2 active workers, got %d", metrics.ActiveWorkers)
	}

	// Test shutdown
	err = pool.Shutdown(5 * time.Second)
	if err != nil {
		t.Fatalf("Failed to shutdown pool: %v", err)
	}

	// Verify shutdown
	finalMetrics := pool.GetMetrics()
	if finalMetrics.ActiveWorkers != 0 {
		t.Errorf("Expected 0 active workers after shutdown, got %d", finalMetrics.ActiveWorkers)
	}
}

func TestWorkerPool_SubmitJob(t *testing.T) {
	config := PoolConfig{
		WorkerCount: 1,
		QueueSize:   2,
		Logger:      logrus.New(),
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Shutdown(1 * time.Second)

	// Create a job that increments a counter
	var executionCount int32
	job := &MockJob{
		BaseJob: BaseJob{
			ID:   "test-job-1",
			Type: CollectionJob,
		},
		ExecuteFunc: func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)
			return nil
		},
	}

	// Submit job
	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Wait for job to complete
	time.Sleep(200 * time.Millisecond)

	if atomic.LoadInt32(&executionCount) != 1 {
		t.Errorf("Expected job to execute once, got %d executions", executionCount)
	}

	metrics := pool.GetMetrics()
	if metrics.JobsProcessed != 1 {
		t.Errorf("Expected 1 job processed, got %d", metrics.JobsProcessed)
	}
}

func TestWorkerPool_JobExecution(t *testing.T) {
	config := PoolConfig{
		WorkerCount: 2,
		QueueSize:   5,
		Logger:      logrus.New(),
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Shutdown(2 * time.Second)

	// Submit multiple jobs
	jobCount := 5
	var completedJobs int32

	for i := 0; i < jobCount; i++ {
		job := &MockJob{
			BaseJob: BaseJob{
				ID:   fmt.Sprintf("test-job-%d", i),
				Type: CollectionJob,
			},
			ExecuteFunc: func(ctx context.Context) error {
				time.Sleep(10 * time.Millisecond) // Simulate work
				atomic.AddInt32(&completedJobs, 1)
				return nil
			},
		}

		err := pool.Submit(job)
		if err != nil {
			t.Errorf("Failed to submit job %d: %v", i, err)
		}
	}

	// Wait for all jobs to complete
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("Timeout waiting for jobs to complete. Completed: %d", atomic.LoadInt32(&completedJobs))
		case <-ticker.C:
			if atomic.LoadInt32(&completedJobs) == int32(jobCount) {
				return // All jobs completed
			}
		}
	}
}

func TestWorkerPool_FailedJob(t *testing.T) {
	config := PoolConfig{
		WorkerCount: 1,
		QueueSize:   2,
		Logger:      logrus.New(),
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Shutdown(1 * time.Second)

	expectedError := errors.New("test error")
	failureSignal := make(chan error, 1)

	job := &MockJob{
		BaseJob: BaseJob{
			ID:   "failing-job",
			Type: CollectionJob,
		},
		ExecuteFunc: func(ctx context.Context) error {
			return expectedError
		},
	}

	job.OnFailureFunc = func(err error) {
		job.BaseJob.OnFailure(err)
		failureSignal <- err
	}

	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	select {
	case err := <-failureSignal:
		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for OnFailure callback")
	}

	metrics := pool.GetMetrics()
	if metrics.JobsFailed != 1 {
		t.Errorf("Expected 1 failed job, got %d", metrics.JobsFailed)
	}
}
