package layli

import (
	"fmt"
	"io"
	"math"
	"strings"

	svg "github.com/ajstarks/svgo"
	"gopkg.in/yaml.v3"
)

type OutputFunc func(output string) error

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
	root := math.Sqrt(float64(len(d.config.Nodes)))
	size := int(math.Round(root))

	if size < int(root) {
		size++
	}
	if len(d.config.Nodes) < 4 {
		size = 2
	}

	gridSpacing := 150
	nodeWidth := 100
	nodeHeight := 80

	w := strings.Builder{}
	canvas := svg.New(&w)
	canvas.Start(gridSpacing*(size+1), gridSpacing*(size+1), "style=\"background-color: white;\"")
	canvas.Gstyle("text-anchor:middle;font-family:sans;fill:none;stroke:black")

	pos := 0
	for y := 0; y < size && pos < len(d.config.Nodes); y++ {
		for x := 0; x < size && pos < len(d.config.Nodes); x++ {
			d.config.Nodes[pos].X = x + 1
			d.config.Nodes[pos].Y = y + 1
			d.config.Nodes[pos].Width = nodeWidth
			d.config.Nodes[pos].Height = nodeHeight

			pos++
		}
	}

	for _, n := range d.config.Nodes {
		n.Draw(canvas, gridSpacing)
	}
	canvas.Gend()

	canvas.End()
	return d.output(w.String())
}
