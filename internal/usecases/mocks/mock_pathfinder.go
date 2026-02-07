package mocks

import (
	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockPathfinder is a mock implementation of Pathfinder.
type MockPathfinder struct {
	mock.Mock
}

// FindPaths implements Pathfinder.FindPaths.
func (m *MockPathfinder) FindPaths(diagram *domain.Diagram) error {
	args := m.Called(diagram)
	return args.Error(0)
}
