package layli

import (
	"testing"

	"github.com/dnnrly/layli/mocks"
)

func TestPoints_Draw(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	p := Points{
		Point{X: 5.5, Y: 4.5},
		Point{X: 8, Y: 4},
		Point{X: 10, Y: 4},
		Point{X: 10, Y: 5},
		Point{X: 12, Y: 5},
		Point{X: 14.5, Y: 4.5},
	}

	drawer.On(
		"Path",
		"M 80 40 L 100 40 L 100 50 L 120 50",
		`class="path-line"`,
	).Once()

	p.Draw(drawer, 10)

	drawer.AssertExpectations(t)
}
