package pathfinding

import (
	"fmt"

	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/layout"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
)

type DijkstraPathfinder struct{}

func NewDijkstraPathfinder() *DijkstraPathfinder {
	return &DijkstraPathfinder{}
}

func (p *DijkstraPathfinder) FindPaths(diagram *domain.Diagram) error {
	cfg := toRootConfig(diagram)

	finder := func(start, end dijkstra.Point) layout.PathFinder {
		return dijkstra.NewPathFinder(start, end)
	}

	layoutObj, err := layout.NewLayoutFromConfig(finder, &cfg)
	if err != nil {
		return fmt.Errorf("finding paths: %w", err)
	}

	for i := range diagram.Edges {
		lp := findMatchingPath(layoutObj.Paths, diagram.Edges[i])
		if lp == nil {
			continue
		}

		points := make([]domain.Position, len(lp.Points))
		for j, pt := range lp.Points {
			points[j] = domain.Position{X: int(pt.X), Y: int(pt.Y)}
		}
		diagram.Edges[i].Path = &domain.Path{Points: points}
	}

	return nil
}

func findMatchingPath(paths layout.LayoutPaths, edge domain.Edge) *layout.LayoutPath {
	for _, lp := range paths {
		if lp.From == edge.From && lp.To == edge.To {
			return &lp
		}
	}
	return nil
}

func toRootConfig(d *domain.Diagram) layout.Config {
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
		Layout:         "absolute",
		LayoutAttempts: d.Config.LayoutAttempts,
		Path: layout.ConfigPath{
			Strategy: d.Config.PathStrategy,
			Attempts: d.Config.PathAttempts,
		},
		NodeWidth:  d.Config.NodeWidth,
		NodeHeight: d.Config.NodeHeight,
		Border:     d.Config.Border,
		Margin:     d.Config.Margin,
		Spacing:    d.Config.Spacing,
		Nodes:      nodes,
		Edges:      edges,
	}
}
