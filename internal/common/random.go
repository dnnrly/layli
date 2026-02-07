package common

import (
	"math/rand"
	"os"
	"sync"
	"time"
)

// Service provides thread-safe random number generation with configurable seeding
type Service struct {
	rng *rand.Rand
	mu  sync.Mutex
}

// New creates a new random service with the given seed
func New(seed int64) *Service {
	return &Service{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// NewDefault creates a random service with time-based seed
func NewDefault() *Service {
	return &Service{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Shuffle shuffles n elements using the provided swap function
func (s *Service) Shuffle(n int, swap func(i, j int)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rng.Shuffle(n, swap)
}

// defaultService is used for backward compatibility
var defaultService = &Service{
	rng: rand.New(rand.NewSource(time.Now().UnixNano())),
}

// Shuffle provides a global shuffle function that respects test seeding
func Shuffle(n int, swap func(i, j int)) {
	// Check if we're in test mode and need deterministic behavior
	if seedStr := os.Getenv("LAYLI_TEST_SEED"); seedStr != "" {
		defaultService.rng = rand.New(rand.NewSource(42))
	}
	defaultService.Shuffle(n, swap)
}
