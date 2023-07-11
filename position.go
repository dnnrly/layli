package layli

import (
	"fmt"
)

type Point struct {
	X float64
	Y float64
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
