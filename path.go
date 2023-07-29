package layli

import (
	"fmt"

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

	return vm
}

func (l *Layout) AddPath(from, to string) error {
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
		return fmt.Errorf("finding shortest path: %w", err)
	}

	path := LayoutPath{}
	for _, p := range points {
		x, y := p.Coordinates()
		path.Points = append(path.Points, Point{X: x, Y: y})
	}

	l.Paths = append(l.Paths, path)

	return nil
}
