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

func TestVertexMap_GetVertexIDs(t *testing.T) {
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
