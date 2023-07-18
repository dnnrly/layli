package layli

import (
	"errors"
	"testing"

	"github.com/RyanCarrier/dijkstra"
	"github.com/dnnrly/layli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPoint_String(t *testing.T) {
	assert.Equal(t, "4.0,7.0", Point{X: 4, Y: 7}.String())
}

func TestPoints_Draw(t *testing.T) {
	p := Points{
		Point{X: 5.5, Y: 4.5},
		Point{X: 8, Y: 4},
		Point{X: 10, Y: 4},
		Point{X: 10, Y: 5},
		Point{X: 12, Y: 5},
		Point{X: 14.5, Y: 4.5},
	}

	assert.Equal(t, "M 80 40 L 100 40 L 100 50 L 120 50", p.Path(10))
}

func TestPoints_AddToGraph(t *testing.T) {
	g := mocks.NewGraph(t)

	points := Points{
		Point{X: 1, Y: 2},
		Point{X: 5, Y: 3},
	}

	g.On("AddMappedVertex", "1.0,2.0").Return(1).Once()
	g.On("AddMappedVertex", "5.0,3.0").Return(2).Once()

	points.AddToGraph(g)

	g.AssertExpectations(t)
}

func TestPoints_NewPointsFromBestPath(t *testing.T) {
	path := dijkstra.BestPath{
		Distance: 99,
		Path:     []int{1, 3, 5},
	}
	g := mocks.NewGraph(t)

	g.On("GetMapped", 1).Return("1.0,2.0", error(nil)).Once()
	g.On("GetMapped", 3).Return("2.0,3.0", error(nil)).Once()
	g.On("GetMapped", 5).Return("5.0,2.0", error(nil)).Once()

	points := NewPointsFromBestPath(g, path)

	g.AssertExpectations(t)

	assert.Equal(t, Points{
		Point{X: 1, Y: 2},
		Point{X: 2, Y: 3},
		Point{X: 5, Y: 2},
	}, points)
}

func TestPoints_NewPointsFromBestPath_PanicsOnUnmapped(t *testing.T) {
	path := dijkstra.BestPath{
		Distance: 99,
		Path:     []int{1},
	}
	g := mocks.NewGraph(t)

	g.On("GetMapped", 1).Return("1.0,2.0", errors.New("some error"))

	assert.Panics(t, func() { NewPointsFromBestPath(g, path) })
}
