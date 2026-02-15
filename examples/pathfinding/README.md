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

## Usage

Run any example:
```bash
./layli examples/pathfinding/astar-euclidean.layli
```

Each example generates the same node layout but uses different pathfinding algorithms to route the edges between nodes.
