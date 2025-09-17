package api

import (
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
