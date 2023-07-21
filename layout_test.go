package layli

import (
	"strings"
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var layoutTestConfig = Config{
	Nodes: ConfigNodes{
		ConfigNode{Id: "1"},
		ConfigNode{Id: "2"},
	},
	Spacing:    20,
	NodeWidth:  5,
	NodeHeight: 3,
	Margin:     2,
	Border:     1,
}

func TestLayout_LayoutSize(t *testing.T) {
	l := Layout{
		pathSpacing: 20,

		nodeWidth:  5,
		nodeHeight: 4,
		nodeMargin: 1,

		layoutBorder: 1,
	}

	l.Nodes = make(LayoutNodes, 2)
	assert.Equal(t, 1+1+5+2+5+1+1, l.LayoutWidth())
	assert.Equal(t, 1+1+4+1+1, l.LayoutHeight())

	l.Nodes = make(LayoutNodes, 5)
	assert.Equal(t, 1+1+5+2+5+2+5+1+1, l.LayoutWidth())
	assert.Equal(t, 1+1+4+2+4+1+1, l.LayoutHeight())

	l.Nodes = make(LayoutNodes, 8)
	assert.Equal(t, 1+1+5+2+5+2+5+1+1, l.LayoutWidth())
	assert.Equal(t, 1+1+4+2+4+2+4+1+1, l.LayoutHeight())
}

func TestLayoutNode_DrawNode(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	n := LayoutNode{
		Id:       "nodeA",
		Contents: "some contents",

		left: 4,
		top:  5,

		width:  3,
		height: 3,
	}

	drawer.On("Roundrect", 160, 200, 80, 80, 3, 3, `id="nodeA"`).Once()
	drawer.On("Textspan", 200, 240, "some contents", `id="nodeA-text"`, "font-size:10px").Once()
	drawer.On("TextEnd").Once()

	n.Draw(drawer, 40)

	drawer.AssertExpectations(t)
}

func TestLayoutNode_IsInside(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3)

	// ...........
	// ...........
	// ...........
	// ...xxxxx...
	// ...xxxxx...
	// ...xxxxx...
	// ...........
	// ...........
	// ...........

	assert.True(t, n.IsInside(3, 3), "3,3 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.True(t, n.IsInside(6, 5), "5,6 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.False(t, n.IsInside(3, 2), "3,2 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.False(t, n.IsInside(8, 3), "8,3 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
}

func TestLayoutNode_IsPort(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3)

	// ...........
	// ...........
	// ...........
	// ....xxx....
	// ...x...x...
	// ....xxx....
	// ...........
	// ...........
	// ...........

	// Inside
	assert.False(t, n.IsPort(4, 4), "4,4 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)

	// Outside
	assert.False(t, n.IsPort(3, 2), "2,2 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.False(t, n.IsPort(4, 8), "4,8 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)

	// Corner
	assert.False(t, n.IsPort(3, 3), "3,7 - corner %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)

	assert.True(t, n.IsPort(3, 4), "3,4 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.True(t, n.IsPort(6, 5), "6,5 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.True(t, n.IsPort(4, 3), "4,3 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
	assert.True(t, n.IsPort(4, 5), "4,5 %d,%d,%d,%d", n.left, n.right, n.top, n.bottom)
}

func TestLayoutNode_GetCentre(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3)

	assert.Equal(t, Point{X: 5.5, Y: 4.5}, n.GetCentre())
}

func TestLayoutNodes_ByID(t *testing.T) {
	nodes := LayoutNodes{
		NewLayoutNode("1", "contents", 3, 7, 5, 3),
		NewLayoutNode("2", "contents", 10, 12, 5, 3),
	}

	assert.Equal(t, NewLayoutNode("1", "contents", 3, 7, 5, 3), *nodes.ByID("1"))
	assert.Equal(t, NewLayoutNode("2", "contents", 10, 12, 5, 3), *nodes.ByID("2"))
}

func TestLayout_InsideAny(t *testing.T) {
	l := NewLayoutFromConfig(layoutTestConfig)

	vm := NewVertexMap(l.LayoutWidth(), l.LayoutHeight())
	vm.MapSet(l.InsideAny)
	assert.Equal(t, strings.ReplaceAll(
		`....................
		....................
		....................
		...xxxxx....xxxxx...
		...xxxxx....xxxxx...
		...xxxxx....xxxxx...
		....................
		....................
		....................`, "	", ""), vm.String(), vm)
}

func TestLayout_IsAnyPort(t *testing.T) {
	l := NewLayoutFromConfig(layoutTestConfig)

	vm := NewVertexMap(l.LayoutWidth(), l.LayoutHeight())
	vm.MapSet(l.IsAnyPort)
	assert.Equal(t, strings.ReplaceAll(
		`....................
		....................
		....................
		....xxx......xxx....
		...x...x....x...x...
		....xxx......xxx....
		....................
		....................
		....................`, "	", ""), vm.String(), vm)
}

func TestLayoutPath_Draw(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	p := LayoutPath{
		points: Points{
			Point{X: 5.5, Y: 4.5},
			Point{X: 8, Y: 4},
			Point{X: 10, Y: 4},
			Point{X: 10, Y: 5},
			Point{X: 12, Y: 5},
			Point{X: 14.5, Y: 4.5},
		},
	}

	drawer.On(
		"Path",
		"M 80 40 L 100 40 L 100 50 L 120 50",
		`class="path-line"`,
	).Once()

	p.Draw(drawer, 10)

	drawer.AssertExpectations(t)
}

func TestLayoutPaths_Draw(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	p := LayoutPaths{
		LayoutPath{points: Points{Point{X: 5.5, Y: 4.5}, Point{X: 8, Y: 4}, Point{X: 12, Y: 4}, Point{X: 14.5, Y: 4.5}}},
		LayoutPath{points: Points{Point{X: 5.5, Y: 4.5}, Point{X: 8, Y: 5}, Point{X: 12, Y: 5}, Point{X: 14.5, Y: 4.5}}},
	}

	drawer.On("Path", mock.Anything, `class="path-line"`).Twice()

	p.Draw(drawer, 10)

	drawer.AssertExpectations(t)
}
