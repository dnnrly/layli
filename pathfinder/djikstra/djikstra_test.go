package djikstra_test

import (
	"testing"

	"github.com/dnnrly/layli"
	"github.com/dnnrly/layli/pathfinder/djikstra"
	"github.com/stretchr/testify/assert"
)

func p(x, y float64) layli.Point {
	return layli.Point{X: x, Y: y}
}

func TestAddVertex(t *testing.T) {
	pf := djikstra.NewPathFinder(p(1, 1), p(2, 2))

	cost := func(from, to layli.Point) int64 { return 1 }

	pf.AddConnection(p(1, 1), cost, p(1, 2))
	pf.AddConnection(p(1, 2), cost, p(2, 2))

	path := pf.BestPath()

	assert.Equal(t, layli.LayoutPath{
		Points: layli.Points{p(1, 1), p(1, 2), p(2, 2)},
	}, path)
}
