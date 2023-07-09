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
func NewDiagramFromFile(r io.ReadCloser, output OutputFunc, showGrid bool) (*Diagram, error) {
	d := Diagram{
		output:   output,
		showGrid: showGrid,
	}
	err := yaml.NewDecoder(r).Decode(&d.config)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	d.layout = NewLayoutFromConfig(d.config)

	return &d, nil
}

// Draw turns the diagram in to an image
func (d *Diagram) Draw() error {
	w := strings.Builder{}

	pathUnitSize := 20

	canvas := svg.New(&w)
	canvas.Start(
		d.layout.LayoutWidth()*pathUnitSize,
		d.layout.LayoutHeight()*pathUnitSize,
		"style=\"background-color: white;\"",
	)
	canvas.Gstyle("text-anchor:middle;font-family:sans;fill:none;stroke:black")

	if d.showGrid {
		d.layout.ShowGrid(canvas, pathUnitSize)
	}

	d.layout.Draw(canvas, pathUnitSize)
	canvas.Gend()

	canvas.End()
	return d.output(w.String())
}
