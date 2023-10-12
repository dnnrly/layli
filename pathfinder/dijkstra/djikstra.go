package dijkstra

import (
	"container/heap"
	"errors"
	"math"
)

type CostFunction func(from, to Point) int64

type PathFinder struct {
	nodes map[Point]map[Point]int64
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
		nodes: make(map[Point]map[Point]int64),
		start: start,
		end:   end,
	}
}

func (pf *PathFinder) AddConnection(from Point, cost CostFunction, to ...Point) {
	if pf.nodes[from] == nil {
		pf.nodes[from] = make(map[Point]int64)
	}
	for _, dest := range to {
		pf.nodes[from][dest] = cost(from, dest)
		if pf.nodes[dest] == nil {
			pf.nodes[dest] = make(map[Point]int64)
		}
	}
}

func (pf *PathFinder) BestPath() ([]Point, error) {
	distance := make(map[Point]int64)
	previous := make(map[Point]Point)
	pq := make(PriorityQueue, 0)

	for node := range pf.nodes {
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

		for neighbor := range pf.nodes[current] {
			cost := distance[current] + pf.nodes[current][neighbor]

			if cost < distance[neighbor] {
				distance[neighbor] = cost
				previous[neighbor] = current
				heap.Push(&pq, &PriorityQueueItem{point: neighbor, cost: cost})
			}
		}
	}

	return nil, errors.New("no path found")
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
