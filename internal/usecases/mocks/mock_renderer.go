package mocks

import (
	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockRenderer is a mock implementation of Renderer.
type MockRenderer struct {
	mock.Mock
}

// Render implements Renderer.Render.
func (m *MockRenderer) Render(diagram *domain.Diagram, outputPath string) error {
	args := m.Called(diagram, outputPath)
	return args.Error(0)
}
