package layout

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/barkimedes/go-deepcopy"
	"github.com/dnnrly/layli/algorithms/tarjan"
	"github.com/dnnrly/layli/algorithms/topological"
)

// LayoutArrangementFunc returns a slice of nodes arranged according to the algorithm implemented
type LayoutArrangementFunc func(c *Config) (LayoutNodes, error)

// selectArrangement maps layout type strings to their implementation functions.
// When adding a new layout algorithm, add a case here and implement the corresponding
// LayoutXxx function. Also add the type to internal/domain/diagram.go and register it
// in the adapter at internal/adapters/layout/engine.go.
// See CONTRIBUTING_LAYOUTS.md for the complete checklist.
func selectArrangement(c *Config) (LayoutArrangementFunc, error) {
	switch c.Layout {
	case "":
		return LayoutFlowSquare, nil

	case "tarjan":
		return LayoutTarjan, nil

	case "flow-square":
		return LayoutFlowSquare, nil

	case "topo-sort":
		return LayoutTopologicalSort, nil

	case "random-shortest-square":
		return LayoutRandomShortestSquare, nil

	case "absolute":
		return LayoutAbsolute, nil
	}

	return nil, errors.New("do not understand layout " + c.Layout)
}

func LayoutFlowSquare(c *Config) (LayoutNodes, error) {
	numNodes := len(c.Nodes)
	nodes := make(LayoutNodes, numNodes)

	root := math.Sqrt(float64(numNodes))
	size := int(math.Ceil(root))

	if size < int(root) {
		size++
	}
	if numNodes < 4 {
		size = 2
	}

	pos := 0
	for y := 0; y < size && pos < numNodes; y++ {
		for x := 0; x < size && pos < numNodes; x++ {
			nodes[pos] = NewLayoutNode(
				c.Nodes[pos].Id,
				c.Nodes[pos].Contents,
				c.Border+
					c.Margin+
					(x*c.NodeWidth)+
					(x*(c.Margin*2)),
				c.Border+
					c.Margin+
					(y*c.NodeHeight)+
					(y*(c.Margin*2)),
				c.NodeWidth, c.NodeHeight,
				c.Nodes[pos].Class,
				c.Nodes[pos].Style,
			)

			pos++
		}
	}

	return nodes, nil
}

// LayoutTopologicalSort arranges nodes in a single row, sorted in topological order
func LayoutTopologicalSort(config *Config) (LayoutNodes, error) {
	layoutNodes := LayoutNodes{}
	graph := topological.NewGraph()

	for _, e := range config.Edges {
		graph.AddEdge(e.From, e.To)
	}

	rankedNodes := graph.RankNodes()

	for i, id := range rankedNodes {
		c := config.Nodes.ByID(id)

		layoutNodes = append(layoutNodes, NewLayoutNode(
			id, c.Contents,
			config.Border+
				config.Margin+
				(i*config.NodeWidth)+
				(i*(config.Margin*2)),
			config.Border+
				config.Margin+
				(0*config.NodeHeight)+
				(0*(config.Margin*2)),
			config.NodeWidth, config.NodeHeight,
			c.Class,
			c.Style,
		))
	}

	return layoutNodes, nil
}

// LayoutTarjan arranges nodes in multiple rows according to Tarhan's algorithm
func LayoutTarjan(config *Config) (LayoutNodes, error) {
	layoutNodes := LayoutNodes{}
	graph := tarjan.NewGraph()

	for _, e := range config.Edges {
		graph.AddEdge(e.From, e.To)
	}

	nodes := graph.RankNodes()

	for row, rNodes := range nodes {
		for col, id := range rNodes {
			c := config.Nodes.ByID(id)

			layoutNodes = append(layoutNodes, NewLayoutNode(
				id, c.Contents,
				config.Border+
					config.Margin+
					(row*config.NodeWidth)+
					(row*(config.Margin*2)),
				config.Border+
					config.Margin+
					(col*config.NodeHeight)+
					(col*(config.Margin*2)),
				config.NodeWidth, config.NodeHeight,
				c.Class,
				c.Style,
			))
		}
	}

	return layoutNodes, nil
}

func LayoutRandomShortestSquare(config *Config) (LayoutNodes, error) {
	return shuffleNodes(config, LayoutFlowSquare)
}

func shuffleNodes(config *Config, arrange LayoutArrangementFunc) (LayoutNodes, error) {
	c := deepcopy.MustAnything(config).(*Config)
	var shortest LayoutNodes
	shortestDist := math.MaxFloat64

	for i := 0; i < config.LayoutAttempts; i++ {
		rand.Shuffle(len(c.Nodes), func(i, j int) { c.Nodes[i], c.Nodes[j] = c.Nodes[j], c.Nodes[i] })
		nodes, _ := arrange(c)
		dist, _ := nodes.ConnectionDistances(c.Edges)
		if dist < shortestDist {
			shortest = nodes
			shortestDist = dist
		}
	}

	return shortest, nil
}

func LayoutAbsolute(c *Config) (LayoutNodes, error) {
	numNodes := len(c.Nodes)
	nodes := make(LayoutNodes, numNodes)

	for i, n := range c.Nodes {
		if n.Position.X < c.Border || n.Position.Y < c.Border {
			return nil, fmt.Errorf("node %s overlaps border", n.Id)
		}
		if n.Position.X < c.Border+c.Margin || n.Position.Y < c.Border+c.Margin {
			return nil, fmt.Errorf("node %s margin overlaps border", n.Id)
		}
		nodes[i] = NewLayoutNode(
			n.Id, n.Contents,
			n.Position.X,
			n.Position.Y,
			c.NodeWidth, c.NodeHeight,
			n.Class,
			n.Style,
		)
	}

	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i != j {
				if nodesOverlap(node1, node2) {
					return nil, fmt.Errorf("nodes %s and %s overlap", node1.Id, node2.Id)
				}
				if marginsOverlap(node1, node2, c.Margin) {
					return nil, fmt.Errorf("nodes %s and %s margins overlap", node1.Id, node2.Id)
				}
			}
		}
	}

	return nodes, nil
}

func nodesOverlap(node1, node2 LayoutNode) bool {
	return !(node1.right <= node2.left ||
		node1.left >= node2.right ||
		node1.bottom <= node2.top ||
		node1.top >= node2.bottom)
}

func marginsOverlap(node1, node2 LayoutNode, margin int) bool {
	return !(node1.right+margin <= node2.left-margin ||
		node1.left-margin >= node2.right+margin ||
		node1.bottom+margin <= node2.top-margin ||
		node1.top-margin >= node2.bottom+margin)
}
