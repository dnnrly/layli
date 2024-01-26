package layli

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func newConfig(nodes, width, height, margin, border int) *Config {
	c := &Config{
		Spacing:    20,
		NodeWidth:  width,
		NodeHeight: height,
		Margin:     margin,
		Border:     1,
	}

	for i := 0; i < nodes; i++ {
		c.Nodes = append(c.Nodes, ConfigNode{Id: fmt.Sprintf("%d", i+1)})
	}

	return c
}

func TestLayout_LayoutSize(t *testing.T) {
	width := 5
	height := 4
	margin := 1
	border := 1
	l := Layout{
		pathSpacing: 20,

		nodeWidth:  width,
		nodeHeight: height,
		nodeMargin: margin,

		layoutBorder: border,
	}

	l.Nodes, _ = LayoutFlowSquare(newConfig(1, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2*1), l.LayoutWidth(), "Expected width: 1 nodes")
	assert.Equal(t, border*2+(height+margin*2)*1, l.LayoutHeight(), "Expected height: 1 node")

	l.Nodes, _ = LayoutFlowSquare(newConfig(2, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2)*2, l.LayoutWidth(), "Expected width: 2 nodes")
	assert.Equal(t, border*2+(height+margin*2), l.LayoutHeight(), "Expected height: 1 node")

	l.Nodes, _ = LayoutFlowSquare(newConfig(4, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2)*2, l.LayoutWidth(), "Expected width: 2 nodes")
	assert.Equal(t, border*2+(height+margin*2)*2, l.LayoutHeight(), "Expected height: 2 nodes")

	l.Nodes, _ = LayoutFlowSquare(newConfig(5, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2)*3, l.LayoutWidth(), "Expected width: 3 nodes")
	assert.Equal(t, border*2+(height+margin*2)*2, l.LayoutHeight(), "Expected height: 2 nodes")

	l.Nodes, _ = LayoutFlowSquare(newConfig(8, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2)*3, l.LayoutWidth(), "Expected width: 3 nodes")
	assert.Equal(t, border*2+(height+margin*2)*3, l.LayoutHeight(), "Expected height: 3 nodes")

	l.Nodes, _ = LayoutFlowSquare(newConfig(9, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2)*3, l.LayoutWidth(), "Expected width: 3 nodes")
	assert.Equal(t, border*2+(height+margin*2)*3, l.LayoutHeight(), "Expected height: 3 nodes")

	l.Nodes, _ = LayoutFlowSquare(newConfig(10, width, height, margin, margin))
	assert.Equal(t, border*2+(width+margin*2)*4, l.LayoutWidth(), "Expected width: 4 nodes")
	assert.Equal(t, border*2+(height+margin*2)*3, l.LayoutHeight(), "Expected height: 3 nodes")
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

	drawer.On("Roundrect", 160, 200, 80, 80, 3, 3, `id="nodeA"`, "").Once()
	drawer.On("Textspan", 200, 240, "some contents", `id="nodeA-text"`, "font-size:10px").Once()
	drawer.On("TextEnd").Once()

	n.Draw(drawer, 40)

	drawer.AssertExpectations(t)
}

func TestLayoutNode_IsInside(t *testing.T) {
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3, "")

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
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3, "")

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
	n := NewLayoutNode("id", "contents", 3, 3, 5, 3, "")

	assert.Equal(t, Point{X: 5.5, Y: 4.5}, n.GetCentre())
}

func TestLayoutNodes_ByID(t *testing.T) {
	nodes := LayoutNodes{
		NewLayoutNode("1", "contents", 3, 7, 5, 3, ""),
		NewLayoutNode("2", "contents", 10, 12, 5, 3, ""),
	}

	assert.Equal(t, NewLayoutNode("1", "contents", 3, 7, 5, 3, ""), *nodes.ByID("1"))
	assert.Equal(t, NewLayoutNode("2", "contents", 10, 12, 5, 3, ""), *nodes.ByID("2"))
	assert.Nil(t, nodes.ByID("unknown"))
}

func TestLayoutNodes_ConnectionDistances_simple(t *testing.T) {
	n := LayoutNodes{
		NewLayoutNode("1", "contents", 1, 1, 3, 3, ""),
		NewLayoutNode("2", "contents", 1, 5, 3, 7, ""),
		NewLayoutNode("3", "contents", 1, 9, 3, 11, ""),
	}
	e := func(f, t string) ConfigEdge {
		return ConfigEdge{From: f, To: t}
	}

	dist, err := n.ConnectionDistances(ConfigEdges{e("1", "2"), e("2", "3"), e("3", "1")})

	assert.NoError(t, err)
	assert.Equal(t, 24.0, dist)
}

