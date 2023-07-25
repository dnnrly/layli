package djikstra

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

	path := pf.BestPath()

	assert.Equal(t,
		fmt.Sprintf("%v", []Point{p(1, 1), p(1, 2), p(2, 2)}),
		fmt.Sprintf("%v", path),
	)
}
