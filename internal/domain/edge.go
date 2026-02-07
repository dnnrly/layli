package domain

import "fmt"

// Edge represents a connection between two nodes.
// Maps to: "And an edge from 'A' to 'B'" in feature files.
type Edge struct {
	ID    string
	From  string  // Node ID
	To    string  // Node ID
	Path  *Path   // Calculated path (may be nil before pathfinding)
	Class string
	Style string
}

// Validate ensures edge invariants.
func (e *Edge) Validate() error {
	if e.From == "" || e.To == "" {
		return fmt.Errorf("edge must have from and to nodes")
	}

	if e.From == e.To {
		return fmt.Errorf("edge cannot connect node to itself")
	}

	return nil
}

// String returns a string representation of the edge.
func (e *Edge) String() string {
	return fmt.Sprintf("Edge(%s -> %s)", e.From, e.To)
}

// Path represents a series of positions forming an edge route.
type Path struct {
	Points []Position
}

// Length returns the total length of the path.
func (p *Path) Length() float64 {
	if p == nil || len(p.Points) < 2 {
		return 0
	}

	total := 0.0
	for i := 1; i < len(p.Points); i++ {
		total += p.Points[i-1].Distance(p.Points[i])
	}
	return total
}

// Corners returns the number of direction changes in the path.
func (p *Path) Corners() int {
	if p == nil || len(p.Points) < 3 {
		return 0
	}

	corners := 0
	for i := 1; i < len(p.Points)-1; i++ {
		prev := p.Points[i-1]
		curr := p.Points[i]
		next := p.Points[i+1]

		// Check if direction changed
		prevDx := curr.X - prev.X
		prevDy := curr.Y - prev.Y
		nextDx := next.X - curr.X
		nextDy := next.Y - curr.Y

		// If the direction vector changed, we have a corner
		if (prevDx != nextDx) || (prevDy != nextDy) {
			corners++
		}
	}

	return corners
}

// String returns a string representation of the path.
func (p *Path) String() string {
	if p == nil || len(p.Points) == 0 {
		return "Path(empty)"
	}
	return fmt.Sprintf("Path(%d points)", len(p.Points))
}
