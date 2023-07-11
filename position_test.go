package layli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
