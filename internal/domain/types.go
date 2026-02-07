package domain

import (
	"fmt"
	"math"
)

// Position represents an X, Y coordinate on the grid.
type Position struct {
	X int
	Y int
}

// String returns a string representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Distance returns the Euclidean distance between two positions.
func (p Position) Distance(other Position) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// Bounds represents a rectangular boundary.
type Bounds struct {
	Min Position // Top-left corner
	Max Position // Bottom-right corner
}

// String returns a string representation of the bounds.
func (b Bounds) String() string {
	return fmt.Sprintf("Bounds(%v to %v)", b.Min, b.Max)
}

// Width returns the width of the bounds.
func (b Bounds) Width() int {
	return b.Max.X - b.Min.X
}

// Height returns the height of the bounds.
func (b Bounds) Height() int {
	return b.Max.Y - b.Min.Y
}

// Contains checks if a position is within the bounds.
func (b Bounds) Contains(p Position) bool {
	return p.X >= b.Min.X && p.X <= b.Max.X &&
		p.Y >= b.Min.Y && p.Y <= b.Max.Y
}
