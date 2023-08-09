package layli

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes   ConfigNodes `yaml:"nodes"`
	Edges   ConfigEdges `yaml:"edges"`
	Spacing int         `yaml:"-"`

	NodeWidth  int `yaml:"width"`
	NodeHeight int `yaml:"height"`
	Border     int `yaml:"border"`
	Margin     int `yaml:"margin"`
}

type ConfigNode struct {
	Id       string `yaml:"id"`
	Contents string `yaml:"contents"`
}

type ConfigNodes []ConfigNode

type ConfigEdge struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type ConfigEdges []ConfigEdge

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
