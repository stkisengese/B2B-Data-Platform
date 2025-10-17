package services

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateJobID creates a unique job identifier
func generateJobID() string {
	timestamp := time.Now().Unix()
	random := rand.Intn(10000)
	return fmt.Sprintf("job_%d_%04d", timestamp, random)
}
