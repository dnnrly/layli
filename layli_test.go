package layli

import (
	"strings"
	"testing"

	"github.com/dnnrly/layli/pathfinder/dijkstra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagram_DrawWithStyleClass(t *testing.T) {
	output := ""
	d := Diagram{
		Output: func(data string) error { output = data; return nil },
		Config: Config{
			Styles: map[string]string{
				".c1": "fill: black; stroke: white;",
				".c2": `stroke: black;
	stroke-width: 2;`},
		},
		Layout:   &Layout{},
		ShowGrid: false,
	}

	err := d.Draw()
	assert.NoError(t, err)
	assert.Contains(t, output, `<style type="text/css">
<![CDATA[
.c1 { fill: black; stroke: white; }
.c2 { stroke: black; stroke-width: 2; }
]]>
</style>
<g`) // Make sure that the style occurs BEFORE the g tag
}

func TestDiagram_DrawWithoutStyleClass(t *testing.T) {
	output := ""
	d := Diagram{
		Output:   func(data string) error { output = data; return nil },
		Layout:   &Layout{},
		ShowGrid: false,
	}

	err := d.Draw()
	assert.NoError(t, err)
	assert.NotContains(t, output, `<style type="text/css"`)
}

func TestDiagram_DrawNodeWithClass(t *testing.T) {
	output := ""
	d := Diagram{
		Output: func(data string) error { output = data; return nil },
		Layout: &Layout{
			Nodes: LayoutNodes{
				NewLayoutNode("node-1", "Text", 10, 10, 20, 20, ".c1"),
			},
		},
		ShowGrid: false,
	}

	err := d.Draw()
	assert.NoError(t, err)
	assert.Regexp(t, `<rect.+id="node-1".+class=".c1".+/>`, output)
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

	t.Run("Styles", func(t *testing.T) {
		check(t, `width: 7
height: 7
margin: 2

nodes:
    - id: a
      contents: Node 1
      style: fill:blue
    - id: b
      contents: Node 2
      class: c1
    - id: c
      contents: Noe 3
    - id: d
      contents: Node 4

edges:
    - from: a
      to: b
      style: stroke:red
    - from: b
      to: c
      class: c2
    - from: c
      to: d

styles:
    c1: fill:red
    c2: stroke-width:2`, "")
	})
}
