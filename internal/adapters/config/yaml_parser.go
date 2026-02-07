package config

import (
	"fmt"

	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/internal/usecases"
	"gopkg.in/yaml.v3"
)

type configPath struct {
	Attempts int    `yaml:"attempts,omitempty"`
	Strategy string `yaml:"strategy,omitempty"`
	Class    string `yaml:"class,omitempty"`
}

type configPosition struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type configNode struct {
	ID       string         `yaml:"id"`
	Contents string         `yaml:"contents"`
	Position configPosition `yaml:"position,omitempty"`
	Class    string         `yaml:"class,omitempty"`
	Style    string         `yaml:"style,omitempty"`
}

type configEdge struct {
	ID    string `yaml:"id,omitempty"`
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Class string `yaml:"class,omitempty"`
	Style string `yaml:"style,omitempty"`
}

type configFile struct {
	Layout         string            `yaml:"layout,omitempty"`
	LayoutAttempts int               `yaml:"layout-attempts,omitempty"`
	Path           configPath        `yaml:"path,omitempty"`
	Nodes          []configNode      `yaml:"nodes"`
	Edges          []configEdge      `yaml:"edges"`
	NodeWidth      int               `yaml:"width"`
	NodeHeight     int               `yaml:"height"`
	Border         int               `yaml:"border"`
	Margin         int               `yaml:"margin"`
	Styles         map[string]string `yaml:"styles,omitempty"`
}

type YAMLParser struct {
	reader usecases.FileReader
}

func NewYAMLParser(reader usecases.FileReader) *YAMLParser {
	return &YAMLParser{reader: reader}
}

func (p *YAMLParser) Parse(path string) (*domain.Diagram, error) {
	data, err := p.reader.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg configFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	applyDefaults(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return toDomain(&cfg), nil
}

func applyDefaults(cfg *configFile) {
	if cfg.Path.Attempts == 0 {
		cfg.Path.Attempts = 20
	}
	if cfg.NodeWidth == 0 {
		cfg.NodeWidth = 5
	}
	if cfg.NodeHeight == 0 {
		cfg.NodeHeight = 3
	}
	if cfg.Margin == 0 {
		cfg.Margin = 2
	}
	if cfg.Border == 0 {
		cfg.Border = 1
	}
	if cfg.LayoutAttempts == 0 {
		cfg.LayoutAttempts = 10
	}

	for i, e := range cfg.Edges {
		if e.ID == "" {
			cfg.Edges[i].ID = fmt.Sprintf("edge-%d", i+1)
		}
	}
}

func validate(cfg *configFile) error {
	if cfg.Path.Attempts > 10000 {
		return fmt.Errorf("cannot specify more that 10000 path attempts")
	}
	if cfg.Margin > 10 {
		return fmt.Errorf("margin cannot be larger than 10")
	}
	if cfg.LayoutAttempts > 10000 {
		return fmt.Errorf("cannot specify more that 10000 layout attempts")
	}
	if len(cfg.Nodes) == 0 {
		return fmt.Errorf("must specify at least 1 node")
	}

	for _, n := range cfg.Nodes {
		if n.ID == "" {
			return fmt.Errorf("all nodes must have an id")
		}
	}

	nodeIDs := make(map[string]bool, len(cfg.Nodes))
	for _, n := range cfg.Nodes {
		nodeIDs[n.ID] = true
	}

	for _, e := range cfg.Edges {
		if e.From == "" || e.To == "" {
			return fmt.Errorf("all edges must have a from and a to")
		}
		if e.From == e.To {
			return fmt.Errorf("edges cannot have the same from and to")
		}
		if !nodeIDs[e.From] || !nodeIDs[e.To] {
			return fmt.Errorf("all edges must have a from and a to that are valid node ids")
		}
	}

	return nil
}

func toDomain(cfg *configFile) *domain.Diagram {
	nodes := make([]domain.Node, len(cfg.Nodes))
	for i, n := range cfg.Nodes {
		nodes[i] = domain.Node{
			ID:       n.ID,
			Contents: n.Contents,
			Position: domain.Position{
				X: n.Position.X,
				Y: n.Position.Y,
			},
			Width:  cfg.NodeWidth,
			Height: cfg.NodeHeight,
			Class:  n.Class,
			Style:  n.Style,
		}
	}

	edges := make([]domain.Edge, len(cfg.Edges))
	for i, e := range cfg.Edges {
		edges[i] = domain.Edge{
			ID:    e.ID,
			From:  e.From,
			To:    e.To,
			Class: e.Class,
			Style: e.Style,
		}
	}

	styles := cfg.Styles
	if styles == nil {
		styles = map[string]string{}
	}

	return &domain.Diagram{
		Nodes: nodes,
		Edges: edges,
		Config: domain.DiagramConfig{
			LayoutType:     domain.LayoutType(cfg.Layout),
			LayoutAttempts: cfg.LayoutAttempts,
			NodeWidth:      cfg.NodeWidth,
			NodeHeight:     cfg.NodeHeight,
			Border:         cfg.Border,
			Margin:         cfg.Margin,
			Spacing:        20,
			PathAttempts:   cfg.Path.Attempts,
			PathStrategy:   cfg.Path.Strategy,
			Styles:         styles,
		},
	}
}
