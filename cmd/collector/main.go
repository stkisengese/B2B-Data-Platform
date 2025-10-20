package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stkisengese/B2B-Data-Platform/internal/api"
	"github.com/stkisengese/B2B-Data-Platform/internal/api/sources"
	"github.com/stkisengese/B2B-Data-Platform/internal/config"
	"github.com/stkisengese/B2B-Data-Platform/internal/database"
	"github.com/stkisengese/B2B-Data-Platform/internal/services"
)

func main() {
	// Setup logging
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	if os.Getenv("LOG_LEVEL") == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewDatabaseConnection(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create storage layer instance to abstract database operations
	storage := database.NewStorage(db)

	// Initialize source manager
	sourceManager := api.NewSourceManager()

	// Register data sources based on configuration
	if cfg.DataSources.CompaniesHouse.Enabled {
		chSource := sources.NewCompaniesHouseSource(cfg.DataSources.CompaniesHouse.APIKey)
		if err := sourceManager.RegisterSource(chSource); err != nil {
			logger.WithError(err).Error("Failed to register Companies House source")
		} else {
			logger.Info("Companies House source registered")
		}
	}

	if cfg.DataSources.OpenCorporates.Enabled {
		ocSource := sources.NewOpenCorporatesSource(cfg.DataSources.OpenCorporates.APIKey)
		if err := sourceManager.RegisterSource(ocSource); err != nil {
			logger.WithError(err).Error("Failed to register OpenCorporates source")
		} else {
			logger.Info("OpenCorporates source registered")
		}
	}

	// Initialize collector service
	collectorConfig := services.CollectorConfig{
		WorkerCount:   5,   // Configurable worker count
		QueueSize:     100, // Configurable queue size
		RetryAttempts: 3,
		RetryDelay:    5 * time.Second,
		Logger:        logger,
	}

	collector := services.NewCollectorService(sourceManager, storage, collectorConfig)

	// Start collector service
	if err := collector.Start(); err != nil {
		log.Fatalf("Failed to start collector service: %v", err)
	}

	logger.Info("Collector service started successfully")

	// Example: Schedule a collection job
	if len(sourceManager.ListSources()) > 0 {
		request := services.CollectionRequest{
			Source: sourceManager.ListSources()[0], // Use first available source
			Params: api.CollectionParams{
				Query:  "technology",
				Limit:  50,
				Offset: 0,
			},
			MaxRetries: 3,
			Timeout:    30 * time.Second,
		}

		jobID, err := collector.ScheduleCollection(request)
		if err != nil {
			logger.WithError(err).Error("Failed to schedule collection job")
		} else {
			logger.WithField("job_id", jobID).Info("Collection job scheduled")

			// Monitor job progress
			go func() {
				ticker := time.NewTicker(5 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					status, err := collector.GetJobStatus(jobID)
					if err != nil {
						logger.WithError(err).Error("Failed to get job status")
						return
					}

					logger.WithFields(logrus.Fields{
						"job_id":      jobID,
						"status":      status.Status,
						"retry_count": status.RetryCount,
					}).Info("Job status update")

					if status.Status == "completed" || status.Status == "failed" {
						return
					}
				}
			}()
		}
	}

	// Metrics reporting
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metrics := collector.GetMetrics()
			logger.WithFields(logrus.Fields{
				"active_workers":    metrics.PoolMetrics.ActiveWorkers,
				"jobs_processed":    metrics.PoolMetrics.JobsProcessed,
				"jobs_completed":    metrics.PoolMetrics.JobsCompleted,
				"jobs_failed":       metrics.PoolMetrics.JobsFailed,
				"average_exec_time": metrics.PoolMetrics.AverageExecTime,
				"queue_length":      metrics.PoolMetrics.QueueLength,
			}).Info("Collector metrics")
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Collector service shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := collector.Stop(30 * time.Second); err != nil {
		logger.WithError(err).Error("Error during collector shutdown")
	}

	select {
	case <-shutdownCtx.Done():
		logger.Warn("Shutdown timeout exceeded")
	default:
		logger.Info("Collector service shutdown completed")
	}
}
