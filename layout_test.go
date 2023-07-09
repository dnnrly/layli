package layli

import (
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/stretchr/testify/assert"
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
	n := NewLayoutNode("id", "contents", 3, 7, 5, 3)

	assert.True(t, n.IsInside(3, 7))
	assert.True(t, n.IsInside(8, 10))
	assert.False(t, n.IsInside(3, 6))
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
