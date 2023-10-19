package layli

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
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

	l, _ := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		gotStart = start
		gotEnd = end
		return finder
	}, &pathTestConfig)

	path, err := l.FindPath("1", "2")
	require.NoError(t, err)

	assert.Equal(t, Point{X: 5.5, Y: 5.5}, gotStart)
	assert.Equal(t, Point{X: 12.5, Y: 5.5}, gotEnd)

	assert.Equal(t, LayoutPath{
		Points: []Point{
			{X: 5.5, Y: 5.5},
			{X: 6, Y: 5},
			{X: 12, Y: 5},
			{X: 12.5, Y: 5.5},
		},
	}, *path)
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
	l, _ := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		return finder
	}, &pathTestConfig)
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
	l, _ := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		return finder
	}, &pathTestConfig)
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

func Test_selectPathStrategy(t *testing.T) {
	tests := []struct {
		name    string
		c       Config
		want    PathStrategy
		wantErr bool
	}{
		{name: "unknown generates error", c: Config{Strategy: "unknown"}, want: nil, wantErr: true},
		{name: "defaults to in-order", c: Config{}, want: findPathsInOrder, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectPathStrategy(&(tt.c))
			if (err != nil) != tt.wantErr {
				t.Errorf("selectPathStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t,
					reflect.ValueOf(tt.want).Pointer(), reflect.ValueOf(got).Pointer(),
					fmt.Sprintf("expected '%v' but got '%v'",
						runtime.FuncForPC(reflect.ValueOf(tt.want).Pointer()).Name(),
						runtime.FuncForPC(reflect.ValueOf(got).Pointer()).Name(),
					))
			}
		})
	}
}

func Test_findPathsInOrder_runsInOrder(t *testing.T) {
	var record []struct {
		from, to string
	}
	paths := LayoutPaths{}
	err := findPathsInOrder(
		ConfigEdges{
			{From: "a", To: "b"},
			{From: "1", To: "2"},
			{From: "2", To: "3"},
			{From: "r", To: "t"},
		},
		&paths,
		func(from, to string) (*LayoutPath, error) {
			record = append(record, struct {
				from string
				to   string
			}{from: from, to: to})
			return &LayoutPath{}, nil
		},
	)

	assert.NoError(t, err)
	assert.Len(t, paths, 4)
	assert.Equal(t, []struct {
		from string
		to   string
	}{{from: "a", to: "b"},
		{from: "1", to: "2"},
		{from: "2", to: "3"},
		{from: "r", to: "t"},
	}, record)
}

func Test_findPathsInOrder_passesErrorThrough(t *testing.T) {
	expectedErr := errors.New("an error")

	paths := LayoutPaths{}
	err := findPathsInOrder(
		ConfigEdges{
			{From: "a", To: "b"},
			{From: "1", To: "2"},
			{From: "2", To: "3"},
			{From: "r", To: "t"},
		},
		&paths,
		func(from, to string) (*LayoutPath, error) {
			return nil, expectedErr
		},
	)

	assert.ErrorIs(t, expectedErr, err)
}
