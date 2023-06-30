package layli

import (
	"testing"

	"github.com/dnnrly/layli/mocks"
)

func TestNode_DrawNode(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	n := LayoutNode{
		Id:       "nodeA",
		Contents: "some contents",
		X:        4,
		Y:        5,
	}

	drawer.On("Roundrect", (233*4)-50, (233*5)-40, 100, 80, 3, 3, `id="nodeA"`).Once()
	drawer.On("Textspan", 233*4, 233*5, "some contents", `id="nodeA-text"`, "font-size:10px").Once()
	drawer.On("TextEnd").Once()

	n.Draw(drawer, 233, 100, 80)

	drawer.AssertExpectations(t)
}
