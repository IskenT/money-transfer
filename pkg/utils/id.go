package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// IDGenerator
type IDGenerator struct {
	mu      sync.Mutex
	counter uint64
}

// NewIDGenerator
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		counter: uint64(time.Now().UnixNano()),
	}
}

// GenerateTransferID
func (g *IDGenerator) GenerateTransferID() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.counter++
	return fmt.Sprintf("TRF%d", g.counter)
}

// GenerateTransactionID
func (g *IDGenerator) GenerateTransactionID() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.counter++
	return fmt.Sprintf("TRX%d", g.counter)
}

// GenerateUUID
func GenerateUUID() string {
	return uuid.New().String()
}
