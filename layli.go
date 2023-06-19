package layli

import (
	"fmt"
	"io"

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
	return d.output("hello world")
}
