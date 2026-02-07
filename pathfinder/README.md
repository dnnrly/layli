# Pathfinder Module

The pathfinder module provides algorithms for finding the shortest path between two points in a graph.

## Available Algorithms

### Dijkstra's Algorithm (`dijkstra/`)

The classic single-source shortest path algorithm. Explores all directions equally from the start node.

**When to use:**
- General shortest path problems
- Graphs where heuristic information isn't available
- When you need to find paths to multiple targets from one source

**Files:**
- `dijkstra.go` - Core implementation
- `dijkstra_test.go` - Tests

### A* (A-star) Algorithm (`dijkstra/astar.go`)

An optimized variant of Dijkstra that uses a heuristic function to guide the search toward the goal.

**When to use:**
- Single-source, single-target pathfinding
- When you have heuristic information about distance to goal
- Large graphs where reducing explored nodes matters
- Grid-based navigation

**Key features:**
- Guaranteed to find shortest path (with admissible heuristic)
- Customizable heuristic functions
- More efficient than Dijkstra when heuristics are good

**Files:**
- `dijkstra/astar.go` - A* implementation
- `dijkstra/astar_test.go` - Tests
- `dijkstra/astar_benchmark_test.go` - Performance benchmarks
- `dijkstra/ASTAR.md` - Detailed documentation

## Interface

All pathfinding algorithms implement the `PathFinder` interface:

```go
type PathFinder interface {
    AddConnection(from Point, cost CostFunction, to ...Point)
    BestPath() ([]Point, error)
}
```

This allows algorithms to be swapped as needed.

## Usage Example

```go
package main

import (
    "github.com/dnnrly/layli/pathfinder/dijkstra"
)

func main() {
    // Create a pathfinder from (1,1) to (5,5)
    pf := dijkstra.NewAStarPathFinder(
        dijkstra.coordinate{x: 1, y: 1},
        dijkstra.coordinate{x: 5, y: 5},
    )

    // Define edge cost (e.g., Euclidean distance)
    cost := dijkstra.PythagoreanDistance

    // Add connections
    pf.AddConnection(
        dijkstra.coordinate{x: 1, y: 1},
        cost,
        dijkstra.coordinate{x: 1, y: 2},
        dijkstra.coordinate{x: 2, y: 1},
    )

    // Find the best path
    path, err := pf.BestPath()
    if err != nil {
        // No path found
        panic(err)
    }

    // Use the path
    for _, point := range path {
        x, y := point.Coordinates()
        println(x, y)
    }
}
```

## Cost Functions

Pass a custom `CostFunction` to `AddConnection` to define edge weights:

```go
type CostFunction func(from, to Point) int64

// Example: Pythagorean distance
func PythagoreanDistance(from, to Point) int64 {
    fromX, fromY := from.Coordinates()
    toX, toY := to.Coordinates()
    dx := fromX - toX
    dy := fromY - toY
    return int64(math.Sqrt(dx*dx + dy*dy))
}

// Example: Uniform cost
func UniformCost(from, to Point) int64 {
    return 1
}

// Example: Manhattan distance
func ManhattanDistance(from, to Point) int64 {
    fromX, fromY := from.Coordinates()
    toX, toY := to.Coordinates()
    return int64(math.Abs(fromX-toX) + math.Abs(fromY-toY))
}
```

## Testing

Run tests for all algorithms:

```bash
go test ./pathfinder/dijkstra -v
```

Run tests for specific algorithm:

```bash
go test ./pathfinder/dijkstra -v -run Dijkstra
go test ./pathfinder/dijkstra -v -run AStar
```

## Benchmarking

Compare algorithm performance:

```bash
go test ./pathfinder/dijkstra -bench=. -benchmem -run=^$
```

## Adding New Algorithms

To add a new pathfinding algorithm:

1. Implement the `PathFinder` interface in a new file
2. Write comprehensive tests following the pattern in `dijkstra_test.go`
3. Add benchmarks to compare with existing algorithms
4. Document the algorithm's characteristics and use cases
5. Update this README

Key interface requirements:

```go
type PathFinder interface {
    // Add directed edges from 'from' to each point in 'to'
    // Cost is determined by the CostFunction
    AddConnection(from Point, cost CostFunction, to ...Point)
    
    // Return the shortest path from start to end
    // Return ErrNotFound if no path exists
    BestPath() ([]Point, error)
}
```

## Integration with Layli

The layout engine uses pathfinding to route edges between nodes. Currently it uses Dijkstra via the `CreateFinder` function in `layout/path.go`.

To switch to A*:

```go
// In layout/path.go, the FindPath method
finder := dijkstra.NewAStarPathFinder(
    nFrom.GetCentre(),
    nTo.GetCentre(),
)
```

See `pathfinder/dijkstra/ASTAR.md` for more details.
