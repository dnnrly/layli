# Layli examples

Here are some examples of how to use Layli.

## Layouts

Sometimes referred to as arrangement, this controls where nodes are positioned in the image.

### Flow Squares

The Flow Square layout arranges the nodes in a square grid, in the order that they were specified in the definition file. This is the default layout. This is the default arrangement if you don't specify anything.

<img src="/examples/simple-flow-square.svg" alt="Simple Flow Square example image" />

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
</details>

### Topological Sort

This layout uses an algorithm to analyse the edges in the graph to arrange the nodes in a single, in the order in which they are connected together. This order takes in to account the order that the paths are specified in.

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

### Tarjan's Algorithm

This is an implementation of the [Tarjan's algorithm](https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm), arranging the nodes in a way that tries to be 'appealing'. This arrangement can sometimes take a long time to process.


<img src="/examples/tarjan.svg" alt="Tarjan's algorithm example image" />

<details>
<summary>Tarjan's algorithm example</summary>

```yaml
layout: tarjan

nodes:
    - id: a
      contents: Node 1
    - id: b
      contents: Node 2
    - id: c
      contents: Node 3
    - id: d
      contents: Node 4
    - id: e
      contents: Node 5
    - id: f
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
    - from: g
      to: e
    - from: d
      to: f
    - from: f
      to: g
    - from: f
      to: h

```
</details>

### Random Shortest Square

This algorithm attempts to arrange the nodes in a square grid, but it does this by randomly selecting the nodes to place many times over. It selects the arrangement with the shortest total distance of all of the specified edges. This distance is just the distance directly between the centre of the 2 nodes on an edge. You can set the number of attempts to find an arrangement with the `layout-attempts` parameter.

<img src="/examples/random-shortest-square.svg" alt="Random Shortest Square example image" />

<details>
<summary>Random Shortest Square example</summary>

```yaml
layout: random-shortest-square
layout-attempts: 1000

nodes:
  - id: node1
    contents: "Node 1"
  - id: node2
    contents: "Node 2"
  - id: node3
    contents: "Node 3"
  - id: node4
    contents: "Node 4"
  - id: node5
    contents: "Node 5"
  - id: node6
    contents: "Node 6"
  - id: node7
    contents: "Node 7"
  - id: node8
    contents: "Node 8"
  - id: node9
    contents: "Node 9"
  - id: node10
    contents: "Node 10"
  - id: node11
    contents: "Node 11"
  - id: node12
    contents: "Node 12"
  - id: node13
    contents: "Node 13"
  - id: node14
    contents: "Node 14"

edges:
  - from: node1
    to: node2
  - from: node2
    to: node3
  - from: node3
    to: node7
  - from: node7
    to: node11
  - from: node11
    to: node10
  - from: node10
    to: node9
  - from: node9
    to: node5
  - from: node5
    to: node1
  - from: node6
    to: node12
```
</details>

### Absolute

This is a fairly straight forward layout, you specify where you would like the nodes to appear on the image and layli will look after the paths. Watch out though, you layli will fail if nodes are too close to the border or each other.

<img src="/examples/absolute.svg" alt="Absolute example image" />

<details>
<summary>Absolute example</summary>

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

## Paths

Paths are defined by selecting a `from` node an a `to` node in the `edges` configuration. To generate the path, `layli` uses [Dijkstr's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm) to find the shortest path across a grid of points that are not covered by a node. You can see this grid by using the `--show-grid` option when you run the command.

Paths connect to nodes on a 'port', which is any grid point that sits on the border of the node but is not a corner.

### Avoiding crossed paths

Layli does **not** allow paths to cross. If one is detected then layli will exit with an error. To avoid this situation, it's possible to select a different path strategy.

#### Random path strategy

You can specify a strategy that will randomly shuffle the order of the paths a number of times and select the one with the shortest total path length. This should hopefully find an arrangement where paths do not cross. It's also possible to specify the number of attempts to find paths that do not cross.

Notice that in this example, the layout can be specified too.

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
  - id: node2
    contents: "Node 2"
  - id: node3
    contents: "Node 3"
  - id: node4
    contents: "Node 4"
  - id: node5
    contents: "Node 5"
  - id: node6
    contents: "Node 6"
  - id: node7
    contents: "Node 7"
  - id: node8
    contents: "Node 8"
  - id: node9
    contents: "Node 9"
  - id: node10
    contents: "Node 10"
  - id: node11
    contents: "Node 11"
  - id: node12
    contents: "Node 12"
  - id: node13
    contents: "Node 13"
  - id: node14
    contents: "Node 14"

edges:
  - from: node1
    to: node2
  - from: node2
    to: node3
  - from: node3
    to: node7
  - from: node7
    to: node11
  - from: node11
    to: node10
  - from: node10
    to: node9
  - from: node9
    to: node5
  - from: node5
    to: node1
  - from: node6
    to: node12

```
</details>

## Size and spacing

It is possible to specify the size of nodes and the spacing between them. It's also possible to specify a margin around the edge of the image where no paths will be drawn.

<img src="/examples/size-and-spacing.svg" alt="Size and spacing example image" />

<details>
<summary>Size and spacing example</summary>

```yaml
width: 7
height: 4
margin: 3

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