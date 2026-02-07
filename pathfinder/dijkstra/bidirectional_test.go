package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBidirectionalPathFinder_SimplePath(t *testing.T) {
	pf := NewBidirectionalPathFinder(pt(1, 1), pt(3, 1))
	cost := uniformCost

	pf.AddConnection(pt(1, 1), cost, pt(2, 1))
	pf.AddConnection(pt(2, 1), cost, pt(3, 1))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, []Point{pt(1, 1), pt(2, 1), pt(3, 1)}, path)
}

func TestBidirectionalPathFinder_LongerPath(t *testing.T) {
	pf := NewBidirectionalPathFinder(pt(1, 1), pt(5, 1))
	cost := uniformCost

	pf.AddConnection(pt(1, 1), cost, pt(2, 1))
	pf.AddConnection(pt(2, 1), cost, pt(3, 1))
	pf.AddConnection(pt(3, 1), cost, pt(4, 1))
	pf.AddConnection(pt(4, 1), cost, pt(5, 1))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, []Point{pt(1, 1), pt(2, 1), pt(3, 1), pt(4, 1), pt(5, 1)}, path)
}

func TestBidirectionalPathFinder_GridPath(t *testing.T) {
	pf := NewBidirectionalPathFinder(pt(1, 1), pt(3, 3))
	cost := uniformCost

	// Create a grid
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

func TestBidirectionalPathFinder_NoPath(t *testing.T) {
	pf := NewBidirectionalPathFinder(pt(1, 1), pt(3, 3))
	cost := uniformCost

	// Create disconnected components
	pf.AddConnection(pt(1, 1), cost, pt(1, 2))
	pf.AddConnection(pt(3, 2), cost, pt(3, 3))

	_, err := pf.BestPath()
	assert.Equal(t, ErrNotFound, err)
}

func TestBidirectionalPathFinder_ComplexPath(t *testing.T) {
	pf := NewBidirectionalPathFinder(pt(0, 0), pt(4, 4))
	cost := uniformCost

	// Create a more complex grid
	pf.AddConnection(pt(0, 0), cost, pt(0, 1), pt(1, 0))
	pf.AddConnection(pt(0, 1), cost, pt(0, 2), pt(1, 1))
	pf.AddConnection(pt(1, 0), cost, pt(2, 0), pt(1, 1))
	pf.AddConnection(pt(0, 2), cost, pt(0, 3), pt(1, 2))
	pf.AddConnection(pt(1, 1), cost, pt(1, 2), pt(2, 1))
	pf.AddConnection(pt(2, 0), cost, pt(3, 0), pt(2, 1))
	pf.AddConnection(pt(0, 3), cost, pt(0, 4), pt(1, 3))
	pf.AddConnection(pt(1, 2), cost, pt(1, 3), pt(2, 2))
	pf.AddConnection(pt(2, 1), cost, pt(2, 2), pt(3, 1))
	pf.AddConnection(pt(3, 0), cost, pt(4, 0), pt(3, 1))
	pf.AddConnection(pt(0, 4), cost, pt(1, 4))
	pf.AddConnection(pt(1, 3), cost, pt(1, 4), pt(2, 3))
	pf.AddConnection(pt(2, 2), cost, pt(2, 3), pt(3, 2))
	pf.AddConnection(pt(3, 1), cost, pt(3, 2), pt(4, 1))
	pf.AddConnection(pt(4, 0), cost, pt(4, 1))
	pf.AddConnection(pt(1, 4), cost, pt(2, 4))
	pf.AddConnection(pt(2, 3), cost, pt(2, 4), pt(3, 3))
	pf.AddConnection(pt(3, 2), cost, pt(3, 3), pt(4, 2))
	pf.AddConnection(pt(4, 1), cost, pt(4, 2))
	pf.AddConnection(pt(2, 4), cost, pt(3, 4))
	pf.AddConnection(pt(3, 3), cost, pt(3, 4), pt(4, 3))
	pf.AddConnection(pt(4, 2), cost, pt(4, 3))
	pf.AddConnection(pt(3, 4), cost, pt(4, 4))
	pf.AddConnection(pt(4, 3), cost, pt(4, 4))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, pt(0, 0), path[0])
	assert.Equal(t, pt(4, 4), path[len(path)-1])
	assert.True(t, len(path) >= 9) // Should be at least 9 steps
}

func TestBidirectionalPathFinder_SameStartEnd(t *testing.T) {
	pf := NewBidirectionalPathFinder(pt(1, 1), pt(1, 1))

	path, err := pf.BestPath()
	require.NoError(t, err)
	assert.Equal(t, []Point{pt(1, 1)}, path)
}
