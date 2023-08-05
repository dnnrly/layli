package layli

import (
	"strings"
	"testing"

	"github.com/dnnrly/layli/mocks"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var pathTestConfig = Config{
	Nodes: ConfigNodes{
		ConfigNode{Id: "1"},
		ConfigNode{Id: "2"},
	},
	Spacing:    20,
	NodeWidth:  3,
	NodeHeight: 3,
	Margin:     2,
	Border:     2,
}

func TestLayout_AddPath_BetweenAdjacentNodes(t *testing.T) {
	finder := mocks.NewPathFinder(t)
	var gotStart dijkstra.Point
	var gotEnd dijkstra.Point

	finder.On("AddConnection", mock.Anything, mock.Anything, mock.Anything)
	finder.On("BestPath").Return([]dijkstra.Point{
		Point{X: 5.5, Y: 5.5},
		Point{X: 6, Y: 5},
		Point{X: 12, Y: 5},
		Point{X: 12.5, Y: 5.5},
	}, nil)

	l := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		gotStart = start
		gotEnd = end
		return finder
	}, pathTestConfig)

	require.NoError(t, l.AddPath("1", "2"))

	assert.Equal(t, Point{X: 5.5, Y: 5.5}, gotStart)
	assert.Equal(t, Point{X: 12.5, Y: 5.5}, gotEnd)

	assert.Len(t, l.Paths, 1)
	assert.Equal(t, LayoutPath{
		Points: []Point{
			{X: 5.5, Y: 5.5},
			{X: 6, Y: 5},
			{X: 12, Y: 5},
			{X: 12.5, Y: 5.5},
		},
	}, l.Paths[0])
}

// func TestLayout_AddPath_ErrorAddingPaths(t *testing.T) {
// 	finder := mocks.NewPathFinder(t)

// 	finder.On("AddConnection", mock.Anything, mock.Anything, mock.Anything)
// 	finder.On("BestPath").Return([]dijkstra.Point{}, errors.New("not paths"))

// 	l := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
// 		return finder
// 	}, pathTestConfig)

// 	require.Error(t, l.AddPath("1", "2"))
// }

func TestLayout_BuildVertexMap(t *testing.T) {
	finder := mocks.NewPathFinder(t)
	l := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		return finder
	}, pathTestConfig)
	vm := BuildVertexMap(l)

	assert.Equal(t, strings.ReplaceAll(
		`..................
		..................
		..xxxxxxxxxxxxxx..
		..xxxxxxxxxxxxxx..
		..xx.x.xxxx.x.xx..
		..xxx.xxxxxx.xxx..
		..xx.x.xxxx.x.xx..
		..xxxxxxxxxxxxxx..
		..xxxxxxxxxxxxxx..
		..................
		..................`, "	", ""), vm.String(), vm)
}

func TestLayout_BuildVertexMapWithPaths(t *testing.T) {
	finder := mocks.NewPathFinder(t)
	l := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		return finder
	}, pathTestConfig)
	l.Paths = append(l.Paths, LayoutPath{
		Points: Points{
			Point{X: 5.5, Y: 5},
			Point{X: 6, Y: 5},
			Point{X: 11, Y: 5},
			Point{X: 11.5, Y: 5},
		},
	})

	vm := BuildVertexMap(l)

	assert.Equal(t, strings.ReplaceAll(
		`..................
		..................
		..xxxxxxxxxxxxxx..
		..xxxxxxxxxxxxxx..
		..xx.x.xxxx.x.xx..
		..xxx........xxx..
		..xx.x.xxxx.x.xx..
		..xxxxxxxxxxxxxx..
		..xxxxxxxxxxxxxx..
		..................
		..................`, "	", ""), vm.String(), vm)
}
