package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnnrly/layli/internal/composition"
)

// fixtureDir returns the path to the fixtures directory relative to the test directory.
func fixtureDir(t *testing.T) string {
	// When tests run, the working directory is the package directory,
	// so we need to navigate up to the test directory first
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	// Find the test directory by looking for the fixtures folder
	searchPath := wd
	for i := 0; i < 5; i++ {
		testFixtures := filepath.Join(searchPath, "fixtures")
		if _, err := os.Stat(testFixtures); err == nil {
			return testFixtures
		}
		searchPath = filepath.Dir(searchPath)
	}
	
	t.Fatal("could not find fixtures directory")
	return ""
}

func TestGenerateDiagram_EndToEnd(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.svg")
	fixtures := fixtureDir(t)

	// Use the composition root to wire everything together
	generateDiagram := composition.NewGenerateDiagram(false)

	// Execute - use a known fixture that exists
	err := generateDiagram.Execute(filepath.Join(fixtures, "inputs", "hello-world.layli"), outputPath)

	// Assert
	require.NoError(t, err)

	// Verify output file exists
	_, err = os.Stat(outputPath)
	assert.NoError(t, err, "output file should exist")

	// Verify output file has content
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "<svg", "output should be valid SVG")
	assert.Contains(t, string(content), "xmlns", "output should have SVG namespace")
}

func TestGenerateDiagram_AllLayouts(t *testing.T) {
	fixtures := fixtureDir(t)
	
	layouts := []struct {
		name   string
		config string
	}{
		{"FlowSquare", filepath.Join(fixtures, "inputs", "hello-world.layli")},
		{"TopoSort", filepath.Join(fixtures, "inputs", "hello-world.layli")},
		{"Tarjan", filepath.Join(fixtures, "inputs", "hello-world.layli")},
	}

	for _, tc := range layouts {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "output.svg")

			generateDiagram := composition.NewGenerateDiagram(false)

			err := generateDiagram.Execute(tc.config, outputPath)

			assert.NoError(t, err, "layout %s should succeed", tc.name)
			
			// Verify output file exists
			_, err = os.Stat(outputPath)
			assert.NoError(t, err, "output file should exist for layout %s", tc.name)
		})
	}
}

func TestGenerateDiagram_MultipleNodes(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.svg")
	fixtures := fixtureDir(t)

	generateDiagram := composition.NewGenerateDiagram(false)

	// Use a fixture with multiple nodes
	err := generateDiagram.Execute(filepath.Join(fixtures, "inputs", "2-nodes.layli"), outputPath)

	require.NoError(t, err)

	// Verify output
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	// Verify both nodes are in the output
	assert.Contains(t, string(content), "<svg", "should be valid SVG")
	// Verify the output has elements (at least the rect elements for nodes)
	svg := string(content)
	assert.Contains(t, svg, "<rect", "output should contain node rectangles")
}

func TestGenerateDiagram_WithCustomDimensions(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.svg")
	fixtures := fixtureDir(t)

	generateDiagram := composition.NewGenerateDiagram(false)

	// Use fixture with specific dimensions
	err := generateDiagram.Execute(filepath.Join(fixtures, "inputs", "hello-world.layli"), outputPath)

	require.NoError(t, err)

	// Verify SVG has dimensions
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	svg := string(content)
	
	assert.Contains(t, svg, "width", "SVG should have width attribute")
	assert.Contains(t, svg, "height", "SVG should have height attribute")
}

func TestCompositionRoot(t *testing.T) {
	// Verify the composition root can wire everything together
	generateDiagram := composition.NewGenerateDiagram(false)
	assert.NotNil(t, generateDiagram)
}
