package dijkstra

import (
	"container/heap"
	"math"
)

// BidirectionalPathFinder implements Dijkstra's algorithm from both directions
type BidirectionalPathFinder struct {
	*PathFinder
}

// BidirectionalSearchState holds the search state for one direction
type BidirectionalSearchState struct {
	distance map[Point]int64
	previous map[Point]Point
	pq       PriorityQueue
	visited  map[Point]bool
	frontier map[Point]bool
}

// NewBidirectionalPathFinder creates a new bidirectional Dijkstra pathfinder
func NewBidirectionalPathFinder(start, end Point) *BidirectionalPathFinder {
	return &BidirectionalPathFinder{
		PathFinder: NewPathFinder(start, end),
	}
}

// BestPath implements bidirectional Dijkstra's algorithm
func (pf *BidirectionalPathFinder) BestPath() ([]Point, error) {
	// Handle the case where start and end are the same
	if pf.start == pf.end {
		return []Point{pf.start}, nil
	}

	// Initialize forward search (from start)
	forward := &BidirectionalSearchState{
		distance: make(map[Point]int64),
		previous: make(map[Point]Point),
		pq:       make(PriorityQueue, 0),
		visited:  make(map[Point]bool),
		frontier: make(map[Point]bool),
	}

	// Initialize backward search (from end)
	backward := &BidirectionalSearchState{
		distance: make(map[Point]int64),
		previous: make(map[Point]Point),
		pq:       make(PriorityQueue, 0),
		visited:  make(map[Point]bool),
		frontier: make(map[Point]bool),
	}

	// Initialize distances
	for _, node := range pf.nodes.Keys() {
		forward.distance[node] = math.MaxInt64
		backward.distance[node] = math.MaxInt64
	}

	// Set start conditions
	forward.distance[pf.start] = 0
	heap.Push(&forward.pq, &PriorityQueueItem{point: pf.start, cost: 0})
	forward.frontier[pf.start] = true

	backward.distance[pf.end] = 0
	heap.Push(&backward.pq, &PriorityQueueItem{point: pf.end, cost: 0})
	backward.frontier[pf.end] = true

	// Track the best meeting point and path length
	var meetingPoint Point
	bestPathLength := int64(math.MaxInt64)

	// Continue searching while both frontiers have nodes
	for len(forward.pq) > 0 && len(backward.pq) > 0 {
		// Alternate between forward and backward search
		if forward.pq[0].cost <= backward.pq[0].cost {
			if point := pf.stepSearch(forward, backward, true); point != nil {
				// Check if this meeting point gives a better path
				pathLength := forward.distance[point] + backward.distance[point]
				if pathLength < bestPathLength {
					bestPathLength = pathLength
					meetingPoint = point
				}
			}
		} else {
			if point := pf.stepSearch(backward, forward, false); point != nil {
				// Check if this meeting point gives a better path
				pathLength := backward.distance[point] + forward.distance[point]
				if pathLength < bestPathLength {
					bestPathLength = pathLength
					meetingPoint = point
				}
			}
		}

		// Early termination: if the smallest node in either frontier is worse
		// than our best path, we can stop
		if len(forward.pq) > 0 && len(backward.pq) > 0 {
			potentialBest := forward.pq[0].cost + backward.pq[0].cost
			if potentialBest >= bestPathLength {
				break
			}
		}
	}

	// If no meeting point was found, no path exists
	if meetingPoint == nil {
		return nil, ErrNotFound
	}

	// Reconstruct the full path
	forwardPath := pf.reconstructPartialPath(forward.previous, meetingPoint, true)
	backwardPath := pf.reconstructPartialPath(backward.previous, meetingPoint, false)

	// Combine paths (exclude the meeting point from the backward path)
	fullPath := append(forwardPath, backwardPath[1:]...)
	return fullPath, nil
}

// stepSearch performs one step of Dijkstra's algorithm in one direction
func (pf *BidirectionalPathFinder) stepSearch(
	active, other *BidirectionalSearchState,
	isForward bool,
) Point {
	if len(active.pq) == 0 {
		return nil
	}

	currentItem := heap.Pop(&active.pq).(*PriorityQueueItem)
	current := currentItem.point

	// Skip if already visited
	if active.visited[current] {
		return nil
	}

	active.visited[current] = true
	delete(active.frontier, current)

	// Check if we've reached a node visited by the other search
	if other.visited[current] {
		return current
	}

	// Explore neighbors
	var neighbors []Point
	if isForward {
		neighbors = (*get(pf.nodes, current)).Keys()
	} else {
		// For backward search, we need to find incoming edges
		// This requires checking all nodes to see if they have edges to current
		for _, node := range pf.nodes.Keys() {
			if node == current {
				continue
			}
			nodeNeighbors := (*get(pf.nodes, node)).Keys()
			for _, neighbor := range nodeNeighbors {
				if neighbor == current {
					neighbors = append(neighbors, node)
					break
				}
			}
		}
	}

	for _, neighbor := range neighbors {
		if active.visited[neighbor] {
			continue
		}

		var edgeCost int64
		if isForward {
			destMapping := get(pf.nodes, current)
			edgeCost = get(destMapping, neighbor)
		} else {
			// For backward search, get the cost from neighbor to current
			destMapping := get(pf.nodes, neighbor)
			edgeCost = get(destMapping, current)
		}

		newDistance := active.distance[current] + edgeCost

		if newDistance < active.distance[neighbor] {
			active.distance[neighbor] = newDistance
			active.previous[neighbor] = current

			// Add to priority queue (may result in duplicate entries, which are skipped when popped)
			heap.Push(&active.pq, &PriorityQueueItem{point: neighbor, cost: newDistance})
			active.frontier[neighbor] = true
		}
	}

	return nil
}

// reconstructPartialPath reconstructs a path from start/end to meeting point
func (pf *BidirectionalPathFinder) reconstructPartialPath(previous map[Point]Point, meetingPoint Point, isForward bool) []Point {
	var path []Point
	current := meetingPoint

	for current != nil {
		if isForward {
			path = append([]Point{current}, path...)
		} else {
			path = append(path, current)
		}
		current = previous[current]
	}

	return path
}
