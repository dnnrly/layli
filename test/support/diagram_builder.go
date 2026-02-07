package support

import "github.com/dnnrly/layli/internal/domain"

// DiagramBuilder provides a fluent API for creating test diagrams.
type DiagramBuilder struct {
	diagram *domain.Diagram
}

// NewDiagram creates a new diagram builder with defaults.
func NewDiagram() *DiagramBuilder {
	return &DiagramBuilder{
		diagram: &domain.Diagram{
			Nodes: []domain.Node{},
			Edges: []domain.Edge{},
			Config: domain.DiagramConfig{
				NodeWidth:      5,
				NodeHeight:     3,
				LayoutType:     domain.LayoutFlowSquare,
				Border:         2,
				Margin:         1,
				Spacing:        1,
				LayoutAttempts: 100,
				PathAttempts:   1000,
				PathStrategy:   "dijkstra",
				Styles:         make(map[string]string),
			},
		},
	}
}

// WithNodeWidth sets the diagram node width.
func (b *DiagramBuilder) WithNodeWidth(w int) *DiagramBuilder {
	b.diagram.Config.NodeWidth = w
	return b
}

// WithNodeHeight sets the diagram node height.
func (b *DiagramBuilder) WithNodeHeight(h int) *DiagramBuilder {
	b.diagram.Config.NodeHeight = h
	return b
}

// WithLayout sets the layout type.
func (b *DiagramBuilder) WithLayout(layout domain.LayoutType) *DiagramBuilder {
	b.diagram.Config.LayoutType = layout
	return b
}

// WithBorder sets the diagram border.
func (b *DiagramBuilder) WithBorder(border int) *DiagramBuilder {
	b.diagram.Config.Border = border
	return b
}

// WithMargin sets the diagram margin.
func (b *DiagramBuilder) WithMargin(margin int) *DiagramBuilder {
	b.diagram.Config.Margin = margin
	return b
}

// AddNode adds a node to the diagram.
func (b *DiagramBuilder) AddNode(id, contents string) *DiagramBuilder {
	b.diagram.Nodes = append(b.diagram.Nodes, domain.Node{
		ID:       id,
		Contents: contents,
		Width:    5,
		Height:   3,
	})
	return b
}

// AddNodeWithDimensions adds a node with specified dimensions.
func (b *DiagramBuilder) AddNodeWithDimensions(id, contents string, width, height int) *DiagramBuilder {
	b.diagram.Nodes = append(b.diagram.Nodes, domain.Node{
		ID:       id,
		Contents: contents,
		Width:    width,
		Height:   height,
	})
	return b
}

// AddEdge adds an edge to the diagram.
func (b *DiagramBuilder) AddEdge(from, to string) *DiagramBuilder {
	b.diagram.Edges = append(b.diagram.Edges, domain.Edge{
		From: from,
		To:   to,
	})
	return b
}

// Build returns the constructed diagram.
func (b *DiagramBuilder) Build() *domain.Diagram {
	return b.diagram
}
