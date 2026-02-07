package adapters

import (
	"testing"

	"github.com/dnnrly/layli/internal/domain"
)

func TestToLayoutConfig(t *testing.T) {
	diagram := &domain.Diagram{
		Config: domain.DiagramConfig{
			LayoutType:     "flow-square",
			LayoutAttempts: 10,
			NodeWidth:      5,
			NodeHeight:     3,
			Border:         1,
			Margin:         2,
			Spacing:        20,
			PathStrategy:   "in-order",
			PathAttempts:   20,
		},
		Nodes: []domain.Node{
			{
				ID:       "node1",
				Contents: "Test Node 1",
				Position: domain.Position{X: 10, Y: 20},
				Width:    5,
				Height:   3,
				Class:    "test-class",
				Style:    "fill: red",
			},
			{
				ID:       "node2",
				Contents: "Test Node 2",
				Position: domain.Position{X: 30, Y: 40},
				Width:    5,
				Height:   3,
			},
		},
		Edges: []domain.Edge{
			{
				ID:    "edge1",
				From:  "node1",
				To:    "node2",
				Class: "edge-class",
				Style:  "stroke: blue",
			},
		},
	}

	config := ToLayoutConfig(diagram)

	// Test basic config fields
	if config.Layout != "flow-square" {
		t.Errorf("Expected Layout 'flow-square', got '%s'", config.Layout)
	}
	if config.LayoutAttempts != 10 {
		t.Errorf("Expected LayoutAttempts 10, got %d", config.LayoutAttempts)
	}
	if config.NodeWidth != 5 {
		t.Errorf("Expected NodeWidth 5, got %d", config.NodeWidth)
	}
	if config.NodeHeight != 3 {
		t.Errorf("Expected NodeHeight 3, got %d", config.NodeHeight)
	}
	if config.Border != 1 {
		t.Errorf("Expected Border 1, got %d", config.Border)
	}
	if config.Margin != 2 {
		t.Errorf("Expected Margin 2, got %d", config.Margin)
	}
	if config.Spacing != 20 {
		t.Errorf("Expected Spacing 20, got %d", config.Spacing)
	}

	// Test nodes
	if len(config.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(config.Nodes))
	}

	node1 := config.Nodes[0]
	if node1.Id != "node1" {
		t.Errorf("Expected node ID 'node1', got '%s'", node1.Id)
	}
	if node1.Contents != "Test Node 1" {
		t.Errorf("Expected node contents 'Test Node 1', got '%s'", node1.Contents)
	}
	if node1.Position.X != 10 {
		t.Errorf("Expected node X position 10, got %d", node1.Position.X)
	}
	if node1.Position.Y != 20 {
		t.Errorf("Expected node Y position 20, got %d", node1.Position.Y)
	}
	if node1.Class != "test-class" {
		t.Errorf("Expected node class 'test-class', got '%s'", node1.Class)
	}
	if node1.Style != "fill: red" {
		t.Errorf("Expected node style 'fill: red', got '%s'", node1.Style)
	}

	// Test edges
	if len(config.Edges) != 1 {
		t.Fatalf("Expected 1 edge, got %d", len(config.Edges))
	}

	edge1 := config.Edges[0]
	if edge1.ID != "edge1" {
		t.Errorf("Expected edge ID 'edge1', got '%s'", edge1.ID)
	}
	if edge1.From != "node1" {
		t.Errorf("Expected edge From 'node1', got '%s'", edge1.From)
	}
	if edge1.To != "node2" {
		t.Errorf("Expected edge To 'node2', got '%s'", edge1.To)
	}
	if edge1.Class != "edge-class" {
		t.Errorf("Expected edge class 'edge-class', got '%s'", edge1.Class)
	}
	if edge1.Style != "stroke: blue" {
		t.Errorf("Expected edge style 'stroke: blue', got '%s'", edge1.Style)
	}
}

func TestToLayoutConfigWithPath(t *testing.T) {
	diagram := &domain.Diagram{
		Config: domain.DiagramConfig{
			LayoutType:   "flow-square",
			PathStrategy: "random",
			PathAttempts: 15,
		},
		Nodes: []domain.Node{
			{
				ID:       "node1",
				Contents: "Test Node 1",
				Position: domain.Position{X: 10, Y: 20},
			},
		},
		Edges: []domain.Edge{},
	}

	config := ToLayoutConfigWithPath(diagram)

	// Test that it includes path configuration
	if config.Path.Strategy != "random" {
		t.Errorf("Expected Path Strategy 'random', got '%s'", config.Path.Strategy)
	}
	if config.Path.Attempts != 15 {
		t.Errorf("Expected Path Attempts 15, got %d", config.Path.Attempts)
	}

	// Test that other fields are still populated
	if config.Layout != "flow-square" {
		t.Errorf("Expected Layout 'flow-square', got '%s'", config.Layout)
	}
	if len(config.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(config.Nodes))
	}
}

func TestToLayoutConfigEmptyDiagram(t *testing.T) {
	diagram := &domain.Diagram{
		Config: domain.DiagramConfig{},
		Nodes:  []domain.Node{},
		Edges:  []domain.Edge{},
	}

	config := ToLayoutConfig(diagram)

	if len(config.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(config.Nodes))
	}
	if len(config.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(config.Edges))
	}
}
