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
