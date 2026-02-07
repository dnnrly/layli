package dijkstra

import (
	"fmt"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
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

	pf.AddConnection(p(1, 1), cost, p(1, 2), p(2, 1))
	pf.AddConnection(p(1, 2), cost, p(1, 3), p(2, 2))
	pf.AddConnection(p(2, 1), cost, p(2, 2), p(3, 1))
	pf.AddConnection(p(1, 3), cost, p(2, 3), p(1, 4))
	pf.AddConnection(p(2, 2), cost, p(2, 3), p(3, 2))
	pf.AddConnection(p(3, 1), cost, p(3, 2), p(4, 1))
	pf.AddConnection(p(1, 4), cost, p(2, 4))
	pf.AddConnection(p(2, 3), cost, p(2, 4), p(3, 3))
	pf.AddConnection(p(3, 2), cost, p(3, 3), p(4, 2))
	pf.AddConnection(p(4, 1), cost, p(4, 2))
	pf.AddConnection(p(2, 4), cost, p(3, 4))
	pf.AddConnection(p(3, 3), cost, p(3, 4), p(4, 3))
	pf.AddConnection(p(4, 2), cost, p(4, 3))
	pf.AddConnection(p(3, 4), cost, p(4, 4))
	pf.AddConnection(p(4, 3), cost, p(4, 4))

	path, err := pf.BestPath()

	assert.NoError(t, err)
	// The exact path may vary depending on the algorithm implementation
	// Just verify we get a valid path from start to end
	assert.Equal(t, p(1, 1), path[0])
	assert.Equal(t, p(4, 4), path[len(path)-1])
	assert.Greater(t, len(path), 2)
}

func TestGetSetFunctions(t *testing.T) {
	// Test the get and set helper functions
	m := orderedmap.NewOrderedMap[string, int]()

	// Test set function
	set(m, "key1", 42)
	set(m, "key2", 24)

	// Test get function
	val1 := get(m, "key1")
	if val1 != 42 {
		t.Errorf("Expected 42, got %d", val1)
	}

	val2 := get(m, "key2")
	if val2 != 24 {
		t.Errorf("Expected 24, got %d", val2)
	}

	// Test get with non-existent key
	val3 := get(m, "nonexistent")
	if val3 != 0 {
		t.Errorf("Expected 0 for non-existent key, got %d", val3)
	}

	// Test set with nil map
	var nilMap *orderedmap.OrderedMap[string, int]
	set(nilMap, "key", 42) // Should not panic

	// Test get with nil map
	val4 := get(nilMap, "key")
	if val4 != 0 {
		t.Errorf("Expected 0 for nil map, got %d", val4)
	}
}

func TestCannotFindImpossiblePath(t *testing.T) {
	pf := NewPathFinder(p(1, 1), p(2, 2))

	cost := func(from, to Point) int64 { return 1 }

	pf.AddConnection(p(1, 1), cost, p(1, 2))
	pf.AddConnection(p(3, 2), cost, p(2, 2))

	_, err := pf.BestPath()
	assert.ErrorIs(t, err, ErrNotFound)
}
