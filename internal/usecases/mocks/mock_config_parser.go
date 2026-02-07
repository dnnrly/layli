package mocks

import (
	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockConfigParser is a mock implementation of ConfigParser.
type MockConfigParser struct {
	mock.Mock
}

// Parse implements ConfigParser.Parse.
func (m *MockConfigParser) Parse(path string) (*domain.Diagram, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagram), args.Error(1)
}
