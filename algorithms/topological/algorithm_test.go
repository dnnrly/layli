package topological_test

import (
	"testing"

	"github.com/dnnrly/layli/algorithms/topological"
	"github.com/stretchr/testify/assert"
)

func TestFromIsToLeft(t *testing.T) {
	g := topological.NewGraph()

	g.AddEdge("A", "B")

	all := g.RankNodes()

	assert.Equal(t, []string{"A", "B"}, all)
}

func TestHandlesCycle(t *testing.T) {
	g := topological.NewGraph()

	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A")

	all := g.RankNodes()

	assert.Equal(t, []string{"A", "B", "C"}, all)
}
