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

func TestFindPath_returnsCorrectErrors(t *testing.T) {
	finder := mocks.NewPathFinder(t)
	expectedErr := errors.New("some error")

	finder.On("AddConnection", mock.Anything, mock.Anything, mock.Anything)
	finder.On("BestPath").Return([]dijkstra.Point{}, expectedErr)

	l, _ := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
		return finder
	}, &pathTestConfig)

	_, err := l.FindPath("1", "2")
	require.ErrorIs(t, err, expectedErr)
}

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
		{name: "unknown generates error", c: Config{Path: ConfigPath{Strategy: "unknown"}}, want: nil, wantErr: true},
		{name: "defaults to in-order", c: Config{}, want: findPathsInOrder, wantErr: false},
		{name: "selects in-order", c: Config{Path: ConfigPath{Strategy: "in-order"}}, want: findPathsInOrder, wantErr: false},
		{name: "selects random", c: Config{Path: ConfigPath{Strategy: "random"}}, want: findPathsRandomlyByOrder, wantErr: false},
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
		Config{
			Edges: ConfigEdges{
				{From: "a", To: "b"},
				{From: "1", To: "2"},
				{From: "2", To: "3"},
				{From: "r", To: "t"},
			},
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

func Test_findPathsInOrder_addsIdCorrectly(t *testing.T) {
	var record []struct {
		from, to string
	}
	paths := LayoutPaths{}
	err := findPathsInOrder(
		Config{
			Edges: ConfigEdges{
				{From: "a", To: "b", ID: "1", Class: "a class"},
				{From: "1", To: "2", ID: "2", Style: "some style"},
				{From: "2", To: "3", ID: "3"},
				{From: "r", To: "t", ID: "4"},
			},
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
	assert.Equal(t, "1", paths[0].ID)
	assert.Equal(t, "a class", paths[0].Class)
	assert.Equal(t, "a", paths[0].From)
	assert.Equal(t, "b", paths[0].To)
	assert.Empty(t, paths[0].Style)

	assert.Equal(t, "2", paths[1].ID)
	assert.Empty(t, paths[1].Class)
	assert.Equal(t, "some style", paths[1].Style)

	assert.Equal(t, "3", paths[2].ID)
	assert.Empty(t, paths[2].Class)
	assert.Empty(t, paths[2].Style)
	assert.Equal(t, "2", paths[2].From)
	assert.Equal(t, "3", paths[2].To)
}

func Test_findPathsInOrder_passesErrorThrough(t *testing.T) {
	expectedErr := errors.New("an error")

	paths := LayoutPaths{}
	err := findPathsInOrder(
		Config{
			Edges: ConfigEdges{
				{From: "a", To: "b"},
				{From: "1", To: "2"},
				{From: "2", To: "3"},
				{From: "r", To: "t"},
			},
		},
		&paths,
		func(from, to string) (*LayoutPath, error) {
			return nil, expectedErr
		},
	)

	assert.ErrorIs(t, expectedErr, err)
}

func Test_findPathsRandomly_addsAllPaths(t *testing.T) {
	config := Config{
		Path: ConfigPath{
			Attempts: 5,
		},
		Edges: ConfigEdges{
			{From: "a", To: "b"},
			{From: "1", To: "2"},
			{From: "2", To: "3"},
			{From: "r", To: "t"},
		},
	}
	records := []ConfigEdges{}
	subStrat := func(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error {
		records = append(records, config.Edges)
		return nil
	}
	err := findPathsRandomly(subStrat)(config, &LayoutPaths{}, func(from, to string) (*LayoutPath, error) { return &LayoutPath{}, nil })

	assert.NoError(t, err)
	assert.Len(t, records, 5)
	for _, r := range records {
		assert.Len(t, r, 4)
		assert.Contains(t, r, ConfigEdge{From: "a", To: "b"})
		assert.Contains(t, r, ConfigEdge{From: "1", To: "2"})
		assert.Contains(t, r, ConfigEdge{From: "2", To: "3"})
		assert.Contains(t, r, ConfigEdge{From: "r", To: "t"})
	}
}

func Test_findPathsRandomly_orderChanges(t *testing.T) {
	config := Config{
		Path: ConfigPath{
			Attempts: 5,
		},
		Edges: ConfigEdges{
			{From: "a", To: "b"}, {From: "1", To: "2"}, {From: "2", To: "3"}, {From: "r", To: "t"},
			{From: "a", To: "d"}, {From: "1", To: "5"}, {From: "2", To: "99"}, {From: "r", To: "d"},
			{From: "5", To: "b"}, {From: "4", To: "2"}, {From: "7", To: "3"}, {From: "i", To: "t"},
		},
	}
	last := ConfigEdges{}
	last = append(last, config.Edges...)
	subStrat := func(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error {
		assert.NotEqual(t, last, config.Edges)
		last = append(last, config.Edges...)
		return nil
	}
	err := findPathsRandomly(subStrat)(config, &LayoutPaths{}, func(from, to string) (*LayoutPath, error) { return &LayoutPath{}, nil })

	assert.NoError(t, err)
}

func Test_findPathsRandomly_selectsShortestPath(t *testing.T) {
	config := Config{Path: ConfigPath{
		Attempts: 5},
		Edges: ConfigEdges{{From: "a", To: "b"}, {From: "1", To: "2"}, {From: "2", To: "3"}, {From: "r", To: "t"}},
	}
	subStrat := func(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error {
		return nil
	}
	err := findPathsRandomly(subStrat)(config, &LayoutPaths{}, func(from, to string) (*LayoutPath, error) { return &LayoutPath{}, nil })

	assert.NoError(t, err)
}

func Test_findPathsRandomly_passesErrorThrough(t *testing.T) {
	expectedErr := errors.New("an error")
	count := 0

	paths := LayoutPaths{}
	err := findPathsRandomly(findPathsInOrder)(
		Config{
			Path: ConfigPath{
				Attempts: 5,
			},
			Edges: ConfigEdges{
				{From: "a", To: "b"},
				{From: "1", To: "2"},
				{From: "2", To: "3"},
				{From: "r", To: "t"},
			},
		},
		&paths,
		func(from, to string) (*LayoutPath, error) {
			count++
			return nil, expectedErr
		},
	)

	assert.ErrorIs(t, expectedErr, err)
	assert.Equal(t, 1, count)
}

func Test_findPathsRandomly_retriesWhenStrugglingToFindPath(t *testing.T) {
	count := 0

	paths := LayoutPaths{}
	err := findPathsRandomly(findPathsInOrder)(
		Config{
			Path: ConfigPath{
				Attempts: 5,
			},
			Edges: ConfigEdges{
				{From: "a", To: "b"},
			},
		},
		&paths,
		func(from, to string) (*LayoutPath, error) {
			count++
			if count%2 == 0 {
				return nil, dijkstra.ErrNotFound
			}
			return &LayoutPath{
				Points: Points{
					Point{X: 1, Y: 0},
					Point{X: 1, Y: 2.0},
				},
			}, nil
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, 2.0, paths.Length())
}

func Test_findPathsRandomly_eventuallyGivesUp(t *testing.T) {
	paths := LayoutPaths{}
	err := findPathsRandomly(findPathsInOrder)(
		Config{
			Path: ConfigPath{
				Attempts: 2,
			},
			Edges: ConfigEdges{
				{From: "a", To: "b"},
				{From: "1", To: "2"},
				{From: "2", To: "3"},
				{From: "r", To: "t"},
			},
		},
		&paths,
		func(from, to string) (*LayoutPath, error) {
			return nil, dijkstra.ErrNotFound
		},
	)

	assert.ErrorIs(t, err, dijkstra.ErrNotFound, "got error: %v", err)
}
