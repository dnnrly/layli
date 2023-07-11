package layli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayout_AddPath_BetweenAdjacentNodes(t *testing.T) {
	l := &Layout{
		Nodes: LayoutNodes{
			NewLayoutNode("1", "contents", 3, 3, 2, 2),
			NewLayoutNode("2", "contents", 8, 3, 2, 2),
		},
	}

	l.AddPath("1", "2")

	assert.Len(t, l.Paths, 1)
	assert.Equal(t,
		LayoutPath{points: Points{
			Point{X: 4, Y: 4},
			Point{X: 5, Y: 4},
			Point{X: 8, Y: 4},
			Point{X: 9, Y: 4},
		}}, l.Paths)
}
