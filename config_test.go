package layli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigStyles_toCSS(t *testing.T) {
	styles := ConfigStyles{
		".c1": "fill: black; stroke: white;",
		".c2": `stroke: black;
	stroke-width: 2;`,
		"#id-1": "fill: red;",
		"rect":  "stroke-width: 5;",
	}
	expected := `#id-1 { fill: red; }
.c1 { fill: black; stroke: white; }
.c2 { stroke: black; stroke-width: 2; }
rect { stroke-width: 5; }`

	assert.Equal(t, expected, styles.toCSS())
}

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

func TestConfig_validate(t *testing.T) {
	check := func(t *testing.T, config string, contained string) {
		_, err := NewConfigFromFile(strings.NewReader(config))

		assert.Error(t, err)
		assert.ErrorContains(t, err, contained)
	}

	t.Run("Edges have from and to", func(t *testing.T) {
		check(t, `nodes:
  - id: a
edges:
  - 0:`, "all edges must have a from and a to")
	})

	t.Run("Nodes require IDs", func(t *testing.T) {
		check(t, `nodes:
  - id: a
edges:
    - from: 0
      to: 00`, "all edges must have a from and a to that are valid node ids")
	})

	t.Run("Non-number for layout attempts", func(t *testing.T) {
		check(t, `nodes:
  - id: a
layout: random-shortest-square
layout-attempts: no-a-number`, "reading config file: yaml: unmarshal errors")
	})

	t.Run("Layout attempts too big", func(t *testing.T) {
		check(t, `nodes:
- id: a
layout: random-shortest-square
layout-attempts: 1e12`, "cannot specify more that 10000 layout attempts")
	})

	t.Run("Nodes require ID", func(t *testing.T) {
		check(t, `nodes:
    - 00: 0
    - id: b
    - id: c
    - id: d
    - id: e
edges:
    - from: a
      to: b
    - from: b
      to: c`, "all nodes must have an id")
	})

	t.Run("Edges cannot have same from and to", func(t *testing.T) {
		check(t, `nodes:
    - id: a
    - id: b
    - id: c
    - id: d
    - id: e
edges:
  - from: a
    to: b
  - from: a
    to: a
  - from: b
    to: c`, "edges cannot have the same from and to")
	})

	t.Run("Require at least 1 node", func(t *testing.T) {
		check(t, `path:
  strategy: random
  attempts: 99`, "must specify at least 1 node")
	})

	t.Run("Too many path attempt", func(t *testing.T) {
		check(t, `path:
  strategy: random
  attempts: 1e12
  nodes:
    - id: a`, "cannot specify more that 10000 path attempts")
	})

	t.Run("Margin too big", func(t *testing.T) {
		check(t, `margin: 20
nodes:
  - id: a`, "margin cannot be larger than 10")
	})

	t.Run("Styles", func(t *testing.T) {
		t.Run("Missing node class definition", func(t *testing.T) {
			check(t, `
nodes:
    - id: a
      contents: Node 1
      style: fill:blue
    - id: b
      contents: Node 2
      class: c1
    - id: c
      contents: Node 3
      class: missing
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
    c2: stroke-width:2`, "node c uses unknown class 'missing'")
		})
		t.Run("Missing path class definition", func(t *testing.T) {
			check(t, `
nodes:
    - id: a
      contents: Node 1
    - id: b
      contents: Node 2

edges:
    - from: a
      to: b
      class: missing

styles:
    c1: fill:red
    c2: stroke-width:2`, "path a to b uses unknown class 'missing'")
		})
	})
}
