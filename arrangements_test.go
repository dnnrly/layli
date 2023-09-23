package layli

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func s(n LayoutNode) string {
	return fmt.Sprintf("L%d R%d T%d B%d", n.left, n.right, n.top, n.bottom)
}

func assertLeftOf(t *testing.T, left, right LayoutNode) {
	assert.Less(t, left.right, right.left, fmt.Sprintf("node '%s' (%s) is not left of node '%s' (%s)", left.Id, s(left), right.Id, s(right)))
}

func assertAbove(t *testing.T, upper, lower LayoutNode) {
	assert.Less(t, upper.bottom, lower.top, fmt.Sprintf("node '%s' (%s) is not above node '%s' (%s)", upper.Id, s(upper), lower.Id, s(lower)))
}

func assertSameRow(t *testing.T, n1, n2 LayoutNode) {
	assert.Equal(t, n1.top, n2.top, fmt.Sprintf("node '%s' (%s) is not on the same row as node '%s' (%s)", n1.Id, s(n1), n2.Id, s(n2)))
}

func assertSameColumn(t *testing.T, n1, n2 LayoutNode) {
	assert.Equal(t, n1.left, n2.left, fmt.Sprintf("node '%s' (%s) is not on the same column as node '%s' (%s)", n1.Id, s(n1), n2.Id, s(n2)))
}

func TestSelectArrangement(t *testing.T) {
	a := func(expected LayoutArrangementFunc, config Config) {
		actual, err := selectArrangement(&config)
		assert.NoError(t, err)
		assert.Equal(t,
			reflect.ValueOf(expected).Pointer(), reflect.ValueOf(actual).Pointer(),
			fmt.Sprintf("expected '%v' but got '%v'",
				runtime.FuncForPC(reflect.ValueOf(expected).Pointer()).Name(),
				runtime.FuncForPC(reflect.ValueOf(actual).Pointer()).Name(),
			))
	}

	a(LayoutFlowSquare, Config{})
	a(LayoutFlowSquare, Config{Layout: "flow-square"})
	a(LayoutTopologicalSort, Config{Layout: "topo-sort"})

	actual, err := selectArrangement(&Config{Layout: "unknown"})
	assert.Error(t, err)
	assert.Nil(t, actual)
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

	assertSameRow(t, *nodes.ByID("1"), *nodes.ByID("2"))
	assertSameRow(t, *nodes.ByID("2"), *nodes.ByID("3"))
}

func TestLayoutTarjan(t *testing.T) {
	nodes := LayoutTarjan(&Config{
		Nodes: ConfigNodes{ConfigNode{Id: "1"}, ConfigNode{Id: "2"}, ConfigNode{Id: "3"}, ConfigNode{Id: "4"}, ConfigNode{Id: "5"}},
		Edges: ConfigEdges{
			ConfigEdge{From: "1", To: "2"},
			ConfigEdge{From: "2", To: "3"},
			ConfigEdge{From: "3", To: "4"},
			ConfigEdge{From: "4", To: "5"},
			ConfigEdge{From: "3", To: "5"},
			ConfigEdge{From: "5", To: "4"},
		},
		Border: 1, Spacing: 1,
		NodeWidth: 1, NodeHeight: 1, Margin: 1,
	})

	assertLeftOf(t, *nodes.ByID("1"), *nodes.ByID("2"))
	assertLeftOf(t, *nodes.ByID("2"), *nodes.ByID("3"))
	assertLeftOf(t, *nodes.ByID("3"), *nodes.ByID("4"))

	assertAbove(t, *nodes.ByID("4"), *nodes.ByID("5"))

	assertSameRow(t, *nodes.ByID("1"), *nodes.ByID("2"))
	assertSameRow(t, *nodes.ByID("1"), *nodes.ByID("3"))
	assertSameRow(t, *nodes.ByID("1"), *nodes.ByID("4"))

	assertSameColumn(t, *nodes.ByID("4"), *nodes.ByID("5"))
}
