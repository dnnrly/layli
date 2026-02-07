package usecases

import "github.com/dnnrly/layli/internal/domain"

// ConfigParser reads configuration files and returns domain diagrams.
// Implementations: YAML parser, JSON parser, etc.
type ConfigParser interface {
	// Parse reads a config file and returns a validated Diagram.
	// Maps to: "Given I have a diagram config 'file.layli'"
	Parse(path string) (*domain.Diagram, error)
}

// LayoutEngine arranges nodes within a diagram.
// Implementations: FlowSquare, TopoSort, Tarjan, Absolute
type LayoutEngine interface {
	// Arrange positions all nodes in the diagram.
	// Maps to: "When I arrange using 'flow-square' layout"
	Arrange(diagram *domain.Diagram) error
}

// Pathfinder calculates edge paths between nodes.
// Implementations: Dijkstra, A*, etc.
type Pathfinder interface {
	// FindPaths calculates paths for all edges in the diagram.
	// Maps to: "And calculate paths for all edges"
	FindPaths(diagram *domain.Diagram) error
}

// Renderer generates output from a positioned diagram.
// Implementations: SVG, PNG, PDF
type Renderer interface {
	// Render writes the diagram to the output path.
	// Maps to: "Then the diagram should be generated"
	Render(diagram *domain.Diagram, outputPath string) error
}

// FileReader abstracts file system reads (for testing).
type FileReader interface {
	Read(path string) ([]byte, error)
}

// FileWriter abstracts file system writes (for testing).
type FileWriter interface {
	Write(path string, data []byte) error
}
