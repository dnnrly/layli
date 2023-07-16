package layli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pathTestConfig = Config{
	Nodes: ConfigNodes{
		ConfigNode{Id: "1"},
		ConfigNode{Id: "2"},
	},
	Spacing:    20,
	NodeWidth:  3,
	NodeHeight: 3,
	Margin:     2,
	Border:     2,
}

func TestLayout_AddPath_BetweenAdjacentNodes(t *testing.T) {
	l := NewLayoutFromConfig(pathTestConfig)

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

func TestLayout_BuildVertexMap(t *testing.T) {
	l := NewLayoutFromConfig(pathTestConfig)

	l.BuildVertexMap()

	assert.Equal(t, strings.ReplaceAll(
		`xxxxxxxxxxxxxxxxxx
		xxxxxxxxxxxxxxxxxx
		xxxxxxxxxxxxxxxxxx
		xxxxxxxxxxxxxxxxxx
		xxxx.x.xxxx.x.xxxx
		xxxxx.xxxxxx.xxxxx
		xxxx.x.xxxx.x.xxxx
		xxxxxxxxxxxxxxxxxx
		xxxxxxxxxxxxxxxxxx
		xxxxxxxxxxxxxxxxxx
		xxxxxxxxxxxxxxxxxx`, "	", ""), l.vertexMap.String(), l.vertexMap)
}
