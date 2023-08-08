package layli

import (
	"strings"
	"testing"

	"github.com/dnnrly/layli/pathfinder/dijkstra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nilCreator(start, end dijkstra.Point) PathFinder { return nil }

func TestNewConfigFromFile(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
    contents: "C1"
  - id: node-2
    contents: "C2"
`)

	config, err := NewConfigFromFile(r)
	require.NoError(t, err)
	assert.Equal(t, Config{
		Nodes: ConfigNodes{
			ConfigNode{
				Id:       "node-1",
				Contents: "C1",
			},
			ConfigNode{
				Id:       "node-2",
				Contents: "C2",
			},
		},
		Spacing:    20,
		Border:     1,
		Margin:     2,
		NodeWidth:  5,
		NodeHeight: 3,
	}, *config)
}

func TestNewConfigFromFile_FailsOnBadYaml(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
-
  `)

	_, err := NewConfigFromFile(r)
	require.Error(t, err)
}

func Test_NewDiagramFromFile_Simple(t *testing.T) {
	d, err := NewDiagramFromFile(nilCreator, &Config{
		Nodes: ConfigNodes{
			ConfigNode{
				Id:       "node-1",
				Contents: "Some content here",
			},
			ConfigNode{
				Id:       "node-2",
				Contents: "More contents",
			},
		},
	}, func(output string) error {
		return nil
	}, false)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(d.config.Nodes))
}

func Test_NewDiagramFromFile_GeneratesSimplestDiagram(t *testing.T) {
	actualOutput := ""

	d, _ := NewDiagramFromFile(nilCreator, &Config{
		Nodes: ConfigNodes{
			ConfigNode{
				Id:       "node-1",
				Contents: "A single box",
			},
		},
	}, func(output string) error {
		actualOutput = output
		return nil
	}, false)

	assert.NoError(t, d.Draw())
	assert.Contains(t, actualOutput, "A single box")
}
