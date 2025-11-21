package workers

import (
	"context"
	// "fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// MockJob for testing
type MockJob struct {
	BaseJob
	ExecuteFunc   func(ctx context.Context) error
	shouldExecute int32
}

func (mj *MockJob) Execute(ctx context.Context) error {
	atomic.AddInt32(&mj.shouldExecute, 1)
	if mj.ExecuteFunc != nil {
		return mj.ExecuteFunc(ctx)
	}
	return nil
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
