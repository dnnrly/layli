package pathfinding

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

func TestDijkstraPathfinder_FindPaths(t *testing.T) {
	t.Run("simple two-node path", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
			},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)

		require.NotNil(t, diagram.Edges[0].Path)
		assert.Greater(t, len(diagram.Edges[0].Path.Points), 0)
	})

	t.Run("multiple edges", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
				{ID: "c", Contents: "C", Position: domain.Position{X: 3, Y: 10}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
				{ID: "e2", From: "a", To: "c"},
			},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)

		require.NotNil(t, diagram.Edges[0].Path)
		assert.Greater(t, len(diagram.Edges[0].Path.Points), 0)

		require.NotNil(t, diagram.Edges[1].Path)
		assert.Greater(t, len(diagram.Edges[1].Path.Points), 0)
	})

	t.Run("path strategy in-order", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()
		cfg.PathStrategy = "in-order"

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
			},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)

		require.NotNil(t, diagram.Edges[0].Path)
		assert.Greater(t, len(diagram.Edges[0].Path.Points), 0)
	})

	t.Run("path strategy random", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()
		cfg.PathStrategy = "random"

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
			},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)

		require.NotNil(t, diagram.Edges[0].Path)
		assert.Greater(t, len(diagram.Edges[0].Path.Points), 0)
	})

	t.Run("path points mapped correctly as integers", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b"},
			},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)

		require.NotNil(t, diagram.Edges[0].Path)
		for _, pt := range diagram.Edges[0].Path.Points {
			assert.Equal(t, pt.X, int(pt.X), "X should be an integer value")
			assert.Equal(t, pt.Y, int(pt.Y), "Y should be an integer value")
		}
	})

	t.Run("no edges produces no paths", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)
	})

	t.Run("edge class and style preserved", func(t *testing.T) {
		pf := NewDijkstraPathfinder()
		cfg := baseDiagramConfig()

		diagram := &domain.Diagram{
			Config: cfg,
			Nodes: []domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 12, Y: 3}, Width: 5, Height: 3},
			},
			Edges: []domain.Edge{
				{ID: "e1", From: "a", To: "b", Class: "my-edge", Style: "stroke: red"},
			},
		}

		err := pf.FindPaths(diagram)
		require.NoError(t, err)

		assert.Equal(t, "my-edge", diagram.Edges[0].Class)
		assert.Equal(t, "stroke: red", diagram.Edges[0].Style)
		require.NotNil(t, diagram.Edges[0].Path)
	})
}
