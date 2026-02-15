# Pathfinding Algorithm Examples

This directory contains examples demonstrating the different pathfinding algorithms available in Layli.

## Available Algorithms

### 1. Dijkstra (Default)

**File:** `default-dijkstra.layli`

Uses the classic Dijkstra algorithm (default when no algorithm is specified). Guarantees shortest path but explores all directions equally.

<img src="default-dijkstra.svg" alt="Dijkstra pathfinding example" />

<details>
<summary>Configuration</summary>

```yaml
path:
  algorithm: dijkstra
```
</details>

### 2. A* with Euclidean Heuristic

**File:** `astar-euclidean.layli`

Uses the A* algorithm with Euclidean distance heuristic. Best for general-purpose pathfinding where straight-line distance is appropriate. More efficient than Dijkstra for most cases.

<img src="astar-euclidean.svg" alt="A* with Euclidean heuristic example" />

<details>
<summary>Configuration</summary>

```yaml
path:
  algorithm: astar
  heuristic: euclidean
```
</details>

### 3. A* with Manhattan Heuristic

**File:** `astar-manhattan.layli`

Uses the A* algorithm with Manhattan distance heuristic. Ideal for grid-based layouts where movement is restricted to horizontal/vertical directions.

<img src="astar-manhattan.svg" alt="A* with Manhattan heuristic example" />

<details>
<summary>Configuration</summary>

```yaml
path:
  algorithm: astar
  heuristic: manhattan
```
</details>

### 4. Bidirectional Dijkstra

**File:** `bidirectional.layli`

Uses bidirectional Dijkstra's algorithm, searching from both start and end points simultaneously. More efficient than standard Dijkstra for large graphs with known endpoints.

<img src="bidirectional.svg" alt="Bidirectional Dijkstra example" />

<details>
<summary>Configuration</summary>

```yaml
path:
  algorithm: bidirectional
```
</details>

---

## Cost Functions

Cost functions work alongside pathfinding algorithms to define how "distance" or "cost" is calculated between points. They influence what the algorithm considers a "short" path.

### Available Cost Functions

#### 1. Pythagorean Distance (Default)

Calculates Euclidean distance between points. Results in paths that minimize straight-line distance, but may have more corners.

```yaml
path:
  cost-function: pythagorean-distance  # or omit for default
  algorithm: dijkstra
```

**Best for:** When you want the geometrically shortest paths.

#### 2. Horizontal-Vertical

Costs 1 for horizontal/vertical moves and 2 for diagonal moves (direction changes). Results in paths that follow straight lines with fewer corners.

```yaml
path:
  cost-function: horizontal-vertical
  algorithm: dijkstra
```

**Best for:** When you prefer paths with fewer turns and more rectilinear routing.

### Combining Algorithms and Cost Functions

You can use any cost function with any algorithm:

```yaml
# A* with horizontal-vertical cost function
path:
  algorithm: astar
  heuristic: euclidean
  cost-function: horizontal-vertical
```

```yaml
# Bidirectional Dijkstra with pythagorean distance
path:
  algorithm: bidirectional
  cost-function: pythagorean-distance
```

---

## Usage

Run any example:
```bash
./layli examples/pathfinding/astar-euclidean.layli
```

Each example generates the same node layout but uses different pathfinding algorithms to route the edges between nodes.

### Experimenting with Different Configurations

To experiment with different pathfinding algorithms and cost functions:

1. Take any of the `.layli` files in this directory
2. Modify the `path` section:

```yaml
path:
  algorithm: astar           # or dijkstra, bidirectional
  heuristic: euclidean       # or manhattan (for astar only)
  cost-function: horizontal-vertical  # or pythagorean-distance
```

3. Run: `./layli your-file.layli`

The combination of algorithm, heuristic, and cost function will affect how paths are routed through your diagram.
