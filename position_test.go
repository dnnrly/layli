package layli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_String(t *testing.T) {
	assert.Equal(t, "4.0,7.0", Point{X: 4, Y: 7}.String())
}

func TestPoint_Coordinates(t *testing.T) {
	x, y := Point{X: 4, Y: 7}.Coordinates()
	assert.Equal(t, 4.0, x)
	assert.Equal(t, 7.0, y)
}

func TestPoint_Distance(t *testing.T) {
	a := Point{X: 2, Y: 5}
	b := Point{X: 6, Y: 7}

	assert.InDelta(t, 4.5, a.Distance(b), 0.1)
	assert.InDelta(t, 4.5, b.Distance(a), 0.1)
	assert.Equal(t, 0.0, a.Distance(a))
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
