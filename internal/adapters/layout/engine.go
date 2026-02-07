package layout

import (
	"fmt"

	"github.com/dnnrly/layli/internal/adapters"
	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/layout"
)

type LayoutAdapter struct{}

func NewLayoutAdapter() *LayoutAdapter {
	return &LayoutAdapter{}
}

func (a *LayoutAdapter) Arrange(diagram *domain.Diagram) error {
	cfg := adapters.ToLayoutConfig(diagram)

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
			return fmt.Errorf("layout engine failed to arrange node: %s", diagram.Nodes[i].ID)
		}
		diagram.Nodes[i].Position.X = ln.Left()
		diagram.Nodes[i].Position.Y = ln.Top()
		diagram.Nodes[i].Width = ln.Width()
		diagram.Nodes[i].Height = ln.Height()
	}

	return nil
}

// selectArranger maps domain layout types to their layout package implementations.
// This is the adapter layer that bridges the domain (domain/diagram.go) with the
// layout implementation (layout/arrangements.go).
//
// When adding a new layout algorithm:
// 1. Add LayoutXxx constant to internal/domain/diagram.go
// 2. Implement LayoutXxx() function in layout/arrangements.go
// 3. Register in selectArrangement() in layout/arrangements.go
// 4. Add case here mapping domain.LayoutXxx to layout.LayoutXxx
// 5. Add tests to layout/arrangements_test.go
//
// See CONTRIBUTING_LAYOUTS.md for the complete guide.
func selectArranger(lt domain.LayoutType) (layout.LayoutArrangementFunc, error) {
	switch lt {
	case "", domain.LayoutFlowSquare:
		return layout.LayoutFlowSquare, nil
	case domain.LayoutTopoSort:
		return layout.LayoutTopologicalSort, nil
	case domain.LayoutTarjan:
		return layout.LayoutTarjan, nil
	case domain.LayoutRandomShortest:
		return layout.LayoutRandomShortestSquare, nil
	case domain.LayoutAbsolute:
		return layout.LayoutAbsolute, nil
	default:
		return nil, fmt.Errorf("unknown layout type: %s", lt)
	}
}
