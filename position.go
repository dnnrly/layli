package layli

import (
	"fmt"
	"math"

	"github.com/dnnrly/layli/pathfinder/dijkstra"
)

type Point struct {
	X float64
	Y float64
}

func PythagoreanDistance(from, to dijkstra.Point) int64 {
	x1, y1 := from.Coordinates()
	x2, y2 := to.Coordinates()

	a := (x1 - x2)
	b := (y1 - y2)

	return int64(math.Sqrt(a*a+b*b) * 100)
}

func (p Point) Distance(to Point) float64 {
	a := (p.X - to.X)
	b := (p.Y - to.Y)
	return math.Sqrt(a*a + b*b)
}

func (p Point) String() string {
	return fmt.Sprintf("%.1f,%.1f", p.X, p.Y)
}

func (p Point) Coordinates() (float64, float64) {
	return p.X, p.Y
}

type Points []Point

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
