package domain

import (
	"fmt"
	"math"
)

// Position represents X/Y coordinates on the grid.
type Position struct {
	X int
	Y int
}

// Distance calculates the Euclidean distance to another position.
func (p Position) Distance(to Position) float64 {
	dx := float64(p.X - to.X)
	dy := float64(p.Y - to.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// Add returns a new position with the offset applied.
func (p Position) Add(dx, dy int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy}
}

// String returns a string representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Bounds represents a rectangular area with min and max positions.
type Bounds struct {
	Min Position
	Max Position
}

// Width returns the width of the bounds.
func (b Bounds) Width() int {
	return b.Max.X - b.Min.X
}

// Height returns the height of the bounds.
func (b Bounds) Height() int {
	return b.Max.Y - b.Min.Y
}

// Contains returns true if the position is within the bounds.
func (b Bounds) Contains(p Position) bool {
	return p.X >= b.Min.X && p.X <= b.Max.X &&
		p.Y >= b.Min.Y && p.Y <= b.Max.Y
}

// Overlaps returns true if this bounds overlaps with another bounds.
func (b Bounds) Overlaps(other Bounds) bool {
	return !(b.Max.X < other.Min.X || b.Min.X > other.Max.X ||
		b.Max.Y < other.Min.Y || b.Min.Y > other.Max.Y)
}

// String returns a string representation of the bounds.
func (b Bounds) String() string {
	return fmt.Sprintf("[%s to %s]", b.Min.String(), b.Max.String())
}
