package layli

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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