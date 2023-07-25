package djikstra

import (
	"strconv"
	"strings"

	dj "github.com/RyanCarrier/dijkstra"
	"github.com/dnnrly/layli"
)

type PathFinder struct {
	graph dj.Graph
	start int
	end   int
}

func NewPathFinder(start, end layli.Point) *PathFinder {
	pf := &PathFinder{
		graph: *dj.NewGraph(),
	}

	pf.start = pf.graph.AddMappedVertex(start.String())
	pf.end = pf.graph.AddMappedVertex(end.String())

	return pf
}

type CostFunction func(from, to layli.Point) int64

func (pf *PathFinder) AddConnection(from layli.Point, cost CostFunction, to layli.Point) {
	pf.graph.AddMappedVertex(from.String())
	pf.graph.AddMappedVertex(to.String())

	pf.graph.AddMappedArc(from.String(), to.String(), cost(from, to))
}

func (pf *PathFinder) BestPath() layli.LayoutPath {
	found, _ := pf.graph.Shortest(pf.start, pf.end)
	path := layli.LayoutPath{}

	for _, id := range found.Path {
		p, _ := pf.graph.GetMapped(id)

		parts := strings.Split(p, ",")
		x, _ := strconv.ParseFloat(parts[0], 64)
		y, _ := strconv.ParseFloat(parts[1], 64)

		path.Points = append(path.Points, layli.Point{X: x, Y: y})
	}

	return path
}
