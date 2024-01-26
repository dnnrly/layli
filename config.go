package layli

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigPath struct {
	Attempts int    `yaml:"attempts"`
	Strategy string `yaml:"strategy"`
	Class    string `yaml:"class"`
}

type ConfigStyles map[string]string

func (styles ConfigStyles) toCSS() string {
	keys := make([]string, 0, len(styles))
	for k := range styles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pattern := regexp.MustCompile(`\n[\r\t ]*`)
	css := []string{}
	for _, k := range keys {
		s := pattern.Split(styles[k], -1)
		css = append(css, fmt.Sprintf("%s { %s }", k, strings.Join(s, " ")))
	}

	return strings.Join(css, "\n")
}

type Config struct {
	Layout         string      `yaml:"layout"`
	LayoutAttempts int         `yaml:"layout-attempts"`
	Path           ConfigPath  `yaml:"path"`
	Nodes          ConfigNodes `yaml:"nodes"`
	Edges          ConfigEdges `yaml:"edges"`
	Spacing        int         `yaml:"-"`

	NodeWidth  int `yaml:"width"`
	NodeHeight int `yaml:"height"`
	Border     int `yaml:"border"`
	Margin     int `yaml:"margin"`

	Styles ConfigStyles `yaml:"styles"`
}

type Position struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type ConfigNode struct {
	Id       string   `yaml:"id"`
	Contents string   `yaml:"contents"`
	Position Position `yaml:"position"`
	Class    string   `yaml:"class"`
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
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Class string `yaml:"class"`
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
	if config.Path.Attempts > 10000 {
		return nil, fmt.Errorf("cannot specify more that 10000 path attempts")
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
	if config.Margin > 10 {
		return nil, fmt.Errorf("margin cannot be larger than 10")
	}
	if config.Border == 0 {
		config.Border = 1
	}
	if config.LayoutAttempts == 0 {
		config.LayoutAttempts = 10
	}
	if config.LayoutAttempts > 10000 {
		return nil, fmt.Errorf("cannot specify more that 10000 layout attempts")
	}

	if len(config.Nodes) == 0 {
		return nil, fmt.Errorf("must specify at least 1 node")
	}

	for _, n := range config.Nodes {
		if n.Id == "" {
			return nil, fmt.Errorf("all nodes must have an id")
		}
	}

	for _, e := range config.Edges {
		if e.From == "" || e.To == "" {
			return nil, fmt.Errorf("all edges must have a from and a to")
		}
		if e.From == e.To {
			return nil, fmt.Errorf("edges cannot have the same from and to")
		}

		if config.Nodes.ByID(e.From) == nil || config.Nodes.ByID(e.To) == nil {
			return nil, fmt.Errorf("all edges must have a from and a to that are valid node ids")
		}
	}

	return &config, nil
}
