package layli

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/dnnrly/layli/pathfinder/dijkstra"
)

type PathFinder interface {
	AddConnection(from dijkstra.Point, cost dijkstra.CostFunction, to ...dijkstra.Point)
	BestPath() ([]dijkstra.Point, error)
}

type CreateFinder func(start, end dijkstra.Point) PathFinder

func BuildVertexMap(l *Layout) VertexMap {
	vm := NewVertexMap(l.LayoutWidth(), l.LayoutHeight())
	vm.MapUnset(l.InsideAny)
	vm.MapOr(l.IsAnyPort)
	vm.Map(func(x, y int, current bool) bool { return x >= l.layoutBorder && current })
	vm.Map(func(x, y int, current bool) bool { return y >= l.layoutBorder && current })
	vm.Map(func(x, y int, current bool) bool { return x < l.LayoutWidth()-l.layoutBorder && current })
	vm.Map(func(x, y int, current bool) bool { return y < l.LayoutHeight()-l.layoutBorder && current })

	for _, v := range l.Paths {
		for p := 1; p < len(v.Points)-2; p++ {
			start := v.Points[p]
			end := v.Points[p+1]

			x := math.Min(start.X, end.X)
			y := math.Min(start.Y, end.Y)

			xMax := math.Max(start.X, end.X)
			yMax := math.Max(start.Y, end.Y)

			vm.Set(int(x), int(y), false)
			for x != xMax || y != yMax {
				if x != xMax {
					x += 1
				}
				if y != yMax {
					y += 1
				}

				vm.Set(int(x), int(y), false)
			}
		}
	}

	return vm
}

func (l *Layout) FindPath(from, to string) (*LayoutPath, error) {
	nFrom := l.Nodes.ByID(from)
	nTo := l.Nodes.ByID(to)

	finder := l.CreateFinder(
		nFrom.GetCentre(),
		nTo.GetCentre(),
	)

	vm := BuildVertexMap(l)

	arcs := vm.GetArcs()
	for _, a := range arcs {
		finder.AddConnection(a.From, PythagoreanDistance, a.To)
	}

	{
		// Add "from" paths
		centre := nFrom.GetCentre()
		ports := nFrom.GetPorts()
		for _, to := range ports {
			finder.AddConnection(centre, PythagoreanDistance, to)
		}
	}

	{
		// Add "to" paths
		centre := nTo.GetCentre()
		ports := nTo.GetPorts()
		for _, from := range ports {
			finder.AddConnection(
				from,
				PythagoreanDistance,
				centre,
			)
		}
	}

	points, err := finder.BestPath()
	if err != nil {
		return nil, fmt.Errorf("finding path between %s and %s: %w", from, to, err)
	}

	path := LayoutPath{}
	for _, p := range points {
		x, y := p.Coordinates()
		path.Points = append(path.Points, Point{X: x, Y: y})
	}

	return &path, nil
}

type PathStrategy func(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error

func selectPathStrategy(c *Config) (PathStrategy, error) {
	switch c.Path.Strategy {
	// Shortest first
	case "random":
		return findPathsRandomlyByOrder, nil
	case "in-order":
		fallthrough
	case "":
		return findPathsInOrder, nil
	default:
		return nil, errors.New("cannot find path strategy " + c.Path.Strategy)
	}
}

func findPathsInOrder(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error {
	for _, p := range config.Edges {
		path, err := find(p.From, p.To)
		if err != nil {
			return err
		}

		*paths = append(*paths, *path)
	}
	return nil
}

func findPathsRandomlyByOrder(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error {
	return findPathsRandomly(findPathsInOrder)(config, paths, find)
}

func findPathsRandomly(subStrategy PathStrategy) PathStrategy {
	return func(config Config, paths *LayoutPaths, find func(from, to string) (*LayoutPath, error)) error {
		shortest := LayoutPaths{LayoutPath{Points: []Point{
			{X: 0.0, Y: 0.0},
			{X: math.MaxFloat64, Y: math.MaxFloat64},
		}}}

		gotPath := false

		for count := 0; count < config.Path.Attempts; count++ {
			rand.Shuffle(len(config.Edges), func(i, j int) { config.Edges[i], config.Edges[j] = config.Edges[j], config.Edges[i] })
			err := subStrategy(config, paths, find)

			if err == nil {
				if paths.Length() < shortest.Length() {
					shortest = *paths
				}
				gotPath = true
			} else {
				if !errors.Is(err, dijkstra.ErrNotFound) {
					return err
				}
			}
		}

		if !gotPath {
			return dijkstra.ErrNotFound
		}

		*paths = shortest

		return nil
	}
}
