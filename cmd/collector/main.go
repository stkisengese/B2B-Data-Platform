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
			log.Printf("Failed to register Companies House source: %v", err)
		}
	}

	if cfg.DataSources.OpenCorporates.Enabled {
		ocSource := sources.NewOpenCorporatesSource(cfg.DataSources.OpenCorporates.APIKey)
		if err := sourceManager.RegisterSource(ocSource); err != nil {
			log.Printf("Failed to register OpenCorporates source: %v", err)
		}
	}

	// Example data collection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := api.CollectionParams{
		Query:  "technology",
		Limit:  10,
		Offset: 0,
	}

	log.Println("Starting data collection from all sources...")
	results, errors := sourceManager.CollectFromAllSources(ctx, params)

	for source, records := range results {
		log.Printf("Collected %d records from %s", len(records), source)
	}

	for source, err := range errors {
		log.Printf("Error collecting from %s: %v", source, err)
	}

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Collector service shutting down...")
}
