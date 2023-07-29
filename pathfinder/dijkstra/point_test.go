package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_String(t *testing.T) {
	p := coordinate{
		x: 2,
		y: 4,
	}

	assert.Equal(t, "2.0,4.0", p.String())
}

func TestCoordinate(t *testing.T) {
	p := coordinate{
		x: 2,
		y: 4,
	}

	x, y := p.Coordinates()

	assert.Equal(t, 2.0, x)
	assert.Equal(t, 4.0, y)
}
