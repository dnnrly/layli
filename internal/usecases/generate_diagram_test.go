package usecases

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/internal/usecases/mocks"
)

func TestGenerateDiagram_Execute_Success(t *testing.T) {
	// Arrange
	diagram := &domain.Diagram{
		Nodes: []domain.Node{
			{ID: "a", Width: 5, Height: 5},
			{ID: "b", Width: 5, Height: 5},
		},
		Edges: []domain.Edge{
			{ID: "e1", From: "a", To: "b"},
		},
		Config: domain.DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			Margin:         1,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}

	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	mockParser.On("Parse", "test.layli").Return(diagram, nil)
	mockLayout.On("Arrange", diagram).Return(nil)
	mockPathfinder.On("FindPaths", diagram).Return(nil)
	mockRenderer.On("Render", diagram, "output.svg").Return(nil)

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("test.layli", "output.svg")

	// Assert
	assert.NoError(t, err)
	mockParser.AssertExpectations(t)
	mockLayout.AssertExpectations(t)
	mockPathfinder.AssertExpectations(t)
	mockRenderer.AssertExpectations(t)
}

func TestGenerateDiagram_Execute_ParseError(t *testing.T) {
	// Arrange
	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	mockParser.On("Parse", "bad.layli").Return(nil, errors.New("syntax error"))

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("bad.layli", "output.svg")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse config")
	mockParser.AssertExpectations(t)
	mockLayout.AssertNotCalled(t, "Arrange", mock.Anything)
	mockPathfinder.AssertNotCalled(t, "FindPaths", mock.Anything)
	mockRenderer.AssertNotCalled(t, "Render", mock.Anything, mock.Anything)
}

func TestGenerateDiagram_Execute_ValidationError(t *testing.T) {
	// Arrange
	diagram := &domain.Diagram{
		Nodes:  []domain.Node{}, // Invalid: no nodes
		Edges:  []domain.Edge{},
		Config: domain.DiagramConfig{},
	}

	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	mockParser.On("Parse", "test.layli").Return(diagram, nil)

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("test.layli", "output.svg")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validate diagram")
	mockParser.AssertExpectations(t)
	mockLayout.AssertNotCalled(t, "Arrange", mock.Anything)
	mockPathfinder.AssertNotCalled(t, "FindPaths", mock.Anything)
	mockRenderer.AssertNotCalled(t, "Render", mock.Anything, mock.Anything)
}

func TestGenerateDiagram_Execute_LayoutError(t *testing.T) {
	// Arrange
	diagram := &domain.Diagram{
		Nodes: []domain.Node{
			{ID: "a", Width: 5, Height: 5},
		},
		Edges: []domain.Edge{},
		Config: domain.DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			Margin:         1,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}

	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	mockParser.On("Parse", "test.layli").Return(diagram, nil)
	mockLayout.On("Arrange", diagram).Return(errors.New("layout failed"))

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("test.layli", "output.svg")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "arrange layout")
	mockParser.AssertExpectations(t)
	mockLayout.AssertExpectations(t)
	mockPathfinder.AssertNotCalled(t, "FindPaths", mock.Anything)
	mockRenderer.AssertNotCalled(t, "Render", mock.Anything, mock.Anything)
}

func TestGenerateDiagram_Execute_PathfindingError(t *testing.T) {
	// Arrange
	diagram := &domain.Diagram{
		Nodes: []domain.Node{
			{ID: "a", Width: 5, Height: 5},
			{ID: "b", Width: 5, Height: 5},
		},
		Edges: []domain.Edge{
			{ID: "e1", From: "a", To: "b"},
		},
		Config: domain.DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			Margin:         1,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}

	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	mockParser.On("Parse", "test.layli").Return(diagram, nil)
	mockLayout.On("Arrange", diagram).Return(nil)
	mockPathfinder.On("FindPaths", diagram).Return(errors.New("no path found"))

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("test.layli", "output.svg")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find paths")
	mockParser.AssertExpectations(t)
	mockLayout.AssertExpectations(t)
	mockPathfinder.AssertExpectations(t)
	mockRenderer.AssertNotCalled(t, "Render", mock.Anything, mock.Anything)
}

func TestGenerateDiagram_Execute_RenderError(t *testing.T) {
	// Arrange
	diagram := &domain.Diagram{
		Nodes: []domain.Node{
			{ID: "a", Width: 5, Height: 5},
			{ID: "b", Width: 5, Height: 5},
		},
		Edges: []domain.Edge{
			{ID: "e1", From: "a", To: "b"},
		},
		Config: domain.DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			Margin:         1,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}

	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	mockParser.On("Parse", "test.layli").Return(diagram, nil)
	mockLayout.On("Arrange", diagram).Return(nil)
	mockPathfinder.On("FindPaths", diagram).Return(nil)
	mockRenderer.On("Render", diagram, "output.svg").Return(errors.New("cannot write file"))

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("test.layli", "output.svg")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "render diagram")
	mockParser.AssertExpectations(t)
	mockLayout.AssertExpectations(t)
	mockPathfinder.AssertExpectations(t)
	mockRenderer.AssertExpectations(t)
}

func TestGenerateDiagram_NewGenerateDiagram(t *testing.T) {
	// Arrange
	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	// Act
	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Assert
	assert.NotNil(t, uc)
	assert.Equal(t, mockParser, uc.configParser)
	assert.Equal(t, mockLayout, uc.layoutEngine)
	assert.Equal(t, mockPathfinder, uc.pathfinder)
	assert.Equal(t, mockRenderer, uc.renderer)
}

func TestGenerateDiagram_CallOrder(t *testing.T) {
	// This test ensures dependencies are called in the correct order
	diagram := &domain.Diagram{
		Nodes: []domain.Node{
			{ID: "a", Width: 5, Height: 5},
		},
		Edges: []domain.Edge{},
		Config: domain.DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			Margin:         1,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}

	mockParser := new(mocks.MockConfigParser)
	mockLayout := new(mocks.MockLayoutEngine)
	mockPathfinder := new(mocks.MockPathfinder)
	mockRenderer := new(mocks.MockRenderer)

	callOrder := []string{}

	mockParser.On("Parse", "test.layli").Run(func(args mock.Arguments) {
		callOrder = append(callOrder, "parse")
	}).Return(diagram, nil)

	mockLayout.On("Arrange", diagram).Run(func(args mock.Arguments) {
		callOrder = append(callOrder, "arrange")
	}).Return(nil)

	mockPathfinder.On("FindPaths", diagram).Run(func(args mock.Arguments) {
		callOrder = append(callOrder, "pathfind")
	}).Return(nil)

	mockRenderer.On("Render", diagram, "output.svg").Run(func(args mock.Arguments) {
		callOrder = append(callOrder, "render")
	}).Return(nil)

	uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)

	// Act
	err := uc.Execute("test.layli", "output.svg")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"parse", "arrange", "pathfind", "render"}, callOrder)
}
