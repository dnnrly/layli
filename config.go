package layli

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type ConfigPath struct {
	Attempts int    `yaml:"attempts"`
	Strategy string `yaml:"strategy"`
}

type Config struct {
	Layout  string      `yaml:"layout"`
	Path    ConfigPath  `yaml:"path"`
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

func (nodes ConfigNodes) ByID(id string) *ConfigNode {
	for _, n := range nodes {
		if n.Id == id {
			return &n
		}
	}
	return nil
}

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

	if config.Path.Attempts == 0 {
		config.Path.Attempts = 20
	}

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
