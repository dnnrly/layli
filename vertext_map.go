package layli

import "strings"

type VertexMap struct {
	points [][]bool
	width  int
	height int
}

func NewVertexMap(x, y int) VertexMap {
	v := VertexMap{
		points: make([][]bool, x),
		width:  x,
		height: y,
	}

	for column := range v.points {
		v.points[column] = make([]bool, y)

	}

	return v
}

func (v *VertexMap) Get(x, y int) bool {
	return v.points[x][y]
}

func (v *VertexMap) Set(x, y int, val bool) {
	v.points[x][y] = val
}

func (v *VertexMap) CountAvailable(available bool) int {
	count := 0
	for x := 0; x < v.width; x++ {
		for y := 0; y < v.height; y++ {
			if v.points[x][y] == available {
				count++
			}
		}
	}

	return count
}

func (v VertexMap) String() string {
	s := make([]string, v.height)

	for x := 0; x < v.width; x++ {
		for y := 0; y < v.height; y++ {
			if v.Get(x, y) {
				s[y] += "x"
			} else {
				s[y] += "."
			}
		}
	}

	return strings.Join(s, "\n")
}

type VertexMapper func(x, y int) bool

func (v *VertexMap) MapAvailable(m VertexMapper) {
	for x := 0; x < v.width; x++ {
		for y := 0; y < v.height; y++ {
			v.points[x][y] = m(x, y)
		}
	}
}
