package tarjan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargan_EdgeOrderCorrect(t *testing.T) {
	g := NewGraph()

	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A")
	g.AddEdge("B", "D")
	g.AddEdge("D", "E")
	g.AddEdge("E", "C")
	g.AddEdge("E", "F")

	nodes := g.RankNodes()

	assert.Equal(t, [][]string{{"A", "B", "C", "D", "E"}, {"F"}}, nodes)
}
