# Layli Examples

Here are examples of how to use Layli, organized from simple to advanced.

## Getting Started

### Minimal Example

The simplest possible Layli diagram with nodes and a single edge:

<details>
<summary>Simple flow-square example</summary>

```yaml
nodes:
    - id: a
      contents: Node 1
    - id: b
      contents: Node 2
    - id: c
      contents: Node 3
    - id: f
      contents: Node 4
    - id: d
      contents: Node 5
    - id: e
      contents: Node 6
    - id: g
      contents: Node 7
    - id: h
      contents: Node 8
    - id: i
      contents: Node 9

edges:
    - from: a
      to: b
    - from: b
      to: c
    - from: c
      to: d
    - from: d
      to: e
    - from: c
      to: e
    - from: e
      to: d
    - from: d
      to: f
    - from: f
      to: g
    - from: f
      to: h
    - from: g
      to: i

width: 7
height: 4
```

**Output:**
<img src="/examples/simple-flow-square.svg" alt="Simple Flow Square example image" />
</details>

---

## Layouts

The layout algorithm controls where nodes are positioned on the diagram. Layli supports several layouts:

### 1. Flow Square (Default)

The Flow Square layout arranges nodes in a grid, filling rows and columns in the order you specify them. This is the default layout.

**File:** `simple-flow-square.layli`

<img src="/examples/simple-flow-square.svg" alt="Flow Square example image" />

**When to use:** When you want a simple grid arrangement of nodes.

### 2. Absolute

Specify exactly where each node should appear on the diagram. You provide x,y coordinates for each node.

**File:** `absolute.layli`

<img src="/examples/absolute.svg" alt="Absolute example image" />

<details>
<summary>Absolute layout example</summary>

```yaml
nodes:
    - id: a
      contents: Node 1
      position: {x: 5, y: 5}
    - id: b
      contents: Node 2
      position: {x: 5, y: 15}
    - id: c
      contents: Node 3
      position: {x: 5, y: 25}
    - id: f
      contents: Node 4
      position: {x: 15, y: 3}
    - id: d
      contents: Node 5
      position: {x: 12, y: 10}
    - id: e
      contents: Node 6
      position: {x: 12, y: 22}
    - id: g
      contents: Node 7
      position: {x: 20, y: 15}
    - id: h
      contents: Node 8
      position: {x: 20, y: 25}
    - id: i
      contents: Node 9
      position: {x: 25, y: 5}

layout: absolute

edges:
    - from: a
      to: b
    - from: b
      to: c
    - from: c
      to: d
    - from: d
      to: e
    - from: c
      to: e
    - from: e
      to: d
    - from: d
      to: f
    - from: f
      to: g
    - from: f
      to: h
    - from: g
      to: i

width: 4
height: 4
```
</details>

**When to use:** When you have a specific diagram layout in mind and want precise control.

### 3. Topological Sort

Arranges nodes in a single row, ordered by their connections in the graph.

**File:** `topological-sort.layli`

<img src="/examples/topological-sort.svg" alt="Topological sort example image" />

<details>
<summary>Topological sort example</summary>

```yaml
layout: topo-sort

nodes:
  - id: node1
    contents: "First Node"
  - id: node2
    contents: "Second Node"
  - id: node3
    contents: "Third Node"
  - id: node4
    contents: "Forth Node"
  - id: node5
    contents: "Fifth Node"

edges:
  - from: node1
    to: node2
  - from: node3
    to: node2
  - from: node3
    to: node4
  - from: node5
    to: node3
  - from: node2
    to: node5
```
</details>

**When to use:** When you want nodes arranged in dependency order.

### 4. Tarjan's Algorithm

Uses Tarjan's strongly connected components algorithm to arrange nodes in an appealing, layered way.

**File:** `tarjan.layli` (unstable - regenerated manually)

<img src="/examples/tarjan.svg" alt="Tarjan's algorithm example image" />

**When to use:** When you want an automatic, aesthetically pleasing layout.

### 5. Random Shortest Square

Tries many random arrangements and selects the one with the shortest total edge length.

**File:** `random-shortest-square.layli` (unstable - regenerated manually)

<img src="/examples/random-shortest-square.svg" alt="Random Shortest Square example image" />

<details>
<summary>Random Shortest Square example</summary>

```yaml
layout: random-shortest-square
layout-attempts: 1000

nodes:
  - id: node1
    contents: "Node 1"
  # ... more nodes ...

edges:
  - from: node1
    to: node2
  # ... more edges ...
```
</details>

**When to use:** When you want automatic layout with optimization for short edges.

---

## Paths & Routing

Paths are the connections between nodes. This section covers how paths work and how to customize them.

### Path Grid

Layli routes paths across a grid of points that don't overlap with nodes. You can visualize this grid with the `--show-grid` flag.

### Cost Functions

The pathfinding algorithm uses a cost function to determine what constitutes a "short" path. Different cost functions produce different routing styles.

#### 1. Pythagorean Distance (Default)

Optimizes for shortest Euclidean distance. May result in paths with more corners.

```yaml
path:
  cost-function: pythagorean-distance  # or omit for default
```

