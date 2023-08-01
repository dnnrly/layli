package layli

import (
	"fmt"
	"io"
	"strings"

	svg "github.com/ajstarks/svgo"
	"gopkg.in/yaml.v3"
)

type OutputFunc func(output string) error

type Diagram struct {
	output   OutputFunc
	config   Config
	layout   *Layout
	showGrid bool
}

// NewDiagramFromFile reads the configuration and parses it in to a Diagram object
func NewDiagramFromFile(cf CreateFinder, r io.ReadCloser, output OutputFunc, showGrid bool) (*Diagram, error) {
	d := Diagram{
		output:   output,
		showGrid: showGrid,
	}
	err := yaml.NewDecoder(r).Decode(&d.config)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	d.config.Spacing = 20

	if d.config.NodeWidth == 0 {
		d.config.NodeWidth = 5
	}
	if d.config.NodeHeight == 0 {
		d.config.NodeHeight = 3
	}
	if d.config.Margin == 0 {
		d.config.Margin = 2
	}
	if d.config.Border == 0 {
		d.config.Border = 1
	}

	d.layout = NewLayoutFromConfig(cf, d.config)

	return &d, nil
}

// Draw turns the diagram in to an image
func (d *Diagram) Draw() error {
	w := strings.Builder{}

	canvas := svg.New(&w)
	canvas.Start(
		(d.layout.LayoutWidth()-1)*d.config.Spacing,
		(d.layout.LayoutHeight()-1)*d.config.Spacing,
		"style=\"background-color: white;\"",
	)
	canvas.Gstyle("text-anchor:middle;font-family:sans;fill:none;stroke:black")
	canvas.Def()
	canvas.Marker("arrow", 10, 5, 7, 7,
		`viewBox="0 0 10 10"`,
		`fill="black"`,
		`orient="auto-start-reverse"`)
	canvas.Path("M 0 0 L 10 5 L 0 10 z")
	canvas.MarkerEnd()
	canvas.DefEnd()

	if d.showGrid {
		d.layout.ShowGrid(canvas, d.config.Spacing)
	}

	d.layout.Draw(canvas, d.config.Spacing)

	canvas.Gend()

	canvas.End()
	return d.output(w.String())
}
