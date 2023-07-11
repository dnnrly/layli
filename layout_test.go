package layli

import (
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	drawer.On("Roundrect", 160, 200, 120, 120, 3, 3, `id="nodeA"`).Once()
	drawer.On("Textspan", 220, 260, "some contents", `id="nodeA-text"`, "font-size:10px").Once()
	drawer.On("TextEnd").Once()

	n.Draw(drawer, 40)

	drawer.AssertExpectations(t)
}

func TestLayoutNode_IsInside(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3)

	assert.True(t, n.IsInside(3, 3))
	assert.True(t, n.IsInside(5, 6))
	assert.False(t, n.IsInside(4, 7))
	assert.False(t, n.IsInside(8, 12))
}

func TestLayoutNode_IsPort(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 7, 5, 3)

	// Inside
	assert.False(t, n.IsPort(4, 8), "4,8 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	// Outside
	assert.False(t, n.IsPort(2, 23), "4,8 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)

	// Corner
	assert.False(t, n.IsPort(3, 7), "3,7 - corner %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(3, 9), "3,9 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(8, 8), "8,7 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(4, 10), "3,10 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(6, 10), "8,10 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
}

func TestLayoutNode_GetPorts(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3)

	ports := n.GetPorts()
	assert.Len(t, ports, 12)

	// Left points
	assert.Contains(t, ports, Point{X: 3, Y: 4})
	assert.Contains(t, ports, Point{X: 3, Y: 5})

	// Right points
	assert.Contains(t, ports, Point{X: 8, Y: 4})
	assert.Contains(t, ports, Point{X: 8, Y: 5})

	// Top points
	assert.Contains(t, ports, Point{X: 4, Y: 3})
	assert.Contains(t, ports, Point{X: 5, Y: 3})
	assert.Contains(t, ports, Point{X: 6, Y: 3})
	assert.Contains(t, ports, Point{X: 7, Y: 3})

	// Bottom points
	assert.Contains(t, ports, Point{X: 4, Y: 6})
	assert.Contains(t, ports, Point{X: 5, Y: 6})
	assert.Contains(t, ports, Point{X: 6, Y: 6})
	assert.Contains(t, ports, Point{X: 7, Y: 6})
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
	l := &Layout{
		Nodes: LayoutNodes{
			NewLayoutNode("1", "contents", 3, 7, 5, 3),
			NewLayoutNode("2", "contents", 10, 12, 5, 3),
		},
	}

	assert.True(t, l.InsideAny(4, 8))
	assert.True(t, l.InsideAny(12, 15))
	assert.False(t, l.InsideAny(9, 10))
}

func TestLayout_IsAnyPort(t *testing.T) {
	l := &Layout{
		Nodes: LayoutNodes{
			NewLayoutNode("1", "contents", 3, 7, 5, 3),
			NewLayoutNode("2", "contents", 10, 12, 5, 3),
		},
	}

	assert.False(t, l.IsAnyPort(3, 7))
	assert.True(t, l.IsAnyPort(10, 8))
	assert.False(t, l.IsAnyPort(9, 9))
	assert.True(t, l.IsAnyPort(10, 13))
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
