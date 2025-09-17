package api

import (
	"fmt"
	"sync"

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
