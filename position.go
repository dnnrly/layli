package layli

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/RyanCarrier/dijkstra"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Distance(to Point) float64 {
	a := (p.X - to.X)
	b := (p.Y - to.Y)
	return math.Sqrt(a*a + b*b)
}

func (p Point) String() string {
	return fmt.Sprintf("%.1f,%.1f", p.X, p.Y)
}

type Points []Point

func NewPointsFromBestPath(g Graph, path dijkstra.BestPath) Points {
	points := Points{}

	for _, id := range path.Path {
		p, err := g.GetMapped(id)
		if err != nil {
			panic(fmt.Sprintf("getting vertex %d from graph: %v", id, err))
		}

		parts := strings.Split(p, ",")
		x, _ := strconv.ParseFloat(parts[0], 64)
		y, _ := strconv.ParseFloat(parts[1], 64)

		points = append(points, Point{X: x, Y: y})
	}

	return points
}

func (p Points) Path(spacing int) string {
	path := fmt.Sprintf(
		"M %d %d",
		int(p[1].X)*spacing,
		int(p[1].Y)*spacing,
	)

	for i := 2; i < len(p)-1; i++ {
		path += fmt.Sprintf(
			" L %d %d",
			int(p[i].X)*spacing,
			int(p[i].Y)*spacing,
		)
	}

	return path
}

func (p Points) AddToGraph(g Graph) {
	for _, v := range p {
		g.AddMappedVertex(v.String())
	}
}
