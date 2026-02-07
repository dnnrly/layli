package composition

import (
	"github.com/dnnrly/layli/internal/adapters/config"
	"github.com/dnnrly/layli/internal/adapters/filesystem"
	"github.com/dnnrly/layli/internal/adapters/layout"
	"github.com/dnnrly/layli/internal/adapters/pathfinding"
	"github.com/dnnrly/layli/internal/adapters/rendering"
	"github.com/dnnrly/layli/internal/usecases"
)

// NewGenerateDiagram wires all adapters together and returns
// a ready-to-use GenerateDiagram use case.
func NewGenerateDiagram(showGrid bool) *usecases.GenerateDiagram {
	reader := filesystem.NewOSFileReader()
	writer := filesystem.NewOSFileWriter()

	parser := config.NewYAMLParser(reader)
	layoutEngine := layout.NewLayoutAdapter()
	pathfinder := pathfinding.NewDijkstraPathfinder()
	renderer := rendering.NewSVGRenderer(writer, showGrid)

	return usecases.NewGenerateDiagram(parser, layoutEngine, pathfinder, renderer)
}
