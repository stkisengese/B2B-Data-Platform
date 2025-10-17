package workers

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// WorkerPool manages a pool of workers for concurrent job processing
type WorkerPool struct {
	workerCount  int                // Number of workers in the pool(controls concurrency)
	jobQueue     chan Job           // Channel for job queue(enforces queue size)
	workers      []*Worker          // Slice of worker instances
	wg           sync.WaitGroup     // WaitGroup to track active workers
	ctx          context.Context    // Context for managing worker lifecycle
	cancel       context.CancelFunc // Cancel function for context cancellation
	logger       *logrus.Logger     // Logger for logging pool activities
	metrics      *PoolMetrics       // Metrics for monitoring pool performance
	shutdownOnce sync.Once          // Ensures shutdown is only performed once
}

// PoolConfig holds configuration for the worker pool
type PoolConfig struct {
	WorkerCount int            // Number of workers in the pool
	QueueSize   int            // Maximum size of the job queue
	Logger      *logrus.Logger // Logger for logging pool activities
}

// PoolMetrics tracks worker pool performance
type PoolMetrics struct {
	JobsProcessed   int64
	JobsCompleted   int64
	JobsFailed      int64
	ActiveWorkers   int32
	QueueLength     int32
	AverageExecTime time.Duration
	LastActivity    time.Time
	mutex           sync.RWMutex
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(config PoolConfig) *WorkerPool {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	// Create cancellable context for worker lifecycle management
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workerCount: config.WorkerCount,
		jobQueue:    make(chan Job, config.QueueSize),
		workers:     make([]*Worker, config.WorkerCount),
		ctx:         ctx,
		cancel:      cancel,
		logger:      config.Logger,
		metrics: &PoolMetrics{
			LastActivity: time.Now(),
		},
	}

	return pool
}

// Start initializes and starts all workers
func (p *WorkerPool) Start() error {
	p.logger.WithField("workers", p.workerCount).Info("Starting worker pool")

	for i := 0; i < p.workerCount; i++ {
		worker := NewWorker(i, p.jobQueue, p.logger, p.metrics)
		p.workers[i] = worker

		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Start(p.ctx)
		}(worker)
	}

	p.logger.Info("Worker pool started successfully")
	return nil
}
