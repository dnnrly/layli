package layli

import (
	"io"
	"strings"
	"testing"

	"github.com/dnnrly/layli/pathfinder/djikstra"
	"github.com/stretchr/testify/assert"
)

func nilCreator(start, end djikstra.Point) PathFinder { return nil }

func Test_NewDiagramFromFile_Simple(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
    contents: "Some content here"
  - id: node-2
    contents: "More contents"
`)

	d, err := NewDiagramFromFile(nilCreator, io.NopCloser(r), func(output string) error {
		return nil
	}, false)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(d.config.Nodes))
}

func Test_NewDiagramFromFile_GeneratesSimplestDiagram(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
    contents: "A single box"
`)
	actualOutput := ""

	d, _ := NewDiagramFromFile(nilCreator, io.NopCloser(r), func(output string) error {
		actualOutput = output
		return nil
	}, false)

	assert.NoError(t, d.Draw())
	assert.Contains(t, actualOutput, "A single box")
}
