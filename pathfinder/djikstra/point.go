package djikstra

import "fmt"

type Point interface {
	String() string
	Coordinates() (float64, float64)
}

type coordinate struct {
	x, y float64
}

func (c coordinate) String() string {
	return fmt.Sprintf("%.1f,%.1f", c.x, c.y)
}

func (c coordinate) Coordinates() (float64, float64) {
	return c.x, c.y
}
