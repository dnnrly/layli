package domain

import (
	"testing"
)

func TestNodeValidate_Success(t *testing.T) {
	tests := []struct {
		name string
		node Node
	}{
		{
			name: "valid node with positive dimensions",
			node: Node{ID: "a", Width: 5, Height: 5},
		},
		{
			name: "valid node with large dimensions",
			node: Node{ID: "node-123", Width: 1000, Height: 2000},
		},
		{
			name: "valid node with zero dimensions",
			node: Node{ID: "b", Width: 0, Height: 0},
		},
		{
			name: "valid node with contents and style",
			node: Node{
				ID:       "c",
				Contents: "Hello World",
				Width:    10,
				Height:   10,
				Class:    "class-1",
				Style:    "fill:red",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.node.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestNodeValidate_EmptyID(t *testing.T) {
	n := Node{ID: "", Width: 5, Height: 5}
	if err := n.Validate(); err == nil {
		t.Fatal("expected error for empty node ID")
	}
}

func TestNodeValidate_NegativeDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"negative width", -1, 5},
		{"negative height", 5, -1},
		{"both negative", -1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{ID: "a", Width: tt.width, Height: tt.height}
			if err := n.Validate(); err == nil {
				t.Fatalf("expected error for width=%d height=%d", tt.width, tt.height)
			}
		})
	}
}

func TestNodeBounds(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		wantMin Position
		wantMax Position
	}{
		{
			name:    "node at origin",
			node:    Node{ID: "a", Position: Position{X: 0, Y: 0}, Width: 5, Height: 10},
			wantMin: Position{X: 0, Y: 0},
			wantMax: Position{X: 5, Y: 10},
		},
		{
			name:    "node at non-origin",
			node:    Node{ID: "b", Position: Position{X: 3, Y: 7}, Width: 5, Height: 10},
			wantMin: Position{X: 3, Y: 7},
			wantMax: Position{X: 8, Y: 17},
		},
		{
			name:    "zero-sized node",
			node:    Node{ID: "c", Position: Position{X: 10, Y: 20}, Width: 0, Height: 0},
			wantMin: Position{X: 10, Y: 20},
			wantMax: Position{X: 10, Y: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.node.Bounds()
			if b.Min != tt.wantMin {
				t.Errorf("bounds.Min = %v, want %v", b.Min, tt.wantMin)
			}
			if b.Max != tt.wantMax {
				t.Errorf("bounds.Max = %v, want %v", b.Max, tt.wantMax)
			}
		})
	}
}

func TestNodeCenter(t *testing.T) {
	tests := []struct {
		name   string
		node   Node
		wantX  int
		wantY  int
	}{
		{
			name:   "node at origin with 10x10 size",
			node:   Node{Position: Position{X: 0, Y: 0}, Width: 10, Height: 10},
			wantX:  5,
			wantY:  5,
		},
		{
			name:   "node at (3,7) with 4x6 size",
			node:   Node{Position: Position{X: 3, Y: 7}, Width: 4, Height: 6},
			wantX:  5,
			wantY:  10,
		},
		{
			name:   "node with odd dimensions",
			node:   Node{Position: Position{X: 0, Y: 0}, Width: 5, Height: 7},
			wantX:  2,
			wantY:  3,
		},
		{
			name:   "zero-sized node",
			node:   Node{Position: Position{X: 10, Y: 20}, Width: 0, Height: 0},
			wantX:  10,
			wantY:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.node.Center()
			if c.X != tt.wantX {
				t.Errorf("center.X = %d, want %d", c.X, tt.wantX)
			}
			if c.Y != tt.wantY {
				t.Errorf("center.Y = %d, want %d", c.Y, tt.wantY)
			}
		})
	}
}

func TestNodeString(t *testing.T) {
	n := Node{
		ID:       "test-node",
		Position: Position{X: 3, Y: 5},
		Width:    10,
		Height:   20,
	}
	result := n.String()
	if result != "Node(test-node at (3,5) [10x20])" {
		t.Errorf("String() = %q, want %q", result, "Node(test-node at (3,5) [10x20])")
	}
}