#### 2. Horizontal-Vertical

Costs 1 for horizontal/vertical moves and 2 for diagonal moves. Results in paths that follow straight lines with fewer direction changes.

**File:** `horizontal-vertical.layli`

<img src="/examples/horizontal-vertical.svg" alt="Horizontal-vertical cost function example image" />

<details>
<summary>Horizontal-vertical cost function example</summary>

```yaml
path:
  cost-function: horizontal-vertical

layout: flow-square

nodes:
  - id: a
    contents: Node A
  - id: b
    contents: Node B
  - id: c
    contents: Node C
  - id: d
    contents: Node D
  - id: e
    contents: Node E
  - id: f
    contents: Node F
  - id: g
    contents: Node G
  - id: h
    contents: Node H
  - id: i
    contents: Node I

edges:
  - from: a
    to: b
  - from: b
    to: c
  - from: c
    to: d
  - from: d
    to: e
  - from: c
    to: e
  - from: e
    to: d
  - from: d
    to: f
  - from: f
    to: g
  - from: f
    to: h
  - from: g
    to: i

width: 7
height: 4
border: 2
margin: 2
```
</details>

### Path Strategies

Layli enforces that paths cannot cross. If crossings occur, you can use a path strategy to find a non-crossing arrangement.

#### In-Order Strategy (Default)

Paths are drawn in the order you specify them.

```yaml
path:
  strategy: in-order
```

#### Random Strategy

Tries multiple random orderings and selects the arrangement with the shortest total path length.

**File:** `random-paths.layli` (unstable - regenerated manually)

<img src="/examples/random-paths.svg" alt="Random paths example image" />

<details>
<summary>Random paths example</summary>

```yaml
path:
  strategy: random
  attempts: 100

layout: random-shortest-square
layout-attempts: 100

nodes:
  - id: node1
    contents: "Node 1"
  # ... more nodes ...

edges:
  - from: node1
    to: node2
  # ... more edges ...
```
</details>

**When to use:** When in-order routing causes path crossings.

---

## Customization

### Size and Spacing

Control node dimensions and spacing.

**File:** `size-and-spacing.layli`

<img src="/examples/size-and-spacing.svg" alt="Size and spacing example image" />

<details>
<summary>Size and spacing example</summary>

```yaml
# Node dimensions (in grid units)
width: 7
height: 4

# Space between nodes (in grid units)
margin: 3

# Border around entire diagram (in grid units)
border: 1

nodes:
    - id: a
      contents: Node 1
    - id: b
      contents: Node 2
    - id: c
      contents: Node 3
    - id: d
      contents: Node 4

edges:
    - from: a
      to: b
    - from: b
      to: c
    - from: c
      to: d
```
</details>

### Styling

Apply CSS styles to nodes and edges.

**File:** `style.layli`

<img src="/examples/style.svg" alt="Style example image" />

<details>
<summary>Adding style</summary>

```yaml
nodes:
    - id: a
      contents: Node 1
      style: "fill:cyan; stroke:red;"
    - id: b
      contents: Node 2
      style: "fill:cyan; stroke:magenta;"
      class: class-2
    - id: c
      contents: Node 3
      class: class-1
    - id: d
      contents: Node 4

edges:
    - id: p1
      from: a
      to: b
      class: class-1
    - id: p2
      from: b
      to: c
      style: stroke:green
    - id: p3
      from: c
      to: d

styles:
    .class-1: >
      fill: azure
      stroke: blue
      stroke-width: 3
    .class-2: stroke:green
```
</details>

---

## Example Files by Feature

| File | Feature | Status |
|------|---------|--------|
| `simple-flow-square.layli` | Basic layout | ✅ Stable |
| `absolute.layli` | Explicit positioning | ✅ Stable |
| `horizontal-vertical.layli` | Cost function | ✅ Stable |
| `size-and-spacing.layli` | Dimensions | ✅ Stable |
| `style.layli` | Styling | ✅ Stable |
| `topological-sort.layli` | Topo-sort layout | ⚠️ Occasionally fails |
| `tarjan.layli` | Tarjan layout | ⚠️ Occasionally fails |
| `random-shortest-square.layli` | Random layout | ⚠️ Nondeterministic |
| `random-paths.layli` | Random path strategy | ⚠️ Nondeterministic |

**Note:** The "occasionally fails" examples are due to complexity in the algorithms. The "nondeterministic" examples use randomization. Run `make examples` to regenerate stable examples only.

---

## Pathfinding Algorithms

For advanced control over how paths are routed, see the [pathfinding examples](pathfinding/) which demonstrate:
- Dijkstra (default)
- A* with Euclidean heuristic
- A* with Manhattan heuristic
- Bidirectional Dijkstra

Each uses the same diagram but with different pathfinding algorithms.

---

## Next Steps

1. Start with `simple-flow-square.layli` to understand basic structure
2. Try different layouts with the same nodes and edges
3. Experiment with cost functions and path strategies
4. Customize sizing, spacing, and styles
5. Explore different pathfinding algorithms in [pathfinding/](pathfinding/)
6. Check out the [main README](../README.md) for more information
