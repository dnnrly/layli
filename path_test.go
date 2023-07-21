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

	if assert.NoError(t, l.AddPath("1", "2"), l.vertexMap.String()) {
		assert.Len(t, l.Paths, 1)
		assert.Contains(t, l.Paths[0].points, Point{X: 5.5, Y: 5.5})
		assert.Contains(t, l.Paths[0].points, Point{X: 12.5, Y: 5.5})
	}
}

func TestLayout_BuildVertexMap(t *testing.T) {
	l := NewLayoutFromConfig(pathTestConfig)
	vm := BuildVertexMap(l)

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
		xxxxxxxxxxxxxxxxxx`, "	", ""), vm.String(), vm)
}
