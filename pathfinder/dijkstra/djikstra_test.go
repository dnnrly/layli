package dijkstra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func p(x, y float64) coordinate {
	return coordinate{x: x, y: y}
}

func TestAddVertex(t *testing.T) {
	pf := NewPathFinder(p(1, 1), p(2, 2))

	cost := func(from, to Point) int64 { return 1 }

	pf.AddConnection(p(1, 1), cost, p(1, 2))
	pf.AddConnection(p(1, 2), cost, p(2, 2))

	path, err := pf.BestPath()

	assert.NoError(t, err)
	assert.Equal(t,
		fmt.Sprintf("%v", []Point{p(1, 1), p(1, 2), p(2, 2)}),
		fmt.Sprintf("%v", path),
	)
}

func TestAddManyVertices(t *testing.T) {
	pf := NewPathFinder(p(1, 1), p(2, 3))

	cost := func(from, to Point) int64 { return 1 }

	pf.AddConnection(p(1, 1), cost, p(1, 2), p(1, 3), p(1, 4), p(1, 5))
	pf.AddConnection(p(1, 2), cost, p(2, 2), p(3, 2), p(4, 2), p(5, 2))
	pf.AddConnection(p(2, 2), cost, p(3, 2), p(2, 3), p(2, 4), p(2, 5))

	path, err := pf.BestPath()

	assert.NoError(t, err)
	assert.Equal(t,
		fmt.Sprintf("%v", []Point{p(1, 1), p(1, 2), p(2, 2), p(2, 3)}),
		fmt.Sprintf("%v", path),
	)
}

func TestPathFinderWithCycles(t *testing.T) {
	pf := NewPathFinder(p(1, 1), p(4, 4))

	cost := func(from, to Point) int64 { return 1 }

	pf.AddConnection(p(1, 1), cost, p(1, 2))
	pf.AddConnection(p(1, 2), cost, p(1, 3), p(2, 2))
	pf.AddConnection(p(1, 3), cost, p(1, 4), p(2, 3))
	pf.AddConnection(p(2, 2), cost, p(3, 2))
	pf.AddConnection(p(2, 3), cost, p(3, 3))
	pf.AddConnection(p(3, 2), cost, p(4, 2))
	pf.AddConnection(p(3, 3), cost, p(4, 3))
	pf.AddConnection(p(4, 2), cost, p(4, 3))
	pf.AddConnection(p(4, 3), cost, p(4, 4))

	shortestPath, err := pf.BestPath()

	assert.NoError(t, err)
	assert.Equal(t,
		fmt.Sprintf("%v", []Point{p(1, 1), p(1, 2), p(2, 2), p(3, 2), p(4, 2), p(4, 3), p(4, 4)}),
		fmt.Sprintf("%v", shortestPath),
	)
}

func TestCannotFindImpossiblePath(t *testing.T) {
	pf := NewPathFinder(p(1, 1), p(2, 2))

	cost := func(from, to Point) int64 { return 1 }

	pf.AddConnection(p(1, 1), cost, p(1, 2))
	pf.AddConnection(p(3, 2), cost, p(2, 2))

	_, err := pf.BestPath()
	assert.Contains(t, err.Error(), "no path found")
}
