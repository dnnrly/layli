package layout

import (
	"testing"

	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseDiagramConfig() domain.DiagramConfig {
	return domain.DiagramConfig{
		NodeWidth:      5,
		NodeHeight:     3,
		Border:         1,
		Margin:         2,
		Spacing:        20,
		LayoutAttempts: 10,
		PathAttempts:   20,
	}
}

func TestLayoutAdapter_Arrange(t *testing.T) {
	t.Run("flow square layout", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutFlowSquare

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
				{ID: "b", Contents: "B"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		assert.Equal(t, 3, diagram.Nodes[0].Position.X)
		assert.Equal(t, 3, diagram.Nodes[0].Position.Y)
		assert.Equal(t, 12, diagram.Nodes[1].Position.X)
		assert.Equal(t, 3, diagram.Nodes[1].Position.Y)
	})

	t.Run("empty layout type defaults to flow-square", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = ""

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
				{ID: "b", Contents: "B"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		assert.Equal(t, 3, diagram.Nodes[0].Position.X)
		assert.Equal(t, 3, diagram.Nodes[0].Position.Y)
		assert.Equal(t, 12, diagram.Nodes[1].Position.X)
		assert.Equal(t, 3, diagram.Nodes[1].Position.Y)
	})

	t.Run("topological sort layout", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutTopoSort

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
				{ID: "b", Contents: "B"},
				{ID: "c", Contents: "C"},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
				{ID: "e2", From: "b", To: "c"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		nodeByID := map[string]domain.Node{}
		for _, n := range diagram.Nodes {
			nodeByID[n.ID] = n
		}

		assert.Less(t, nodeByID["a"].Position.X, nodeByID["b"].Position.X)
		assert.Less(t, nodeByID["b"].Position.X, nodeByID["c"].Position.X)
	})

	t.Run("tarjan layout", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutTarjan

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
				{ID: "b", Contents: "B"},
				{ID: "c", Contents: "C"},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
				{ID: "e2", From: "b", To: "c"},
				{ID: "e3", From: "c", To: "a"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		for _, n := range diagram.Nodes {
			assert.Greater(t, n.Width, 0, "node %s should have width", n.ID)
			assert.Greater(t, n.Height, 0, "node %s should have height", n.ID)
		}
	})

	t.Run("absolute layout", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutAbsolute

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		assert.Equal(t, 3, diagram.Nodes[0].Position.X)
		assert.Equal(t, 3, diagram.Nodes[0].Position.Y)
		assert.Equal(t, 12, diagram.Nodes[1].Position.X)
		assert.Equal(t, 3, diagram.Nodes[1].Position.Y)
	})

	t.Run("random shortest square", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutRandomShortest

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
				{ID: "b", Contents: "B"},
				{ID: "c", Contents: "C"},
				{ID: "d", Contents: "D"},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
				{ID: "e2", From: "c", To: "d"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		for _, n := range diagram.Nodes {
			assert.Greater(t, n.Width, 0, "node %s should have width", n.ID)
			assert.Greater(t, n.Height, 0, "node %s should have height", n.ID)
			assert.GreaterOrEqual(t, n.Position.X, 0, "node %s X should be non-negative", n.ID)
			assert.GreaterOrEqual(t, n.Position.Y, 0, "node %s Y should be non-negative", n.ID)
		}
	})

	t.Run("unknown layout type returns error", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = "unknown-layout"

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
			},
		}

		err := adapter.Arrange(diagram)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown layout type")
	})

	t.Run("node dimensions set from config", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutFlowSquare
		cfg.NodeWidth = 8
		cfg.NodeHeight = 4

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A"},
				{ID: "b", Contents: "B"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		for _, n := range diagram.Nodes {
			assert.Equal(t, 8, n.Width, "node %s should have width 8", n.ID)
			assert.Equal(t, 4, n.Height, "node %s should have height 4", n.ID)
		}
	})

	t.Run("class and style preservation", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutFlowSquare

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Class: "my-class", Style: "fill: red"},
				{ID: "b", Contents: "B", Class: "other-class", Style: "stroke: blue"},
			},
		}

		err := adapter.Arrange(diagram)
		require.NoError(t, err)

		assert.Equal(t, "my-class", diagram.Nodes[0].Class)
		assert.Equal(t, "fill: red", diagram.Nodes[0].Style)
		assert.Equal(t, "other-class", diagram.Nodes[1].Class)
		assert.Equal(t, "stroke: blue", diagram.Nodes[1].Style)
	})

	t.Run("absolute layout errors on overlapping nodes", func(t *testing.T) {
		adapter := NewLayoutAdapter()
		cfg := baseDiagramConfig()
		cfg.LayoutType = domain.LayoutAbsolute

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}},
				{ID: "b", Contents: "B", Position: domain.Position{X: 4, Y: 3}},
			},
		}

		err := adapter.Arrange(diagram)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "overlap")
	})
}
