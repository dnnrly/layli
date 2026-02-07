package dijkstra

import (
	"container/heap"
	"math"
)

// HeuristicFunction estimates the cost from a point to the goal
type HeuristicFunction func(from, to Point) int64

// EuclideanDistance returns the straight-line distance between two points
func EuclideanDistance(from, to Point) int64 {
	fromX, fromY := from.Coordinates()
	toX, toY := to.Coordinates()
	dx := fromX - toX
	dy := fromY - toY
	return int64(math.Sqrt(dx*dx + dy*dy))
}

// ManhattanDistance returns the Manhattan distance between two points
func ManhattanDistance(from, to Point) int64 {
	fromX, fromY := from.Coordinates()
	toX, toY := to.Coordinates()
	return int64(math.Abs(fromX-toX) + math.Abs(fromY-toY))
}

// AStarPathFinder implements the A* algorithm
type AStarPathFinder struct {
	*PathFinder
	heuristic HeuristicFunction
	goal      Point
}

// AStarPriorityQueueItem represents an item in the A* priority queue
type AStarPriorityQueueItem struct {
	point     Point
	gCost     int64 // Cost from start to this point
	fCost     int64 // Estimated total cost (gCost + heuristic)
	index     int   // For heap interface
}

// AStarPriorityQueue implements the heap interface for A*
type AStarPriorityQueue []*AStarPriorityQueueItem

func (pq AStarPriorityQueue) Len() int { return len(pq) }

func (pq AStarPriorityQueue) Less(i, j int) bool {
	return pq[i].fCost < pq[j].fCost
}

func (pq AStarPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *AStarPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*AStarPriorityQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *AStarPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// NewAStarPathFinder creates a new A* pathfinder with the specified heuristic
func NewAStarPathFinder(start, end Point, heuristic HeuristicFunction) *AStarPathFinder {
	return &AStarPathFinder{
		PathFinder: NewPathFinder(start, end),
		heuristic:  heuristic,
		goal:       end,
	}
}

// NewAStarPathFinderWithEuclidean creates a new A* pathfinder with Euclidean distance heuristic
func NewAStarPathFinderWithEuclidean(start, end Point) *AStarPathFinder {
	return NewAStarPathFinder(start, end, EuclideanDistance)
}

// NewAStarPathFinderWithManhattan creates a new A* pathfinder with Manhattan distance heuristic
func NewAStarPathFinderWithManhattan(start, end Point) *AStarPathFinder {
	return NewAStarPathFinder(start, end, ManhattanDistance)
}

// BestPath implements the A* algorithm to find the shortest path
func (pf *AStarPathFinder) BestPath() ([]Point, error) {
	// gScore represents the cost of the cheapest path from start to each point
	gScore := make(map[Point]int64)
	// fScore represents gScore + heuristic estimate to goal
	fScore := make(map[Point]int64)
	// cameFrom reconstructs the path
	cameFrom := make(map[Point]Point)
	// openSet is the priority queue of nodes to evaluate
	openSet := make(AStarPriorityQueue, 0)
	// closedSet contains nodes already evaluated
	closedSet := make(map[Point]bool)

	// Initialize scores
	for _, node := range pf.nodes.Keys() {
		gScore[node] = math.MaxInt64
		fScore[node] = math.MaxInt64
	}

	// Start node
	gScore[pf.start] = 0
	fScore[pf.start] = pf.heuristic(pf.start, pf.goal)
	
	// Add start to open set
	startItem := &AStarPriorityQueueItem{
		point: pf.start,
		gCost: 0,
		fCost: fScore[pf.start],
	}
	heap.Push(&openSet, startItem)

	for len(openSet) > 0 {
		// Get the node with lowest fScore
		current := heap.Pop(&openSet).(*AStarPriorityQueueItem).point

		// Check if we reached the goal
		if current == pf.end {
			return pf.reconstructPath(cameFrom), nil
		}

		// Mark as evaluated
		closedSet[current] = true

		// Check all neighbors
		for _, neighborPoint := range (*get(pf.nodes, current)).Keys() {
			if closedSet[neighborPoint] {
				continue // Skip already evaluated nodes
			}

			// Calculate tentative gScore
			destMapping := get(pf.nodes, current)
			tentativeGScore := gScore[current] + get(destMapping, neighborPoint)

			// If this path to neighbor is better than any previous one
			if tentativeGScore < gScore[neighborPoint] {
				cameFrom[neighborPoint] = current
				gScore[neighborPoint] = tentativeGScore
				fScore[neighborPoint] = tentativeGScore + pf.heuristic(neighborPoint, pf.goal)

				// Check if neighbor is in open set
				var found bool
				for _, item := range openSet {
					if item.point == neighborPoint {
						// Update existing item
						item.gCost = tentativeGScore
						item.fCost = fScore[neighborPoint]
						heap.Fix(&openSet, item.index)
						found = true
						break
					}
				}

				if !found {
					// Add to open set
					neighborItem := &AStarPriorityQueueItem{
						point: neighborPoint,
						gCost: tentativeGScore,
						fCost: fScore[neighborPoint],
					}
					heap.Push(&openSet, neighborItem)
				}
			}
		}
	}

	return nil, ErrNotFound
}

// reconstructPath builds the path from start to end using the cameFrom map
func (pf *AStarPathFinder) reconstructPath(cameFrom map[Point]Point) []Point {
	var path []Point
	current := pf.end

	for current != nil {
		path = append([]Point{current}, path...)
		current = cameFrom[current]
	}

	return path
}
