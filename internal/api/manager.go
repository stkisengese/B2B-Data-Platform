package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SourceManager manages multiple data sources and their API clients
type SourceManager struct {
	sources map[string]DataSource
	logger  *logrus.Logger
	mutex   sync.RWMutex
}

// NewSourceManager creates a new source manager
func NewSourceManager() *SourceManager {
	return &SourceManager{
		sources: make(map[string]DataSource),
		logger:  logrus.New(),
	}
}

// RegisterSource adds a new data source to the manager
func (sm *SourceManager) RegisterSource(source DataSource) error {
	if err := source.Validate(); err != nil {
		return fmt.Errorf("source validation failed: %w", err)
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.sources[source.GetName()] = source
	sm.logger.WithField("source", source.GetName()).Info("Data source registered")

	return nil
}

// GetSource retrieves a data source by name
func (sm *SourceManager) GetSource(name string) (DataSource, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	source, exists := sm.sources[name]
	if !exists {
		return nil, fmt.Errorf("data source '%s' not found", name)
	}

	return source, nil
}

// ListSources returns all registered data source names
func (sm *SourceManager) ListSources() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var names []string
	for name := range sm.sources {
		names = append(names, name)
	}

	return names
}

// CollectFromSource collects data from a specific source
func (sm *SourceManager) CollectFromSource(ctx context.Context, sourceName string, params CollectionParams) ([]RawRecord, error) {
	source, err := sm.GetSource(sourceName)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	records, err := source.Collect(ctx, params)
	duration := time.Since(startTime)

	sm.logger.WithFields(logrus.Fields{
		"source":   sourceName,
		"duration": duration,
		"records":  len(records),
		"error":    err,
	}).Info("Data collection completed")

	return records, err
}

// CollectFromAllSources collects data from all registered sources concurrently
func (sm *SourceManager) CollectFromAllSources(ctx context.Context, params CollectionParams) (map[string][]RawRecord, map[string]error) {
	sm.mutex.RLock()
	sources := make(map[string]DataSource, len(sm.sources))
	for name, source := range sm.sources {
		sources[name] = source
	}
	sm.mutex.RUnlock()

	results := make(map[string][]RawRecord)
	errors := make(map[string]error)

	var wg sync.WaitGroup
	var resultMutex sync.Mutex

	for name, source := range sources {
		wg.Add(1)
		go func(sourceName string, src DataSource) {
			defer wg.Done()

			records, err := sm.CollectFromSource(ctx, sourceName, params)

			resultMutex.Lock()
			if err != nil {
				errors[sourceName] = err
			} else {
				results[sourceName] = records
			}
			resultMutex.Unlock()
		}(name, source)
	}

	wg.Wait()
	return results, errors
}
