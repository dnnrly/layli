package layli

import (
	"strings"
	"testing"

	"github.com/dnnrly/layli/pathfinder/dijkstra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		Path: ConfigPath{
			Attempts: 20,
		},
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
		LayoutAttempts: 10,
		Spacing:        20,
		Border:         1,
		Margin:         2,
		NodeWidth:      5,
		NodeHeight:     3,
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

func TestLayliFullFlow(t *testing.T) {
	check := func(t *testing.T, input string, contains string) {
		config, err := NewConfigFromFile(strings.NewReader(input))
		require.NoError(t, err)

		layout, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
			return dijkstra.NewPathFinder(start, end)
		}, config)

		if err != nil {
			require.ErrorContains(t, err, contains)
		}
		d := Diagram{
			Output:   func(data string) error { return nil },
			ShowGrid: false,
			Config:   *config,
			Layout:   layout,
		}
		err = d.Draw()
		if err != nil {
			require.ErrorContains(t, err, contains)
		}

		if contains == "" {
			assert.NoError(t, err)
		}
	}

	t.Run("Normal random path", func(t *testing.T) {
		check(t, `path:
  strategy: random
  attempts: 10
nodes:
  - id: a
  - id: b
edges:
  - from: a
    to: b`, "")
	})

	t.Run("Long running", func(t *testing.T) {
		check(t, `width: 7
height: 7
margin: 8

nodes:
    - id: a
      contents: Node 1
    - id: b
      cojtents: Node 2
    - id: c
      contents: Noe 3
    - id: d
      contents: Node 4

edges:
    - from: a
      to: b
    - from: b
      to: c
    - from: c
      to: d
    - from: d
      to: a `, "")
	})
}
