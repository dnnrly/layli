package layli

import (
	"fmt"
	"io"
	"strings"

	svg "github.com/ajstarks/svgo"
	"gopkg.in/yaml.v3"
)

type OutputFunc func(output string) error

type Node struct {
	Id       string `yaml:"id"`
	Contents string `yaml:"contents"`
}

type DiagramConfig struct {
	Nodes []Node `yaml:"nodes"`
}

type Diagram struct {
	output OutputFunc
	config DiagramConfig
}

// NewDiagramFromFile reads the configuration and parses it in to a Diagram object
func NewDiagramFromFile(r io.ReadCloser, output OutputFunc) (*Diagram, error) {
	d := Diagram{
		output: output,
	}
	err := yaml.NewDecoder(r).Decode(&d.config)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return &d, nil
}

// Draw turns the diagram in to an image
func (d *Diagram) Draw() error {
	width := 100
	height := 100

	w := strings.Builder{}
	canvas := svg.New(&w)
	canvas.Start(width, height, "style=\"background-color: white;\"")
	canvas.Gstyle("text-anchor:middle;font-family:sans;fill:none;stroke:black")
	canvas.Roundrect(
		width/10, height/10,
		(width/10)*8, (height/10)*8,
		3, 3,
	)
	canvas.Textspan(
		width/2,
		height/2,
		d.config.Nodes[0].Contents,
		"font-size:10px",
	)
	canvas.TextEnd()
	canvas.Gend()

	canvas.End()
	return d.output(w.String())
}
