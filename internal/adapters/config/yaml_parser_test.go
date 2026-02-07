package config

import (
	"fmt"
	"testing"

	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFileReader struct {
	files map[string][]byte
}

func (m *mockFileReader) Read(path string) ([]byte, error) {
	data, ok := m.files[path]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return data, nil
}

func newParser(files map[string][]byte) *YAMLParser {
	return NewYAMLParser(&mockFileReader{files: files})
}

func TestYAMLParser_Parse(t *testing.T) {
	t.Run("minimal config with just nodes", func(t *testing.T) {
		parser := newParser(map[string][]byte{
			"test.layli": []byte(`
nodes:
  - id: node-1
    contents: "C1"
  - id: node-2
    contents: "C2"
`),
		})

		diagram, err := parser.Parse("test.layli")
		require.NoError(t, err)

		assert.Len(t, diagram.Nodes, 2)
		assert.Equal(t, "node-1", diagram.Nodes[0].ID)
		assert.Equal(t, "C1", diagram.Nodes[0].Contents)
		assert.Equal(t, "node-2", diagram.Nodes[1].ID)
		assert.Equal(t, "C2", diagram.Nodes[1].Contents)

		assert.Equal(t, 5, diagram.Config.NodeWidth)
		assert.Equal(t, 3, diagram.Config.NodeHeight)
		assert.Equal(t, 1, diagram.Config.Border)
		assert.Equal(t, 2, diagram.Config.Margin)
		assert.Equal(t, 20, diagram.Config.Spacing)
		assert.Equal(t, 20, diagram.Config.PathAttempts)
		assert.Equal(t, 10, diagram.Config.LayoutAttempts)
	})

	t.Run("full config with all fields", func(t *testing.T) {
		parser := newParser(map[string][]byte{
			"full.layli": []byte(`
layout: topo-sort
layout-attempts: 50
width: 10
height: 6
border: 2
margin: 5
path:
  attempts: 100
  strategy: random
  class: path-class
nodes:
  - id: node-1
    contents: "C1"
    position:
      x: 1
      y: 2
    class: my-class
    style: "fill: red"
  - id: node-2
    contents: "C2"
edges:
  - id: e1
    from: node-1
    to: node-2
    class: edge-class
    style: "stroke: blue"
styles:
  ".c1": "fill: black;"
`),
		})

		diagram, err := parser.Parse("full.layli")
		require.NoError(t, err)

		assert.Equal(t, domain.LayoutType("topo-sort"), diagram.Config.LayoutType)
		assert.Equal(t, 50, diagram.Config.LayoutAttempts)
		assert.Equal(t, 10, diagram.Config.NodeWidth)
		assert.Equal(t, 6, diagram.Config.NodeHeight)
		assert.Equal(t, 2, diagram.Config.Border)
		assert.Equal(t, 5, diagram.Config.Margin)
		assert.Equal(t, 100, diagram.Config.PathAttempts)
		assert.Equal(t, "random", diagram.Config.PathStrategy)
		assert.Equal(t, map[string]string{".c1": "fill: black;"}, diagram.Config.Styles)

		require.Len(t, diagram.Nodes, 2)
		assert.Equal(t, "node-1", diagram.Nodes[0].ID)
		assert.Equal(t, "C1", diagram.Nodes[0].Contents)
		assert.Equal(t, domain.Position{X: 1, Y: 2}, diagram.Nodes[0].Position)
		assert.Equal(t, "my-class", diagram.Nodes[0].Class)
		assert.Equal(t, "fill: red", diagram.Nodes[0].Style)
		assert.Equal(t, 10, diagram.Nodes[0].Width)
		assert.Equal(t, 6, diagram.Nodes[0].Height)

		require.Len(t, diagram.Edges, 1)
		assert.Equal(t, "e1", diagram.Edges[0].ID)
		assert.Equal(t, "node-1", diagram.Edges[0].From)
		assert.Equal(t, "node-2", diagram.Edges[0].To)
		assert.Equal(t, "edge-class", diagram.Edges[0].Class)
		assert.Equal(t, "stroke: blue", diagram.Edges[0].Style)
	})

	t.Run("edge ID auto-generation", func(t *testing.T) {
		parser := newParser(map[string][]byte{
			"edges.layli": []byte(`
nodes:
  - id: node-1
  - id: node-2
  - id: node-3
edges:
  - from: node-1
    to: node-2
    class: a-class
  - id: 2-to-3
    from: node-2
    to: node-3
    style: some-style
  - from: node-3
    to: node-1
`),
		})

		diagram, err := parser.Parse("edges.layli")
		require.NoError(t, err)

		require.Len(t, diagram.Edges, 3)
		assert.Equal(t, "edge-1", diagram.Edges[0].ID)
		assert.Equal(t, "a-class", diagram.Edges[0].Class)
		assert.Equal(t, "2-to-3", diagram.Edges[1].ID)
		assert.Equal(t, "some-style", diagram.Edges[1].Style)
		assert.Equal(t, "edge-3", diagram.Edges[2].ID)
	})

	t.Run("default values applied correctly", func(t *testing.T) {
		parser := newParser(map[string][]byte{
			"defaults.layli": []byte(`
nodes:
  - id: a
`),
		})

		diagram, err := parser.Parse("defaults.layli")
		require.NoError(t, err)

		assert.Equal(t, 5, diagram.Config.NodeWidth)
		assert.Equal(t, 3, diagram.Config.NodeHeight)
		assert.Equal(t, 1, diagram.Config.Border)
		assert.Equal(t, 2, diagram.Config.Margin)
		assert.Equal(t, 20, diagram.Config.Spacing)
		assert.Equal(t, 20, diagram.Config.PathAttempts)
		assert.Equal(t, 10, diagram.Config.LayoutAttempts)
	})

	t.Run("styles map is empty when not provided", func(t *testing.T) {
		parser := newParser(map[string][]byte{
			"no-styles.layli": []byte(`
nodes:
  - id: a
`),
		})

		diagram, err := parser.Parse("no-styles.layli")
		require.NoError(t, err)

		assert.NotNil(t, diagram.Config.Styles)
		assert.Empty(t, diagram.Config.Styles)
	})

	t.Run("node dimensions come from config", func(t *testing.T) {
		parser := newParser(map[string][]byte{
			"dims.layli": []byte(`
width: 8
height: 4
nodes:
  - id: a
  - id: b
`),
		})

		diagram, err := parser.Parse("dims.layli")
		require.NoError(t, err)

		assert.Equal(t, 8, diagram.Nodes[0].Width)
		assert.Equal(t, 4, diagram.Nodes[0].Height)
		assert.Equal(t, 8, diagram.Nodes[1].Width)
		assert.Equal(t, 4, diagram.Nodes[1].Height)
	})
}

func TestYAMLParser_Parse_Validation(t *testing.T) {
	check := func(t *testing.T, content string, contained string) {
		t.Helper()
		parser := newParser(map[string][]byte{
			"test.layli": []byte(content),
		})
		_, err := parser.Parse("test.layli")
		assert.Error(t, err)
		assert.ErrorContains(t, err, contained)
	}

	t.Run("bad YAML", func(t *testing.T) {
		check(t, `
nodes:
  - id: node-1
-
  `, "reading config file:")
	})

	t.Run("no nodes", func(t *testing.T) {
		check(t, `
path:
  strategy: random
  attempts: 99
`, "must specify at least 1 node")
	})

	t.Run("node without ID", func(t *testing.T) {
		check(t, `
nodes:
  - 00: 0
  - id: b
  - id: c
  - id: d
  - id: e
edges:
  - from: a
    to: b
  - from: b
    to: c
`, "all nodes must have an id")
	})

	t.Run("edge without from and to", func(t *testing.T) {
		check(t, `
nodes:
  - id: a
edges:
  - 0:
`, "all edges must have a from and a to")
	})

	t.Run("edge referencing non-existent node", func(t *testing.T) {
		check(t, `
nodes:
  - id: a
edges:
  - from: 0
    to: 00
`, "all edges must have a from and a to that are valid node ids")
	})

	t.Run("edge with same from and to", func(t *testing.T) {
		check(t, `
nodes:
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
    to: c
`, "edges cannot have the same from and to")
	})

	t.Run("path attempts too high", func(t *testing.T) {
		check(t, `
path:
  strategy: random
  attempts: 1e12
nodes:
  - id: a
`, "cannot specify more that 10000 path attempts")
	})

	t.Run("layout attempts too high", func(t *testing.T) {
		check(t, `
nodes:
  - id: a
layout: random-shortest-square
layout-attempts: 1e12
`, "cannot specify more that 10000 layout attempts")
	})

	t.Run("margin too big", func(t *testing.T) {
		check(t, `
margin: 20
nodes:
  - id: a
`, "margin cannot be larger than 10")
	})

	t.Run("file not found", func(t *testing.T) {
		parser := newParser(map[string][]byte{})
		_, err := parser.Parse("missing.layli")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "reading config file:")
	})

	t.Run("non-number for layout attempts", func(t *testing.T) {
		check(t, `
nodes:
  - id: a
layout: random-shortest-square
layout-attempts: no-a-number
`, "reading config file: yaml: unmarshal errors")
	})
}
