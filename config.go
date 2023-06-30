package layli

type Config struct {
	Nodes ConfigNodes `yaml:"nodes"`
}

type ConfigNode struct {
	Id       string `yaml:"id"`
	Contents string `yaml:"contents"`
}

type ConfigNodes []ConfigNode
