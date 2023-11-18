package layli

import (
	"errors"
	"math"
	"math/rand"

	"github.com/barkimedes/go-deepcopy"
	"github.com/dnnrly/layli/algorithms/tarjan"
	"github.com/dnnrly/layli/algorithms/topological"
)

// LayoutArrangementFunc returns a slice of nodes arranged according to the algorithm implemented
type LayoutArrangementFunc func(c *Config) LayoutNodes

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

// LayoutTarjan arranges nodes in multiple rows according to Tarhan's algorithm
func LayoutTarjan(config *Config) LayoutNodes {
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
			))
		}
	}

	return layoutNodes
}

func LayoutRandomShortestSquare(config *Config) LayoutNodes {
	return shuffleNodes(config, LayoutFlowSquare)
}

func shuffleNodes(config *Config, arrange func(config *Config) LayoutNodes) LayoutNodes {
	c := deepcopy.MustAnything(config).(*Config)
	var shortest LayoutNodes
	shortestDist := math.MaxFloat64

	for i := 0; i < config.LayoutAttempts; i++ {
		rand.Shuffle(len(c.Nodes), func(i, j int) { c.Nodes[i], c.Nodes[j] = c.Nodes[j], c.Nodes[i] })
		nodes := arrange(c)
		dist, _ := nodes.ConnectionDistances(c.Edges)
		if dist < shortestDist {
			shortest = nodes
			shortestDist = dist
		}
	}

	return shortest
}
