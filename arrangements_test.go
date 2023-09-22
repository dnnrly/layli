package layli

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertLeftOf(t *testing.T, left, right LayoutNode) {
	s := func(n LayoutNode) string {
		return fmt.Sprintf("L%dR%dT%dB%d", n.left, n.right, n.top, n.bottom)
	}
	assert.Less(t, left.left, right.left, fmt.Sprintf("node '%s' (%s) is not left of node '%s' (%s)", left.Id, s(left), right.Id, s(right)))
	assert.Less(t, left.right, right.right, fmt.Sprintf("node '%s' (%s) is not left of node '%s' (%s)", left.Id, s(left), right.Id, s(right)))
	assert.Less(t, left.right, right.left, fmt.Sprintf("node '%s' (%s) is not left of node '%s' (%s)", left.Id, s(left), right.Id, s(right)))
}

// func assertAbove(top, bottom LayoutNode) {
// 	assert.Less(t, top.bottom, bottom.top, fmt.Sprintf("node '%s' is not above node '%s'", top.Id, bottom.Id))
// }

func assertSameRow(t *testing.T, n1, n2 LayoutNode) {
	assert.Equal(t, n1.top, n2.top, fmt.Sprintf("node '%s' is not on the same row aw node '%s'", n1.Id, n2.Id))
}

func TestSelectArrangement(t *testing.T) {
	a := func(expected, actual LayoutArrangementFunc) {
		assert.Equal(t,
			reflect.ValueOf(expected).Pointer(), reflect.ValueOf(actual).Pointer(),
			fmt.Sprintf("expected '%v' but got '%v'",
				runtime.FuncForPC(reflect.ValueOf(expected).Pointer()).Name(),
				runtime.FuncForPC(reflect.ValueOf(actual).Pointer()).Name(),
			))
	}

	a(nil, selectArrangement(&Config{Layout: "unknown"}))
	a(LayoutFlowSquare, selectArrangement(&Config{}))
	a(LayoutFlowSquare, selectArrangement(&Config{Layout: "flow-square"}))
	a(LayoutTopologicalSort, selectArrangement(&Config{Layout: "topo-sort"}))
}

func TestLayoutFlowSquare(t *testing.T) {
	{
		l := LayoutFlowSquare(newConfig(2, 5, 3, 1, 1))

		require.Len(t, l, 2)
		assert.EqualValues(t, LayoutNode{"1", "", 5, 3, 2, 4, 2, 6}, l[0])
		assert.EqualValues(t, LayoutNode{"2", "", 5, 3, 2, 4, 9, 13}, l[1])
	}

	{
		l := LayoutFlowSquare(newConfig(4, 5, 3, 1, 1))

		require.Len(t, l, 4)
		assert.EqualValues(t, LayoutNode{"1", "", 5, 3, 2, 4, 2, 6}, l[0])
		assert.EqualValues(t, LayoutNode{"2", "", 5, 3, 2, 4, 9, 13}, l[1])
		assert.EqualValues(t, LayoutNode{"3", "", 5, 3, 7, 9, 2, 6}, l[2])
		assert.EqualValues(t, LayoutNode{"4", "", 5, 3, 7, 9, 9, 13}, l[3])
	}

	{
		l := LayoutFlowSquare(newConfig(4, 5, 3, 2, 1))

		require.Len(t, l, 4)
		assert.EqualValues(t, LayoutNode{"1", "", 5, 3, 3, 5, 3, 7}, l[0])
		assert.EqualValues(t, LayoutNode{"2", "", 5, 3, 3, 5, 12, 16}, l[1])
		assert.EqualValues(t, LayoutNode{"3", "", 5, 3, 10, 12, 3, 7}, l[2])
		assert.EqualValues(t, LayoutNode{"4", "", 5, 3, 10, 12, 12, 16}, l[3])
	}

	{
		l := LayoutFlowSquare(newConfig(8, 5, 3, 2, 1))

		require.Len(t, l, 8)
		assert.EqualValues(t, LayoutNode{"1", "", 5, 3, 3, 5, 3, 7}, l[0])
		assert.EqualValues(t, LayoutNode{"3", "", 5, 3, 3, 5, 21, 25}, l[2])
		assert.EqualValues(t, LayoutNode{"4", "", 5, 3, 10, 12, 3, 7}, l[3])
		assert.EqualValues(t, LayoutNode{"6", "", 5, 3, 10, 12, 21, 25}, l[5])
		assert.EqualValues(t, LayoutNode{"8", "", 5, 3, 17, 19, 12, 16}, l[7])
	}

	{
		l := LayoutFlowSquare(newConfig(4, 5, 4, 2, 2))

		require.Len(t, l, 4)
		assert.EqualValues(t, LayoutNode{"1", "", 5, 4, 3, 6, 3, 7}, l[0])
		assert.EqualValues(t, LayoutNode{"4", "", 5, 4, 11, 14, 12, 16}, l[3])
	}
}

func TestLayoutTopologicalSort_simpleLine(t *testing.T) {
	nodes := LayoutTopologicalSort(&Config{
		Nodes: ConfigNodes{ConfigNode{Id: "1"}, ConfigNode{Id: "2"}, ConfigNode{Id: "3"}},
		Edges: ConfigEdges{
			ConfigEdge{From: "1", To: "3"},
			ConfigEdge{From: "3", To: "2"},
		},
		Border: 1, Spacing: 1,
		NodeWidth: 1, NodeHeight: 1, Margin: 1,
	})

	assert.Len(t, nodes, 3)

	assertLeftOf(t, *nodes.ByID("1"), *nodes.ByID("3"))
	assertLeftOf(t, *nodes.ByID("3"), *nodes.ByID("2"))
}
