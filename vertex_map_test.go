package layli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestVertexMap_Get(t *testing.T) {
	m := NewVertexMap(20, 20)

	assert.Equal(t, 400, m.CountAvailable(false))
	assert.Equal(t, 0, m.CountAvailable(true))

	m.Set(1, 2, true)
	m.Set(1, 3, true)
	m.Set(1, 4, true)

	assert.True(t, m.Get(1, 2))
	assert.False(t, m.Get(2, 1))
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

func TestVertexMap_GetArcs_CheckDistances(t *testing.T) {
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

	// Some representative points
	assert.Equal(t, 4000, arcs.Get(Point{X: 2, Y: 0}, Point{X: 2, Y: 4}).Distance)
	assert.Equal(t, 3000, arcs.Get(Point{X: 4, Y: 4}, Point{X: 4, Y: 1}).Distance)
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

func TestArcs_Get(t *testing.T) {
	arcs := Arcs{}

	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}, 1)
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 5}, 2)

	assert.Equal(t, Arc{From: Point{X: 1, Y: 2}, To: Point{X: 1, Y: 3}, Distance: 1}, arcs.Get(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}))
	assert.Equal(t, Arc{}, arcs.Get(Point{X: 1, Y: 2}, Point{X: 1, Y: 9}))
}

func TestArcs_String(t *testing.T) {
	arcs := Arcs{}
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 3}, 101)
	arcs.Add(Point{X: 1, Y: 2}, Point{X: 1, Y: 5}, 103)

	assert.Equal(t, `1.0,2.0-1.0,3.0-101
1.0,2.0-1.0,5.0-103`, arcs.String())
}
