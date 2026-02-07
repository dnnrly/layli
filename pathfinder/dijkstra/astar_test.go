package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pt(x, y float64) Point {
	return coordinate{x: x, y: y}
}

func uniformCost(from, to Point) int64 {
	return 1
}

func TestAStarPathFinder_SimplePath(t *testing.T) {
	pf := NewAStarPathFinder(pt(1, 1), pt(3, 1), EuclideanDistance)
	cost := uniformCost

	pf.AddConnection(pt(1, 1), cost, pt(2, 1))
	pf.AddConnection(pt(2, 1), cost, pt(3, 1))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, []Point{pt(1, 1), pt(2, 1), pt(3, 1)}, path)
}

func TestAStarPathFinder_WithEuclideanHeuristic(t *testing.T) {
	pf := NewAStarPathFinderWithEuclidean(pt(1, 1), pt(3, 3))
	cost := uniformCost

	// Create a simple grid
	pf.AddConnection(pt(1, 1), cost, pt(1, 2), pt(2, 1))
	pf.AddConnection(pt(1, 2), cost, pt(1, 3), pt(2, 2))
	pf.AddConnection(pt(2, 1), cost, pt(2, 2), pt(3, 1))
	pf.AddConnection(pt(2, 2), cost, pt(2, 3), pt(3, 2))
	pf.AddConnection(pt(1, 3), cost, pt(2, 3))
	pf.AddConnection(pt(3, 1), cost, pt(3, 2))
	pf.AddConnection(pt(2, 3), cost, pt(3, 3))
	pf.AddConnection(pt(3, 2), cost, pt(3, 3))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, 5, len(path)) // Should find a path of length 5
	assert.Equal(t, pt(1, 1), path[0])
	assert.Equal(t, pt(3, 3), path[len(path)-1])
}

func TestAStarPathFinder_WithManhattanHeuristic(t *testing.T) {
	pf := NewAStarPathFinderWithManhattan(pt(1, 1), pt(3, 3))
	cost := uniformCost

	// Create a simple grid
	pf.AddConnection(pt(1, 1), cost, pt(1, 2), pt(2, 1))
	pf.AddConnection(pt(1, 2), cost, pt(1, 3), pt(2, 2))
	pf.AddConnection(pt(2, 1), cost, pt(2, 2), pt(3, 1))
	pf.AddConnection(pt(2, 2), cost, pt(2, 3), pt(3, 2))
	pf.AddConnection(pt(1, 3), cost, pt(2, 3))
	pf.AddConnection(pt(3, 1), cost, pt(3, 2))
	pf.AddConnection(pt(2, 3), cost, pt(3, 3))
	pf.AddConnection(pt(3, 2), cost, pt(3, 3))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, 5, len(path)) // Should find a path of length 5
	assert.Equal(t, pt(1, 1), path[0])
	assert.Equal(t, pt(3, 3), path[len(path)-1])
}

func TestAStarPathFinder_NoPath(t *testing.T) {
	pf := NewAStarPathFinder(pt(1, 1), pt(3, 3), EuclideanDistance)
	cost := uniformCost

	// Create disconnected components
	pf.AddConnection(pt(1, 1), cost, pt(1, 2))
	pf.AddConnection(pt(3, 2), cost, pt(3, 3))

	_, err := pf.BestPath()
	assert.Equal(t, ErrNotFound, err)
}

func TestEuclideanDistance(t *testing.T) {
	from := pt(0, 0)
	to := pt(3, 4)
	
	distance := EuclideanDistance(from, to)
	assert.Equal(t, int64(5), distance) // 3-4-5 triangle
}

func TestManhattanDistance(t *testing.T) {
	from := pt(0, 0)
	to := pt(3, 4)
	
	distance := ManhattanDistance(from, to)
	assert.Equal(t, int64(7), distance) // |3-0| + |4-0| = 7
}
