package rendering

import (
	"fmt"

	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/internal/usecases"
	"github.com/dnnrly/layli/layout"
)

var _ usecases.Renderer = (*SVGRenderer)(nil)

type SVGRenderer struct {
	writer   usecases.FileWriter
	showGrid bool
}

func NewSVGRenderer(writer usecases.FileWriter, showGrid bool) *SVGRenderer {
	return &SVGRenderer{writer: writer, showGrid: showGrid}
}

func (r *SVGRenderer) Render(diagram *domain.Diagram, outputPath string) error {
	cfg := buildConfig(diagram)
	nodes := buildLayoutNodes(diagram)
	paths := buildLayoutPaths(diagram)

	layoutObj := layout.NewLayout(
		nodes, paths,
		diagram.Config.NodeWidth,
		diagram.Config.NodeHeight,
		diagram.Config.Margin,
		diagram.Config.Border,
		diagram.Config.Spacing,
	)

	var svgOutput string
	rootDiagram := layout.Diagram{
		Output: func(output string) error {
			svgOutput = output
			return nil
		},
		Config:   cfg,
		Layout:   layoutObj,
		ShowGrid: r.showGrid,
	}

	if err := rootDiagram.Draw(); err != nil {
		return fmt.Errorf("rendering SVG: %w", err)
	}

	return r.writer.Write(outputPath, []byte(svgOutput))
}

func buildConfig(diagram *domain.Diagram) layout.Config {
	styles := layout.ConfigStyles{}
	for k, v := range diagram.Config.Styles {
		styles[k] = v
	}

	return layout.Config{
		Layout:     string(diagram.Config.LayoutType),
		NodeWidth:  diagram.Config.NodeWidth,
		NodeHeight: diagram.Config.NodeHeight,
		Border:     diagram.Config.Border,
		Margin:     diagram.Config.Margin,
		Spacing:    diagram.Config.Spacing,
		Styles:     styles,
	}
}

func buildLayoutNodes(diagram *domain.Diagram) layout.LayoutNodes {
	nodes := make(layout.LayoutNodes, len(diagram.Nodes))
	for i, n := range diagram.Nodes {
		width := n.Width
		if width == 0 {
			width = diagram.Config.NodeWidth
		}
		height := n.Height
		if height == 0 {
			height = diagram.Config.NodeHeight
		}
		nodes[i] = layout.NewLayoutNode(
			n.ID, n.Contents,
			n.Position.X, n.Position.Y,
			width, height,
			n.Class, n.Style,
		)
	}
	return nodes
}

func buildLayoutPaths(diagram *domain.Diagram) layout.LayoutPaths {
	var paths layout.LayoutPaths
	for _, e := range diagram.Edges {
		// Skip edges without calculated paths.
		// All edges should have paths after pathfinding, so this shouldn't happen in normal operation.
		if e.Path == nil {
			continue
		}
		points := make(layout.Points, len(e.Path.Points))
		for j, p := range e.Path.Points {
			points[j] = layout.Point{X: float64(p.X), Y: float64(p.Y)}
		}
		paths = append(paths, layout.LayoutPath{
			ID:     e.ID,
			From:   e.From,
			To:     e.To,
			Points: points,
			Class:  e.Class,
			Style:  e.Style,
		})
	}
	return paths
}
