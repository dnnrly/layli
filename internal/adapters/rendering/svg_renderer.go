package rendering

import (
	"fmt"

	"github.com/dnnrly/layli"
	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/internal/usecases"
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

	layout := layli.NewLayout(
		nodes, paths,
		diagram.Config.NodeWidth,
		diagram.Config.NodeHeight,
		diagram.Config.Margin,
		diagram.Config.Border,
		diagram.Config.Spacing,
	)

	var svgOutput string
	rootDiagram := layli.Diagram{
		Output: func(output string) error {
			svgOutput = output
			return nil
		},
		Config:   cfg,
		Layout:   layout,
		ShowGrid: r.showGrid,
	}

	if err := rootDiagram.Draw(); err != nil {
		return fmt.Errorf("rendering SVG: %w", err)
	}

	return r.writer.Write(outputPath, []byte(svgOutput))
}

func buildConfig(diagram *domain.Diagram) layli.Config {
	styles := layli.ConfigStyles{}
	for k, v := range diagram.Config.Styles {
		styles[k] = v
	}

	return layli.Config{
		Layout:     string(diagram.Config.LayoutType),
		NodeWidth:  diagram.Config.NodeWidth,
		NodeHeight: diagram.Config.NodeHeight,
		Border:     diagram.Config.Border,
		Margin:     diagram.Config.Margin,
		Spacing:    diagram.Config.Spacing,
		Styles:     styles,
	}
}

func buildLayoutNodes(diagram *domain.Diagram) layli.LayoutNodes {
	nodes := make(layli.LayoutNodes, len(diagram.Nodes))
	for i, n := range diagram.Nodes {
		width := n.Width
		if width == 0 {
			width = diagram.Config.NodeWidth
		}
		height := n.Height
		if height == 0 {
			height = diagram.Config.NodeHeight
		}
		nodes[i] = layli.NewLayoutNode(
			n.ID, n.Contents,
			n.Position.X, n.Position.Y,
			width, height,
			n.Class, n.Style,
		)
	}
	return nodes
}

func buildLayoutPaths(diagram *domain.Diagram) layli.LayoutPaths {
	var paths layli.LayoutPaths
	for _, e := range diagram.Edges {
		if e.Path == nil {
			continue
		}
		points := make(layli.Points, len(e.Path.Points))
		for j, p := range e.Path.Points {
			points[j] = layli.Point{X: float64(p.X), Y: float64(p.Y)}
		}
		paths = append(paths, layli.LayoutPath{
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
