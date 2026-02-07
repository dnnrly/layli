package domain

import "fmt"

// Node represents a box/component in the diagram.
// Maps to: "And a node 'A'" in feature files.
type Node struct {
	ID       string
	Contents string
	Position Position
	Width    int
	Height   int
	Class    string
	Style    string
}

// Validate ensures node invariants.
func (n *Node) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}

	if n.Width < 0 || n.Height < 0 {
		return fmt.Errorf("node dimensions must be non-negative")
	}

	return nil
}

// Bounds returns the rectangular bounds of the node.
func (n *Node) Bounds() Bounds {
	return Bounds{
		Min: n.Position,
		Max: Position{
			X: n.Position.X + n.Width,
			Y: n.Position.Y + n.Height,
		},
	}
}

// Center returns the center point of the node.
func (n *Node) Center() Position {
	return Position{
		X: n.Position.X + n.Width/2,
		Y: n.Position.Y + n.Height/2,
	}
}

// String returns a string representation of the node.
func (n *Node) String() string {
	return fmt.Sprintf("Node(%s at %s [%dx%d])", n.ID, n.Position.String(), n.Width, n.Height)
}
