package layli

import (
	"errors"
	"math"

	"github.com/dnnrly/layli/algorithms/topological"
)

// LayoutArrangementFunc returns a slice of nodes arranged according to the algorithm implemented
type LayoutArrangementFunc func(c *Config) LayoutNodes

func selectArrangement(c *Config) (LayoutArrangementFunc, error) {
	switch c.Layout {
	case "":
		return LayoutFlowSquare, nil

	case "flow-square":
		return LayoutFlowSquare, nil

	case "topo-sort":
		return LayoutTopologicalSort, nil
	}

	return nil, errors.New("do not understand layout " + c.Layout)
}

func LayoutFlowSquare(c *Config) LayoutNodes {
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
			)

			pos++
		}
	}

	return nodes
}

// LayoutTopologicalSort arranges nodes in a single row, sorted in topological order
func LayoutTopologicalSort(config *Config) LayoutNodes {
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
		))
	}

	return layoutNodes
}
