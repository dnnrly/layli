package domain

import "fmt"

// LayoutType enumerates available layout algorithms.
// When adding a new layout algorithm, add the constant here, then:
// 1. Implement LayoutXxx function in layout/arrangements.go
// 2. Register in selectArrangement() in layout/arrangements.go
// 3. Register in selectArranger() in internal/adapters/layout/engine.go
// 4. Add tests to layout/arrangements_test.go
// See CONTRIBUTING_LAYOUTS.md for detailed instructions.
type LayoutType string

const (
	LayoutFlowSquare       LayoutType = "flow-square"
	LayoutTopoSort         LayoutType = "topo-sort"
	LayoutTarjan           LayoutType = "tarjan"
	LayoutAbsolute         LayoutType = "absolute"
	LayoutRandomShortest   LayoutType = "random-shortest-square"
)

// DiagramConfig holds diagram-level settings.
type DiagramConfig struct {
	LayoutType      LayoutType
	LayoutAttempts  int
	NodeWidth       int
	NodeHeight      int
	Border          int
	Margin          int
	Spacing         int
	PathAttempts    int
	PathStrategy    string
	Styles          map[string]string
}

// Diagram represents the complete diagram specification.
// Maps to: "Given a diagram with..." in feature files.
type Diagram struct {
	Nodes  []Node
	Edges  []Edge
	Config DiagramConfig
}

// Validate ensures diagram invariants are met.
func (d *Diagram) Validate() error {
	if len(d.Nodes) == 0 {
		return fmt.Errorf("must specify at least 1 node")
	}

	// Validate each node
	for _, n := range d.Nodes {
		if err := n.Validate(); err != nil {
			return err
		}
	}

	// Validate each edge
	for _, e := range d.Edges {
		if err := e.Validate(); err != nil {
			return err
		}

		// Ensure from and to nodes exist
		fromFound := false
		toFound := false
		for _, n := range d.Nodes {
			if n.ID == e.From {
				fromFound = true
			}
			if n.ID == e.To {
				toFound = true
			}
		}

		if !fromFound {
			return fmt.Errorf("edge references non-existent node: %s", e.From)
		}
		if !toFound {
			return fmt.Errorf("edge references non-existent node: %s", e.To)
		}
	}

	// Validate config
	if d.Config.NodeWidth <= 0 || d.Config.NodeHeight <= 0 {
		return fmt.Errorf("node dimensions must be positive")
	}

	if d.Config.Margin < 0 || d.Config.Margin > 10 {
		return fmt.Errorf("margin must be between 0 and 10")
	}

	if d.Config.PathAttempts <= 0 || d.Config.PathAttempts > 10000 {
		return fmt.Errorf("path attempts must be between 1 and 10000")
	}

	if d.Config.LayoutAttempts <= 0 || d.Config.LayoutAttempts > 10000 {
		return fmt.Errorf("layout attempts must be between 1 and 10000")
	}

	return nil
}
