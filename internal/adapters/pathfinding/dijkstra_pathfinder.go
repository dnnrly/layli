package pathfinding

import (
	"fmt"

	"github.com/dnnrly/layli"
	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
)

type DijkstraPathfinder struct{}

func NewDijkstraPathfinder() *DijkstraPathfinder {
	return &DijkstraPathfinder{}
}

func (p *DijkstraPathfinder) FindPaths(diagram *domain.Diagram) error {
	cfg := toRootConfig(diagram)

	finder := func(start, end dijkstra.Point) layli.PathFinder {
		return dijkstra.NewPathFinder(start, end)
	}

	layout, err := layli.NewLayoutFromConfig(finder, &cfg)
	if err != nil {
		return fmt.Errorf("finding paths: %w", err)
	}

	for i := range diagram.Edges {
		lp := findMatchingPath(layout.Paths, diagram.Edges[i])
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

func findMatchingPath(paths layli.LayoutPaths, edge domain.Edge) *layli.LayoutPath {
	for _, lp := range paths {
		if lp.From == edge.From && lp.To == edge.To {
			return &lp
		}
	}
	return nil
}

func toRootConfig(d *domain.Diagram) layli.Config {
	nodes := make(layli.ConfigNodes, len(d.Nodes))
	for i, n := range d.Nodes {
		nodes[i] = layli.ConfigNode{
			Id:       n.ID,
			Contents: n.Contents,
			Position: layli.Position{
				X: n.Position.X,
				Y: n.Position.Y,
			},
			Class: n.Class,
			Style: n.Style,
		}
	}

	edges := make(layli.ConfigEdges, len(d.Edges))
	for i, e := range d.Edges {
		edges[i] = layli.ConfigEdge{
			ID:    e.ID,
			From:  e.From,
			To:    e.To,
			Class: e.Class,
			Style: e.Style,
		}
	}

	return layli.Config{
		Layout:         "absolute",
		LayoutAttempts: d.Config.LayoutAttempts,
		Path: layli.ConfigPath{
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
