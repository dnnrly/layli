# Pathfinding Algorithm Examples

This directory contains examples demonstrating the different pathfinding algorithms available in Layli.

## Available Algorithms

### 1. A* with Euclidean Heuristic
**File:** `astar-euclidean.layli`

Uses the A* algorithm with Euclidean distance heuristic. Best for general-purpose pathfinding where straight-line distance is appropriate.

### 2. A* with Manhattan Heuristic  
**File:** `astar-manhattan.layli`

Uses the A* algorithm with Manhattan distance heuristic. Ideal for grid-based layouts where movement is restricted to horizontal/vertical directions.

### 3. Bidirectional Dijkstra
**File:** `bidirectional.layli`

Uses bidirectional Dijkstra's algorithm, searching from both start and end points simultaneously. More efficient for large graphs with known endpoints.

### 4. Default Dijkstra
**File:** `default-dijkstra.layli`

Uses the classic Dijkstra algorithm (default when no algorithm is specified). Guarantees shortest path but explores all directions equally.

## Usage

Run any example:
```bash
./layli examples/pathfinding/astar-euclidean.layli
```

Each example generates the same node layout but uses different pathfinding algorithms to route the edges between nodes.
