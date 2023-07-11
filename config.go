package layli

type Config struct {
	Nodes   ConfigNodes `yaml:"nodes"`
	Edges   ConfigEdges `yaml:"edges"`
	Spacing int         `yaml:"-"`
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
