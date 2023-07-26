package djikstra

import (
	"fmt"
	"strconv"
	"strings"

	dj "github.com/RyanCarrier/dijkstra"
)

type PathFinder struct {
	graph dj.Graph
	start int
	end   int
}

func NewPathFinder(start, end Point) *PathFinder {
	pf := &PathFinder{
		graph: *dj.NewGraph(),
	}

	pf.start = pf.graph.AddMappedVertex(start.String())
	pf.end = pf.graph.AddMappedVertex(end.String())

	return pf
}

type CostFunction func(from, to Point) int64

func (pf *PathFinder) AddConnection(from Point, cost CostFunction, to Point) {
	pf.graph.AddMappedVertex(from.String())
	pf.graph.AddMappedVertex(to.String())

	pf.graph.AddMappedArc(from.String(), to.String(), cost(from, to))
}

func (pf *PathFinder) BestPath() ([]Point, error) {
	found, err := pf.graph.Shortest(pf.start, pf.end)
	if err != nil {
		return nil, fmt.Errorf("calculating shortest path: %w", err)
	}

	path := []Point{}

	for _, id := range found.Path {
		p, _ := pf.graph.GetMapped(id)

		parts := strings.Split(p, ",")
		x, _ := strconv.ParseFloat(parts[0], 64)
		y, _ := strconv.ParseFloat(parts[1], 64)

		path = append(path, Point(&coordinate{x: x, y: y}))
	}

	return path, nil
}
