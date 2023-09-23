package topological

type Graph struct {
	Nodes     []string
	Edges     map[string][]string
	Visited   map[string]bool
	NodeRanks map[string]int
	Sorted    []string
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:     []string{},
		Edges:     make(map[string][]string),
		Visited:   make(map[string]bool),
		NodeRanks: make(map[string]int),
		Sorted:    []string{},
	}
}

func (g *Graph) RankNodes() []string {
	g.Nodes = removeDuplicates[string](g.Nodes)
	for _, node := range g.Nodes {
		if !g.Visited[node] {
			g.dfs(node)
		}
	}
	g.sortNodes()
	for i, j := 0, len(g.Sorted)-1; i < j; i, j = i+1, j-1 {
		g.Sorted[i], g.Sorted[j] = g.Sorted[j], g.Sorted[i]
	}

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
	g.Nodes = append(g.Nodes, from)
	g.Nodes = append(g.Nodes, to)

	if _, ok := g.Edges[from]; !ok {
		g.Edges[from] = []string{}
	}
	g.Edges[from] = append(g.Edges[from], to)
}

func removeDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	uniqueSlice := []T{}

	for _, element := range slice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}

	return uniqueSlice
}
