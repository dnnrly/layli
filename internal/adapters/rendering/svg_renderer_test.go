package rendering

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dnnrly/layli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFileWriter struct {
	written map[string][]byte
	err     error
}

func (m *mockFileWriter) Write(path string, data []byte) error {
	if m.err != nil {
		return m.err
	}
	m.written[path] = data
	return nil
}

func newTestDiagram(nodes []domain.Node, edges []domain.Edge) *domain.Diagram {
	return &domain.Diagram{
		Nodes: nodes,
		Edges: edges,
		Config: domain.DiagramConfig{
			NodeWidth:  5,
			NodeHeight: 3,
			Border:     1,
			Margin:     2,
			Spacing:    20,
		},
	}
}

func TestSVGRenderer_Render(t *testing.T) {
	t.Run("renders simple diagram with two nodes", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{ID: "node-1", Contents: "First", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "node-2", Contents: "Second", Position: domain.Position{X: 10, Y: 3}, Width: 5, Height: 3},
			},
			nil,
		)

		err := renderer.Render(diagram, "output.svg")
		require.NoError(t, err)

		svg := string(writer.written["output.svg"])
		assert.Contains(t, svg, `id="node-1"`)
		assert.Contains(t, svg, `id="node-2"`)
		assert.Contains(t, svg, "First")
		assert.Contains(t, svg, "Second")
		assert.Contains(t, svg, "<svg")
	})

	t.Run("renders with edges and paths", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{ID: "a", Contents: "A", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "b", Contents: "B", Position: domain.Position{X: 10, Y: 3}, Width: 5, Height: 3},
			},
			[]domain.Edge{
				{
					ID:   "edge-1",
					From: "a",
					To:   "b",
					Path: &domain.Path{
						Points: []domain.Position{
							{X: 8, Y: 4},
							{X: 9, Y: 4},
							{X: 10, Y: 4},
						},
					},
				},
			},
		)

		err := renderer.Render(diagram, "output.svg")
		require.NoError(t, err)

		svg := string(writer.written["output.svg"])
		assert.Contains(t, svg, `id="edge-1"`)
		assert.Contains(t, svg, `data-from="a"`)
		assert.Contains(t, svg, `data-to="b"`)
	})

	t.Run("renders with styles", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{ID: "s1", Contents: "Styled", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
			},
			nil,
		)
		diagram.Config.Styles = map[string]string{
			".highlight": "fill: yellow;",
		}

		err := renderer.Render(diagram, "styled.svg")
		require.NoError(t, err)

		svg := string(writer.written["styled.svg"])
		assert.Contains(t, svg, "fill: yellow;")
	})

	t.Run("renders with node class and style", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{
					ID: "classy", Contents: "Classy",
					Position: domain.Position{X: 3, Y: 3},
					Width: 5, Height: 3,
					Class: "my-class",
					Style: "fill: red",
				},
			},
			nil,
		)

		err := renderer.Render(diagram, "class.svg")
		require.NoError(t, err)

		svg := string(writer.written["class.svg"])
		assert.Contains(t, svg, `class="my-class"`)
		assert.Contains(t, svg, `style="fill: red"`)
	})

	t.Run("write error is propagated", func(t *testing.T) {
		writer := &mockFileWriter{
			written: map[string][]byte{},
			err:     fmt.Errorf("disk full"),
		}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{ID: "n1", Contents: "N1", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
			},
			nil,
		)

		err := renderer.Render(diagram, "fail.svg")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "disk full")
	})

	t.Run("empty nodes still renders", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(nil, nil)

		err := renderer.Render(diagram, "empty.svg")
		require.NoError(t, err)

		svg := string(writer.written["empty.svg"])
		assert.Contains(t, svg, "<svg")
	})

	t.Run("edges without paths are skipped", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{ID: "x", Contents: "X", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "y", Contents: "Y", Position: domain.Position{X: 10, Y: 3}, Width: 5, Height: 3},
			},
			[]domain.Edge{
				{ID: "no-path", From: "x", To: "y", Path: nil},
			},
		)

		err := renderer.Render(diagram, "no-paths.svg")
		require.NoError(t, err)

		svg := string(writer.written["no-paths.svg"])
		assert.False(t, strings.Contains(svg, `id="no-path"`))
	})

	t.Run("edge path class and style are rendered", func(t *testing.T) {
		writer := &mockFileWriter{written: map[string][]byte{}}
		renderer := NewSVGRenderer(writer, false)

		diagram := newTestDiagram(
			[]domain.Node{
				{ID: "p", Contents: "P", Position: domain.Position{X: 3, Y: 3}, Width: 5, Height: 3},
				{ID: "q", Contents: "Q", Position: domain.Position{X: 10, Y: 3}, Width: 5, Height: 3},
			},
			[]domain.Edge{
				{
					ID: "styled-edge", From: "p", To: "q",
					Class: "special",
					Style: "stroke: green",
					Path: &domain.Path{
						Points: []domain.Position{
							{X: 8, Y: 4},
							{X: 9, Y: 4},
							{X: 10, Y: 4},
						},
					},
				},
			},
		)

		err := renderer.Render(diagram, "styled-edge.svg")
		require.NoError(t, err)

		svg := string(writer.written["styled-edge.svg"])
		assert.Contains(t, svg, "special")
		assert.Contains(t, svg, "stroke: green")
	})
}

func TestNewSVGRenderer_implements_Renderer(t *testing.T) {
	writer := &mockFileWriter{written: map[string][]byte{}}
	var _ interface {
		Render(diagram *domain.Diagram, outputPath string) error
	} = NewSVGRenderer(writer, false)
}
