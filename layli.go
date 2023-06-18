package layli

import "io"

type Diagram struct{}

func NewDiagram(r io.ReadCloser, output func(output string) error) (*Diagram, error) {
	return nil, nil
}

func (d *Diagram) Draw() error {
	return nil
}
