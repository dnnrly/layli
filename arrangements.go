package layli

import "math"

// LayoutArrangementFunc returns a slice of nodes arranged according to the algorithm implemented
type LayoutArrangementFunc func(c *Config) LayoutNodes

func selectArrangement(c *Config) LayoutArrangementFunc {
	switch c.Layout {
	case "":
		return LayoutFlowSquare

	case "flow-square":
		return LayoutFlowSquare

	case "topo-sort":
		return LayoutTopologicalSort
	}

	return nil
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

func LayoutTopologicalSort(c *Config) LayoutNodes {
	return LayoutNodes{}
}