func TestLayoutNodes_ConnectionDistances_notFound(t *testing.T) {
	n := LayoutNodes{
		NewLayoutNode("1", "contents", 1, 1, 3, 3, ""),
		NewLayoutNode("2", "contents", 1, 5, 3, 7, ""),
		NewLayoutNode("3", "contents", 1, 9, 3, 11, ""),
	}
	e := func(f, t string) ConfigEdge {
		return ConfigEdge{From: f, To: t}
	}

	_, err := n.ConnectionDistances(ConfigEdges{e("1", "X"), e("2", "3"), e("3", "1")})
	assert.Error(t, err)

	_, err = n.ConnectionDistances(ConfigEdges{e("1", "2"), e("2", "3"), e("5", "1")})
	assert.Error(t, err)
}

func TestLayout_ErrorsOnBadLayoutName(t *testing.T) {
	_, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder { return nil }, &Config{Layout: "bad name"})
	require.Error(t, err)

}

func TestLayout_ErrorsOnBadPathStrategy(t *testing.T) {
	_, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder { return nil }, &Config{Path: ConfigPath{Strategy: "unknown"}})
	require.Error(t, err)
}

func TestLayout_PassesArrangementErrorsBack(t *testing.T) {
	// Not too happy with this, it relies on the implementation of a single arrangement algorithm to generate the error.
	// It'll do for now...
	config := &Config{
		Layout: "absolute",
		Nodes: ConfigNodes{
			ConfigNode{Id: "1", Position: Position{X: 10, Y: 10}},
			ConfigNode{Id: "2", Position: Position{X: 10, Y: 10}},
		},

		Spacing:    1,
		NodeWidth:  5,
		NodeHeight: 4,
		Margin:     2,
		Border:     1,
	}
	_, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder { return nil }, config)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "arranging nodes")
}

func TestLayout_InsideAny(t *testing.T) {
	finder := mocks.NewPathFinder(t)
	l, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder { return finder }, &layoutTestConfig)
	require.NoError(t, err)

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
	finder := mocks.NewPathFinder(t)
	l, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder { return finder }, &layoutTestConfig)
	require.NoError(t, err)

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
		Points: Points{
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
		`marker-end="url(#arrow)"`,
	).Once()

	p.Draw(drawer, 10)

	drawer.AssertExpectations(t)
}

func TestLayoutPath_Length(t *testing.T) {
	path := LayoutPath{Points: Points{Point{X: 5, Y: 4}, Point{X: 8, Y: 4}, Point{X: 12, Y: 4}, Point{X: 12, Y: 3}}}
	assert.Equal(t, 8.0, path.Length())
}

func TestLayoutPath_Length_veryShort(t *testing.T) {
	path := LayoutPath{Points: Points{Point{X: 5, Y: 4}}}
	assert.Equal(t, 0.0, path.Length())
}

func TestLayoutPaths_Draw(t *testing.T) {
	drawer := mocks.NewLayoutDrawer(t)

	p := LayoutPaths{
		LayoutPath{Points: Points{Point{X: 5.5, Y: 4.5}, Point{X: 8, Y: 4}, Point{X: 12, Y: 4}, Point{X: 14.5, Y: 4.5}}},
		LayoutPath{Points: Points{Point{X: 5.5, Y: 4.5}, Point{X: 8, Y: 5}, Point{X: 12, Y: 5}, Point{X: 14.5, Y: 4.5}}},
	}

	drawer.On("Path", mock.Anything, `class="path-line"`, `marker-end="url(#arrow)"`).Twice()

	p.Draw(drawer, 10)

	drawer.AssertExpectations(t)
}

func TestLayoutPaths_Length(t *testing.T) {

	t.Run("empty paths", func(t *testing.T) {
		var paths LayoutPaths
		length := paths.Length()
		assert.Equal(t, 0.0, length)
	})

	t.Run("single empty path", func(t *testing.T) {
		paths := LayoutPaths{LayoutPath{}}
		length := paths.Length()
		assert.Equal(t, 0.0, length)
	})

	t.Run("single path", func(t *testing.T) {
		path := LayoutPath{
			Points: Points{
				{X: 0, Y: 0},
				{X: 3, Y: 4},
			},
		}
		paths := LayoutPaths{path}
		length := paths.Length()

		assert.Equal(t, 5.0, length)
	})

	t.Run("multiple paths", func(t *testing.T) {
		paths := LayoutPaths{
			{
				Points: Points{{X: 0, Y: 0}, {X: 3, Y: 4}},
			},
			{
				Points: Points{{X: 0, Y: 0}, {X: 1, Y: 1}},
			},
		}

		length := paths.Length()

		assert.Equal(t, 5+math.Sqrt(2), length)
	})

}
