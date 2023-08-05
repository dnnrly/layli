package layli

import (
	"fmt"
	"math"

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
		// return fmt.Errorf("cannot find a path between %w", err)
		return fmt.Errorf("cannot find a path between %s and %s", from, to)
	}

	path := LayoutPath{}
	for _, p := range points {
		x, y := p.Coordinates()
		path.Points = append(path.Points, Point{X: x, Y: y})
	}

	l.Paths = append(l.Paths, path)

	return nil
}
