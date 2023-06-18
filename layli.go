package layli

import "io"

type OutputFunc func(output string) error

type Diagram struct {
	output OutputFunc
}

// NewDiagram reads the configuration and parses it in to a Diagram object
func NewDiagram(r io.ReadCloser, output OutputFunc) (*Diagram, error) {
	return &Diagram{
		output: output,
	}, nil
}

// Draw turns the diagram in to an image
func (d *Diagram) Draw() error {
	return d.output("hello world")
}
