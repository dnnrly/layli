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

func NewConfigFromFile(r io.Reader) (*Config, error) {
	config := Config{}
	err := yaml.NewDecoder(r).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	config.Spacing = 20

	if config.NodeWidth == 0 {
		config.NodeWidth = 5
	}
	if config.NodeHeight == 0 {
		config.NodeHeight = 3
	}
	if config.Margin == 0 {
		config.Margin = 2
	}
	if config.Border == 0 {
		config.Border = 1
	}

	return &config, nil
}

// NewDiagramFromFile reads the configuration and parses it in to a Diagram object
func NewDiagramFromFile(cf CreateFinder, config *Config, output OutputFunc, showGrid bool) (*Diagram, error) {
	d := Diagram{
		output:   output,
		showGrid: showGrid,
	}
	d.config = *config

	layout, err := NewLayoutFromConfig(cf, d.config)
	if err != nil {
		return nil, err
	}

	d.layout = layout

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
