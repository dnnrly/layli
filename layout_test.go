package layli

import (
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLayoutNode_DrawNode(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	n := LayoutNode{
		Id:       "nodeA",
		Contents: "some contents",
		X:        4 * 233,
		Y:        5 * 233,

		spacing: 233,
		left:    (233 * 4) - 50,
		top:     (233 * 5) - 40,
	}

	drawer.On("Roundrect", (233*4)-50, (233*5)-40, 100, 80, 3, 3, `id="nodeA"`).Once()
	drawer.On("Textspan", 233*4, 233*5, "some contents", `id="nodeA-text"`, "font-size:10px").Once()
	drawer.On("TextEnd").Once()

	n.Draw(drawer, 233, 100, 80)

	drawer.AssertExpectations(t)
}

func TestLayoutNode_IsInside(t *testing.T) {
	n := NewLayoutNode("id", "contents", 200, 400, 20, 80, 80)

	assert.True(t, n.IsInside(200, 400))
	assert.True(t, n.IsInside(160, 400))
	assert.False(t, n.IsInside(200, 200))
	assert.False(t, n.IsInside(159, 200))
}

func TestLayoutNode_IsPort(t *testing.T) {
	n := NewLayoutNode("id", "contents", 200, 400, 20, 80, 80)

	assert.False(t, n.IsPort(200, 400), "200,400 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.False(t, n.IsPort(201, 360), "201,360 - not a valid position")
	assert.False(t, n.IsPort(160, 360), "160,360 - corner %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(160, 420), "160,420 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(240, 380), "240,380 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(180, 360), "180,360 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
	assert.True(t, n.IsPort(220, 440), "220,440 %d,%d,%d,%d", n.top, n.bottom, n.left, n.right)
}

func TestLayout_InsideAny(t *testing.T) {
	l := &Layout{
		Nodes: LayoutNodes{
			NewLayoutNode("1", "contents", 200, 200, 20, 80, 80),
			NewLayoutNode("2", "contents", 200, 400, 20, 80, 80),
		},
	}

	assert.True(t, l.InsideAny(200, 200))
	assert.True(t, l.InsideAny(200, 400))
	assert.False(t, l.InsideAny(200, 300))
}

func TestLayout_IsAnyPort(t *testing.T) {
	l := &Layout{
		Nodes: LayoutNodes{
			NewLayoutNode("1", "contents", 200, 200, 20, 80, 80),
			NewLayoutNode("2", "contents", 200, 400, 20, 80, 80),
		},
	}

	assert.False(t, l.IsAnyPort(200, 200))
	assert.True(t, l.IsAnyPort(160, 200))
	assert.False(t, l.IsAnyPort(220, 300))
	assert.True(t, l.IsAnyPort(180, 360))
}
