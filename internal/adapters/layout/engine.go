package layout

import (
	"fmt"

	"github.com/dnnrly/layli/internal/domain"
	layoutpkg "github.com/dnnrly/layli/internal/layout"
)

type LayoutAdapter struct{}

func NewLayoutAdapter() *LayoutAdapter {
	return &LayoutAdapter{}
}

func (a *LayoutAdapter) Arrange(diagram *domain.Diagram) error {
	cfg := toRootConfig(diagram)

	arranger, err := selectArranger(diagram.Config.LayoutType)
	if err != nil {
		return err
	}

	nodes, err := arranger(&cfg)
	if err != nil {
		return fmt.Errorf("arranging nodes: %w", err)
	}

	for i := range diagram.Nodes {
		ln := nodes.ByID(diagram.Nodes[i].ID)
		if ln == nil {
			continue
		}
		diagram.Nodes[i].Position.X = ln.Left()
		diagram.Nodes[i].Position.Y = ln.Top()
		diagram.Nodes[i].Width = ln.Width()
		diagram.Nodes[i].Height = ln.Height()
	}

	return nil
}

func toRootConfig(d *domain.Diagram) layoutpkg.Config {
	nodes := make(layoutpkg.ConfigNodes, len(d.Nodes))
	for i, n := range d.Nodes {
		nodes[i] = layoutpkg.ConfigNode{
			Id:       n.ID,
			Contents: n.Contents,
			Position: layoutpkg.Position{
				X: n.Position.X,
				Y: n.Position.Y,
			},
			Class: n.Class,
			Style: n.Style,
		}
	}

	edges := make(layoutpkg.ConfigEdges, len(d.Edges))
	for i, e := range d.Edges {
		edges[i] = layoutpkg.ConfigEdge{
			ID:    e.ID,
			From:  e.From,
			To:    e.To,
			Class: e.Class,
			Style: e.Style,
		}
	}

	return layoutpkg.Config{
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

func selectArranger(lt domain.LayoutType) (layoutpkg.LayoutArrangementFunc, error) {
	switch lt {
	case "", domain.LayoutFlowSquare:
		return layoutpkg.LayoutFlowSquare, nil
	case domain.LayoutTopoSort:
		return layoutpkg.LayoutTopologicalSort, nil
	case domain.LayoutTarjan:
		return layoutpkg.LayoutTarjan, nil
	case domain.LayoutRandomShortest:
		return layoutpkg.LayoutRandomShortestSquare, nil
	case domain.LayoutAbsolute:
		return layoutpkg.LayoutAbsolute, nil
	default:
		return nil, fmt.Errorf("unknown layout type: %s", lt)
	}
}
