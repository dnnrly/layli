package layli

import (
	"fmt"
	"math"
	"strings"
)

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

func (v *VertexMap) MapSet(m VertexMapper) {
	for x := 0; x < v.width; x++ {
		for y := 0; y < v.height; y++ {
			v.points[x][y] = m(x, y)
		}
	}
}

func (v *VertexMap) MapUnset(m VertexMapper) {
	v.MapSet(func(x, y int) bool {
		return !m(x, y)
	})
}

func (v *VertexMap) MapOr(m VertexMapper) {
	v.MapSet(func(x, y int) bool {
		return m(x, y) || v.points[x][y]
	})
}

func (v *VertexMap) MapAnd(m VertexMapper) {
	v.MapSet(func(x, y int) bool {
		return m(x, y) && v.points[x][y]
	})
}

func (v VertexMap) GetVertexPoints() Points {
	p := Points{}

	for x := 0; x < v.width; x++ {
		for y := 0; y < v.height; y++ {
			if v.points[x][y] {
				p = append(p, Point{X: float64(x), Y: float64(y)})
			}
		}
	}

	return p
}

func (vm VertexMap) GetArcs() Arcs {
	arcs := Arcs{}

	points := vm.GetVertexPoints()

	for _, from := range points {
		for _, to := range points {
			if from != to {
				if from.X == to.X {
					ok := true
					for p := int(math.Min(from.Y, to.Y)); p <= int(math.Max(from.Y, to.Y)); p++ {
						ok = ok && vm.points[int(from.X)][p]
					}

					if ok {
						arcs.Add(from, to, 1)
					}
				}

				if from.Y == to.Y {
					ok := true
					for p := int(math.Min(from.X, to.X)); p <= int(math.Max(from.X, to.X)); p++ {
						ok = ok && vm.points[p][int(from.Y)]
					}

					if ok {
						arcs.Add(from, to, 1)
					}
				}
			}
		}
	}

	return arcs
}

type Arc struct {
	From     Point
	To       Point
	Distance int
}

type Arcs []Arc

func (all *Arcs) Add(from Point, to Point, distance int) {
	*all = append(*all, Arc{
		From:     from,
		To:       to,
		Distance: distance,
	})
}

func (all Arcs) Exists(from Point, to Point) bool {
	for _, v := range all {
		if v.From == from && v.To == to {
			return true
		}
	}

	return false
}

func (all Arcs) AddToGraph(g Graph) {
	for _, v := range all {
		g.AddMappedArc(v.From.String(), v.To.String(), int64(v.Distance))
	}
}

func (all Arcs) String() string {
	str := []string{}

	for _, v := range all {
		str = append(str, fmt.Sprintf("%s-%s-%d", v.From, v.To, v.Distance))
	}

	return strings.Join(str, "\n")
}
