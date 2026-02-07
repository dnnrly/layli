package mocks

import (
	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockLayoutEngine is a mock implementation of LayoutEngine.
type MockLayoutEngine struct {
	mock.Mock
}

// Arrange implements LayoutEngine.Arrange.
func (m *MockLayoutEngine) Arrange(diagram *domain.Diagram) error {
	args := m.Called(diagram)
	return args.Error(0)
}
