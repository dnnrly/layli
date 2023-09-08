package topological

type Graph struct {
	Nodes     map[string]bool
	Edges     map[string][]string
	Visited   map[string]bool
	NodeRanks map[string]int
	Sorted    []string
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:     make(map[string]bool),
		Edges:     make(map[string][]string),
		Visited:   make(map[string]bool),
		NodeRanks: make(map[string]int),
		Sorted:    []string{},
	}
}

func (g *Graph) RankNodes() []string {
	for node := range g.Nodes {
		if !g.Visited[node] {
			g.dfs(node)
		}
	}
	g.sortNodes()
	return g.Sorted
}

func (g *Graph) dfs(node string) int {
	if g.Visited[node] {
		return g.NodeRanks[node]
	}
	g.Visited[node] = true
	maxRank := -1

	for _, neighbor := range g.Edges[node] {
		neighborRank := g.dfs(neighbor)
		if neighborRank > maxRank {
			maxRank = neighborRank
		}
	}

	if maxRank == -1 {
		g.NodeRanks[node] = 0
	} else {
		g.NodeRanks[node] = maxRank + 1
	}
	g.Sorted = append(g.Sorted, node)
	return g.NodeRanks[node]
}

func (g *Graph) sortNodes() {
	for i := 0; i < len(g.Sorted)-1; i++ {
		for j := 0; j < len(g.Sorted)-i-1; j++ {
			node1, node2 := g.Sorted[j], g.Sorted[j+1]
			if g.NodeRanks[node1] > g.NodeRanks[node2] {
				g.Sorted[j], g.Sorted[j+1] = g.Sorted[j+1], g.Sorted[j]
			}
		}
	}
}

func (g *Graph) AddEdge(from, to string) {
	g.Nodes[from] = true
	g.Nodes[to] = true

	if _, ok := g.Edges[from]; !ok {
		g.Edges[from] = []string{}
	}
	g.Edges[from] = append(g.Edges[from], to)
}
