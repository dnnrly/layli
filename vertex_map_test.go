package layli

import (
	"strings"
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/stretchr/testify/assert"
)

var vertexTestConfig = Config{
	Nodes: ConfigNodes{
		ConfigNode{Id: "1"},
		ConfigNode{Id: "2"},
		ConfigNode{Id: "3"},
	},
	Spacing:    20,
	NodeWidth:  5,
	NodeHeight: 3,
	Margin:     2,
	Border:     1,
}

func TestVertexMap_Count(t *testing.T) {
	m := NewVertexMap(20, 20)

	assert.Equal(t, 400, m.CountAvailable(false))
	assert.Equal(t, 0, m.CountAvailable(true))

	m.Set(1, 2, true)
	m.Set(1, 3, true)
	m.Set(1, 4, true)

	assert.Equal(t, 397, m.CountAvailable(false))
	assert.Equal(t, 3, m.CountAvailable(true))
}

func TestVertexMap_String(t *testing.T) {
	m := NewVertexMap(20, 20)

	assert.Equal(t, strings.ReplaceAll(
		`....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................`, "	", ""), m.String(), m)
	assert.Equal(t, 0, m.CountAvailable(true))

	m.Set(1, 2, true)
	m.Set(1, 3, true)
	m.Set(1, 4, true)

	assert.Equal(t, strings.ReplaceAll(
		`....................
		....................
		.x..................
		.x..................
		.x..................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................
		....................`, "	", ""), m.String(), m)
}

func TestVertexMap_MapSet(t *testing.T) {
	m := NewVertexMap(20, 20)

	m.MapSet(func(x, y int) bool {
		return (x+y)%2 == 0
	})

	assert.Equal(t, strings.ReplaceAll(
		`x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x`, "	", ""), m.String(), m)
}

func TestVertexMap_MapUnset(t *testing.T) {
	m := NewVertexMap(20, 20)

	m.MapUnset(func(x, y int) bool {
		return (x+y)%2 == 0
	})

	assert.Equal(t, strings.ReplaceAll(
		`.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.
		.x.x.x.x.x.x.x.x.x.x
		x.x.x.x.x.x.x.x.x.x.`, "	", ""), m.String(), m)
}

func TestVertexMap_MapOr(t *testing.T) {
	m := NewVertexMap(5, 5)

	m.Set(2, 1, true)
	m.Set(2, 2, true)
	m.Set(2, 3, true)

	m.MapOr(func(x, y int) bool {
		return x == 2
	})

	assert.Equal(t, strings.ReplaceAll(
		`..x..
		..x..
		..x..
		..x..
		..x..`, "	", ""), m.String(), m)
}

func TestVertexMap_MapAnd(t *testing.T) {
	m := NewVertexMap(5, 5)

	m.Set(2, 1, true)
	m.Set(2, 2, true)
	m.Set(2, 3, true)

	m.MapAnd(func(x, y int) bool {
		return y == 2
	})

	assert.Equal(t, strings.ReplaceAll(
		`.....
		.....
		..x..
		.....
		.....`, "	", ""), m.String(), m)
}

func TestVertexMap_GetArcs(t *testing.T) {
	m := NewVertexMap(5, 5)

	// .....
	// .x.x.
	// .....
	// .x.x.
	// .....

	m.MapSet(func(x, y int) bool { return true })
	m.Set(1, 1, false)
	m.Set(1, 3, false)
	m.Set(3, 1, false)
	m.Set(3, 3, false)

	arcs := m.GetArcs()

	assert.Len(t, arcs, 120)

	// Some representative points
	assert.True(t, arcs.Exists(Point{X: 2, Y: 0}, Point{X: 2, Y: 4}))
	assert.True(t, arcs.Exists(Point{X: 4, Y: 4}, Point{X: 4, Y: 1}))

	// A reversed point - we go back as wekk as forwards
	assert.True(t, arcs.Exists(Point{X: 2, Y: 4}, Point{X: 2, Y: 0}))
	assert.True(t, arcs.Exists(Point{X: 2, Y: 4}, Point{X: 2, Y: 0}))

	// A point that crosses invalid points
	assert.False(t, arcs.Exists(Point{X: 0, Y: 1}, Point{X: 2, Y: 1}))

	// From and To are the same point
	assert.False(t, arcs.Exists(Point{X: 0, Y: 1}, Point{X: 0, Y: 1}))
}

func TestVertexMap_GetVertexPoints(t *testing.T) {
	m := NewVertexMap(5, 5)

	m.Set(2, 1, true)
	m.Set(2, 2, true)
	m.Set(2, 3, true)

	points := m.GetVertexPoints()

	assert.Len(t, points, 3)
	assert.Contains(t, points, Point{X: 2, Y: 1})
	assert.Contains(t, points, Point{X: 2, Y: 2})
	assert.Contains(t, points, Point{X: 2, Y: 3})
}

func TestArcs_Add(t *testing.T) {
	arcs := Arcs{}

	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}, 1)
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 5}, 1)

	assert.Len(t, arcs, 2)
	assert.Contains(t, arcs, Arc{From: Point{X: 1, Y: 2}, To: Point{X: 1, Y: 3}, Distance: 1})
	assert.Contains(t, arcs, Arc{From: Point{X: 1, Y: 2}, To: Point{X: 1, Y: 5}, Distance: 1})
}

func TestArcs_Exists(t *testing.T) {
	arcs := Arcs{}

	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}, 1)
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 5}, 1)

	assert.True(t, arcs.Exists(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}))
	assert.False(t, arcs.Exists(Point{X: 1, Y: 2}, Point{X: 1, Y: 9}))
}

func TestArcs_AddToGraph(t *testing.T) {
	g := mocks.NewGraph(t)

	arcs := Arcs{}
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}, 1)
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 5}, 1)

	g.On("AddMappedArc", "1.0,2.0", "1.0,3.0", int64(1)).Return(nil).Once()
	g.On("AddMappedArc", "1.0,2.0", "1.0,5.0", int64(1)).Return(nil).Once()

	arcs.AddToGraph(g)

	g.AssertExpectations(t)
}
