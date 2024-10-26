package dijkstra

import (
	"container/heap"
	"errors"
	"math"

	"github.com/elliotchance/orderedmap/v2"
)

var (
	ErrNotFound = errors.New("no path found")
)

type CostFunction func(from, to Point) int64

type costMapping *orderedmap.OrderedMap[Point, int64]

func newCostMapping() costMapping {
	return orderedmap.NewOrderedMap[Point, int64]()
}

func get[K comparable, V any](m *orderedmap.OrderedMap[K, V], k K) V {
	var none V

	if m == nil {
		return none
	}

	v, ok := (*m).Get(k)
	if !ok {
		return none
	}

	return v
}

func set[K comparable, V any](m *orderedmap.OrderedMap[K, V], k K, v V) {
	if m == nil {
		return
	}
	(*m).Set(k, v)
}

type PathFinder struct {
	nodes *orderedmap.OrderedMap[Point, costMapping]
	start Point
	end   Point
}

type PriorityQueueItem struct {
	point Point
	cost  int64
}

type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PriorityQueueItem)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func NewPathFinder(start, end Point) *PathFinder {
	return &PathFinder{
		nodes: orderedmap.NewOrderedMap[Point, costMapping](),
		start: start,
		end:   end,
	}
}

func (pf *PathFinder) AddConnection(from Point, cost CostFunction, to ...Point) {
	if get(pf.nodes, from) == nil {
		set(pf.nodes, from, costMapping(orderedmap.NewOrderedMap[Point, int64]()))
	}
	for _, dest := range to {
		fromMapping := get(pf.nodes, from)
		set(fromMapping, dest, cost(from, dest))
		if get(pf.nodes, dest) == nil {
			set(pf.nodes, dest, newCostMapping())
		}
	}
}

func (pf *PathFinder) BestPath() ([]Point, error) {
	distance := make(map[Point]int64)
	previous := make(map[Point]Point)
	pq := make(PriorityQueue, 0)

	for _, node := range pf.nodes.Keys() {
		distance[node] = math.MaxInt64
	}
	distance[pf.start] = 0
	heap.Push(&pq, &PriorityQueueItem{point: pf.start, cost: 0})

	for len(pq) > 0 {
		currentItem := heap.Pop(&pq).(*PriorityQueueItem)
		current := currentItem.point

		if current == pf.end {
			return pf.reconstructPath(previous), nil
		}

		for _, neighborPoint := range (*get(pf.nodes, current)).Keys() {
			destMapping := get(pf.nodes, current)
			cost := distance[current] + get(destMapping, neighborPoint)

			if cost < distance[neighborPoint] {
				distance[neighborPoint] = cost
				previous[neighborPoint] = current
				heap.Push(&pq, &PriorityQueueItem{point: neighborPoint, cost: cost})
			}
		}
	}

	return nil, ErrNotFound
}

func (pf *PathFinder) reconstructPath(previous map[Point]Point) []Point {
	var path []Point
	current := pf.end

	for current != nil {
		path = append([]Point{current}, path...)
		current = previous[current]
	}

	return path
}
