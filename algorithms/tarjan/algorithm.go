package tarjan

import "slices"

type Graph struct {
	V     int
	Edges map[string][]string
	nodes []string
}

func NewGraph() *Graph {
	return &Graph{
		V:     0,
		Edges: make(map[string][]string),
		nodes: []string{},
	}
}

func (g *Graph) AddEdge(u, v string) {
	if _, ok := g.Edges[u]; !ok {
		g.nodes = append(g.nodes, u)
		g.V++
		g.Edges[u] = make([]string, 0)
	}
	g.Edges[u] = append(g.Edges[u], v)
}

func (g *Graph) RankNodes() [][]string {
	return tarjanSCC(g)
}

func tarjanSCC(graph *Graph) [][]string {
	index := 0
	stack := make([]string, 0)
	onStack := make(map[string]bool)
	lowLink := make(map[string]int)
	indexMap := make(map[string]int)
	sccs := [][]string{}

	var strongConnect func(node string)

	strongConnect = func(node string) {
		indexMap[node] = index
		lowLink[node] = index
		index++
		stack = append(stack, node)
		onStack[node] = true

		for _, neighbor := range graph.Edges[node] {
			if _, visited := indexMap[neighbor]; !visited {
				strongConnect(neighbor)
				lowLink[node] = min(lowLink[node], lowLink[neighbor])
			} else if onStack[neighbor] {
				lowLink[node] = min(lowLink[node], indexMap[neighbor])
			}
		}

		if lowLink[node] == indexMap[node] {
			scc := []string{}
			for {
				popped := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[popped] = false
				scc = append(scc, popped)
				if popped == node {
					break
				}
			}
			sccs = append(sccs, scc)
		}
	}

	for _, node := range graph.nodes {
		if _, visited := indexMap[node]; !visited {
			strongConnect(node)
		}
	}

	for _, s := range sccs {
		slices.Reverse(s)
	}
	slices.Reverse(sccs)

	return sccs
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
