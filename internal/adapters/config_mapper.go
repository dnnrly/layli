package adapters

import (
	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/layout"
)

// ToLayoutConfig converts a domain diagram to a layout config.
// This is used by multiple adapters to ensure consistent type mapping.
func ToLayoutConfig(d *domain.Diagram) layout.Config {
	nodes := make(layout.ConfigNodes, len(d.Nodes))
	for i, n := range d.Nodes {
		nodes[i] = layout.ConfigNode{
			Id:       n.ID,
			Contents: n.Contents,
			Position: layout.Position{
				X: n.Position.X,
				Y: n.Position.Y,
			},
			Class: n.Class,
			Style: n.Style,
		}
	}

	edges := make(layout.ConfigEdges, len(d.Edges))
	for i, e := range d.Edges {
		edges[i] = layout.ConfigEdge{
			ID:    e.ID,
			From:  e.From,
			To:    e.To,
			Class: e.Class,
			Style: e.Style,
		}
	}

	return layout.Config{
		Layout:         string(d.Config.LayoutType),
		LayoutAttempts: d.Config.LayoutAttempts,
		NodeWidth:      d.Config.NodeWidth,
		NodeHeight:     d.Config.NodeHeight,
		Border:         d.Config.Border,
		Margin:         d.Config.Margin,
		Spacing:        d.Config.Spacing,
		Nodes:          nodes,
		Edges:          edges,
	}
}

// ToLayoutConfigWithPath is like ToLayoutConfig but includes the path configuration.
// Used by the pathfinder adapter.
func ToLayoutConfigWithPath(d *domain.Diagram) layout.Config {
	cfg := ToLayoutConfig(d)
	cfg.Path = layout.ConfigPath{
		Strategy: d.Config.PathStrategy,
		Attempts: d.Config.PathAttempts,
	}
	return cfg
}

// ToLayoutConfigWithFullPaths is like ToLayoutConfigWithPath but includes pathfinding algorithm settings.
func ToLayoutConfigWithFullPaths(d *domain.Diagram) layout.Config {
	cfg := ToLayoutConfig(d)
	cfg.Path = layout.ConfigPath{
		Strategy:  d.Config.PathStrategy,
		Attempts:  d.Config.PathAttempts,
		Algorithm: string(d.Config.Pathfinding.Algorithm),
		Heuristic: string(d.Config.Pathfinding.Heuristic),
	}
	return cfg
}
