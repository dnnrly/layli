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

type Node struct {
	Id       string `yaml:"id"`
	Contents string `yaml:"contents"`

	X int `yaml:"-"`
	Y int `yaml:"-"`
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
	root := math.Sqrt(float64(len(d.config.Nodes)))
	size := int(math.Round(root))

	if size < int(root) {
		size++
	}
	if len(d.config.Nodes) < 4 {
		size = 2
	}

	nodeWidth := 100
	nodeHeight := 80

	w := strings.Builder{}
	canvas := svg.New(&w)
	canvas.Start(nodeWidth*size, nodeHeight*size, "style=\"background-color: white;\"")
	canvas.Gstyle("text-anchor:middle;font-family:sans;fill:none;stroke:black")

	pos := 0
	for y := 0; y < size && pos < len(d.config.Nodes); y++ {
		for x := 0; x < size && pos < len(d.config.Nodes); x++ {
			fmt.Printf("x=%d y=%d pos=%d\n", x, y, pos)
			d.config.Nodes[pos].X = x
			d.config.Nodes[pos].Y = y
			pos++
		}
	}

	for i, n := range d.config.Nodes {
		canvas.Roundrect(
			(nodeWidth*n.X)+(nodeWidth/10), (nodeHeight*n.Y)+nodeHeight/10,
			(nodeWidth/10)*8, (nodeHeight/10)*8,
			3, 3,
			fmt.Sprintf(`id="node%0d-rect"`, i+1),
		)
		canvas.Textspan(
			(nodeWidth*n.X)+(nodeWidth/2),
			(nodeHeight*n.Y)+(nodeHeight/2),
			n.Contents,
			fmt.Sprintf(`id="node%0d-text"`, i+1),
			"font-size:10px",
		)
		canvas.TextEnd()
	}
	canvas.Gend()

	canvas.End()
	return d.output(w.String())
}
