# Adding Features to Layli

This guide explains how to add new features to layli following BDD and Clean Architecture principles.

## Table of Contents

1. [Adding a New Layout Algorithm](#adding-a-new-layout-algorithm)
2. [Adding a New Output Format](#adding-a-new-output-format)
3. [Adding a New Configuration Format](#adding-a-new-configuration-format)
4. [General Workflow](#general-workflow)

## Adding a New Layout Algorithm

### Overview

Layout algorithms determine how nodes are positioned in the diagram. Each algorithm is an adapter implementing the `LayoutEngine` interface.

### Step-by-Step Guide

#### 1. Write the Gherkin Scenario

Add a new scenario to `test/features/layouts.feature`:

```gherkin
Scenario: Circular layout arranges nodes in a circle
  Given I have a diagram with 6 nodes
  When I arrange using "circular" layout
  Then the nodes should be positioned in a circle
  And the circle radius should be optimal
```

#### 2. Run Tests (They Will Fail)

```bash
make acceptance-test
# Output: undefined step: "the nodes should be positioned in a circle"
```

#### 3. Implement Step Definitions

Edit `test/steps/layout_steps.go` and add step definitions:

```go
func (ctx *layoutContext) theNodesShouldBePositionedInACircle() error {
    // Verify nodes are positioned in a circle pattern
    // Use ctx.svgDoc to access rendered SVG
    // Example:
    // - Check all nodes are equidistant from center
    // - Check angles are evenly distributed
    return nil
}

func (ctx *layoutContext) theCircleRadiusShouldBeOptimal() error {
    // Verify the circle is as small as possible
    // while not overlapping with nodes
    return nil
}
```

#### 4. Add the Layout Type to Domain

Edit `internal/domain/diagram.go`:

```go
const (
    LayoutFlowSquare       LayoutType = "flow-square"
    // ... existing layouts ...
    LayoutCircular         LayoutType = "circular"  // NEW
)
```

#### 5. Create the Adapter

Create `internal/adapters/layout/circular.go`:

```go
package layout

import (
    "fmt"
    "math"

    "github.com/dnnrly/layli/internal/domain"
)

// Circular implements LayoutEngine for circular layout.
type Circular struct {
    // Optional configuration
}

// NewCircular creates a new circular layout engine.
func NewCircular() *Circular {
    return &Circular{}
}

// Arrange positions nodes in a circle.
func (c *Circular) Arrange(diagram *domain.Diagram) error {
    if len(diagram.Nodes) == 0 {
        return fmt.Errorf("cannot arrange: no nodes in diagram")
    }

    numNodes := len(diagram.Nodes)
    centerX := float64(diagram.Config.NodeWidth + diagram.Config.Spacing)
    centerY := float64(diagram.Config.NodeHeight + diagram.Config.Spacing)

    // Calculate radius based on node size
    radius := float64(numNodes * diagram.Config.NodeWidth / 2)

    // Position each node on the circle
    for i, node := range diagram.Nodes {
        angle := 2 * math.Pi * float64(i) / float64(numNodes)
        x := centerX + radius*math.Cos(angle)
        y := centerY + radius*math.Sin(angle)

        diagram.Nodes[i].Position = domain.Position{
            X:     int(x),
            Y:     int(y),
            Width: node.Width,
            Height: node.Height,
        }
    }

    return nil
}
```

#### 6. Register the Adapter in Factory

Edit `internal/adapters/layout/engine.go` and update the factory:

```go
// NewLayoutAdapter creates a layout engine selector.
func NewLayoutAdapter() *LayoutAdapter {
    return &LayoutAdapter{
        engines: map[domain.LayoutType]LayoutEngine{
            domain.LayoutFlowSquare:     NewFlowSquare(),
            domain.LayoutTopoSort:       NewTopoSort(),
            domain.LayoutTarjan:         NewTarjan(),
            domain.LayoutAbsolute:       NewAbsolute(),
            domain.LayoutRandomShortest: NewRandomShortest(),
            domain.LayoutCircular:       NewCircular(), // ADD THIS
        },
    }
}
```

#### 7. Test Until Green

```bash
make acceptance-test
# Iterate on step definitions and implementation until all tests pass
```

#### 8. Create an Example

Create `examples/circular.layli`:

```yaml
layout: circular
nodes:
  - id: a
    contents: Node A
  - id: b
    contents: Node B
  - id: c
    contents: Node C
  - id: d
    contents: Node D
  - id: e
    contents: Node E
  - id: f
    contents: Node F
edges:
  - from: a
    to: b
  - from: b
    to: c
```

Generate it:
```bash
layli examples/circular.layli -o examples/circular.svg
```

#### 9. Add Unit Tests (Optional)

Create `internal/adapters/layout/circular_test.go`:

```go
package layout

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/dnnrly/layli/internal/domain"
    "github.com/dnnrly/layli/test/support"
)

func TestCircularLayout(t *testing.T) {
    diagram := support.NewDiagram().
        WithLayout(domain.LayoutCircular).
        AddNode("a", "Node A").
        AddNode("b", "Node B").
        AddNode("c", "Node C").
        Build()

    layout := NewCircular()
    err := layout.Arrange(diagram)

    require.NoError(t, err)
    assert.Equal(t, 3, len(diagram.Nodes))
    // Assert nodes are positioned in a circle
}
```

#### 10. Commit

```bash
git add .
git commit -m "feat(layout): add circular layout algorithm

- Create Circular layout adapter
- Implements LayoutEngine interface
- Positions nodes in a circle based on count
- Adds circular layout type to domain
- Includes acceptance test and example

Tests: 26/26 acceptance tests passing"
```

---

## Adding a New Output Format

### Overview

Output formats are adapters implementing the `Renderer` interface. Currently only SVG is supported; here's how to add JSON output.

### Step-by-Step Guide

#### 1. Add Scenario to Feature File

Edit `test/features/rendering.feature` (or create it):

```gherkin
Scenario: Outputs diagram as JSON
  Given I have a diagram with 2 nodes
  When I generate the diagram as JSON
  Then the JSON file contains node positions
  And the JSON file contains edge connections
```

#### 2. Create the Domain Type

If needed, update `internal/domain/` to support JSON serialization. The domain entities should be JSON-serializable:

```go
type Diagram struct {
    Nodes  []Node          `json:"nodes"`
    Edges  []Edge          `json:"edges"`
    Config DiagramConfig   `json:"config"`
}
```

#### 3. Create the Adapter

Create `internal/adapters/rendering/json_renderer.go`:

```go
package rendering

import (
    "encoding/json"
    "fmt"

    "github.com/dnnrly/layli/internal/domain"
    "github.com/dnnrly/layli/internal/usecases"
)

// JSONRenderer implements Renderer for JSON output.
type JSONRenderer struct {
    writer usecases.FileWriter
}

// NewJSONRenderer creates a new JSON renderer.
func NewJSONRenderer(writer usecases.FileWriter) *JSONRenderer {
    return &JSONRenderer{
        writer: writer,
    }
}

// Render generates JSON output.
func (r *JSONRenderer) Render(diagram *domain.Diagram, outputPath string) error {
    data, err := json.MarshalIndent(diagram, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling JSON: %w", err)
    }

    if err := r.writer.Write(outputPath, data); err != nil {
        return fmt.Errorf("writing JSON file: %w", err)
    }

    return nil
}
```

#### 4. Register in Composition Root

Update `internal/composition/generate_diagram.go` to support format selection:

```go
func NewGenerateDiagramWithFormat(showGrid bool, format string) *usecases.GenerateDiagram {
    reader := filesystem.NewOSFileReader()
    writer := filesystem.NewOSFileWriter()

    parser := config.NewYAMLParser(reader)
    layoutEngine := layout.NewLayoutAdapter()
    pathfinder := pathfinding.NewDijkstraPathfinder()

    var renderer usecases.Renderer
    switch format {
    case "json":
        renderer = rendering.NewJSONRenderer(writer)
    default:
        renderer = rendering.NewSVGRenderer(writer, showGrid)
    }

    return usecases.NewGenerateDiagram(parser, layoutEngine, pathfinder, renderer)
}
```

#### 5. Update CLI

Edit `cmd/layli/main.go` to accept format flag:

```go
var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a diagram",
    RunE: func(cmd *cobra.Command, args []string) error {
        format, _ := cmd.Flags().GetString("format")
        uc := composition.NewGenerateDiagramWithFormat(showGrid, format)
        return uc.Execute(inputPath, outputPath)
    },
}

func init() {
    generateCmd.Flags().StringP("format", "f", "svg", "Output format (svg, json)")
}
```

#### 6. Test and Commit

```bash
make acceptance-test
git add .
git commit -m "feat(rendering): add JSON output format

- Create JSONRenderer adapter
- Implements Renderer interface
- Serializes diagram structure to JSON
- Update CLI to support --format flag

Tests: X/26 acceptance tests passing"
```

---

## Adding a New Configuration Format

### Overview

Configuration parsers read user input and create domain `Diagram` objects. Currently YAML is supported; here's how to add JSON.

### Step-by-Step Guide

#### 1. Define the Schema

Create a JSON schema or just document the format. Example:

```json
{
  "layout": "flow-square",
  "nodes": [
    {
      "id": "a",
      "contents": "Node A"
    }
  ],
  "edges": [
    {
      "from": "a",
      "to": "b"
    }
  ]
}
```

#### 2. Create the Adapter

Create `internal/adapters/config/json_parser.go`:

```go
package config

import (
    "encoding/json"
    "fmt"

    "github.com/dnnrly/layli/internal/domain"
    "github.com/dnnrly/layli/internal/usecases"
)

// JSONParser implements ConfigParser for JSON format.
type JSONParser struct {
    reader usecases.FileReader
}

// NewJSONParser creates a new JSON config parser.
func NewJSONParser(reader usecases.FileReader) *JSONParser {
    return &JSONParser{
        reader: reader,
    }
}

// Parse reads a JSON configuration file.
func (p *JSONParser) Parse(path string) (*domain.Diagram, error) {
    data, err := p.reader.Read(path)
    if err != nil {
        return nil, fmt.Errorf("reading config file: %w", err)
    }

    var cfg struct {
        Layout string `json:"layout"`
        Nodes  []struct {
            ID       string `json:"id"`
            Contents string `json:"contents"`
        } `json:"nodes"`
        Edges []struct {
            From string `json:"from"`
            To   string `json:"to"`
        } `json:"edges"`
    }

    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parsing JSON: %w", err)
    }

    // Convert to domain types
    diagram := &domain.Diagram{
        Config: domain.DiagramConfig{
            LayoutType: domain.LayoutType(cfg.Layout),
            // ... other defaults ...
        },
    }

    for _, n := range cfg.Nodes {
        diagram.Nodes = append(diagram.Nodes, domain.Node{
            ID:       n.ID,
            Contents: n.Contents,
            Width:    5,
            Height:   3,
        })
    }

    for _, e := range cfg.Edges {
        diagram.Edges = append(diagram.Edges, domain.Edge{
            From: e.From,
            To:   e.To,
        })
    }

    if err := diagram.Validate(); err != nil {
        return nil, fmt.Errorf("validating diagram: %w", err)
    }

    return diagram, nil
}
```

#### 3. Test and Commit

Create tests in `internal/adapters/config/json_parser_test.go` and commit.

---

## General Workflow

All feature development follows this pattern:

### 1. Red Phase (Feature Fails)
```bash
# Write feature/test first
# Run tests
make acceptance-test
# Tests fail (red)
```

### 2. Green Phase (Implement)
```bash
# Implement feature
# Update tests pass
make acceptance-test
# Tests pass (green)
```

### 3. Improve Code Quality
```bash
# Improve code quality
make acceptance-test
# Ensure tests still pass
```

### 4. Commit
```bash
git add .
git commit -m "feat(scope): description"
```

## Testing Guidelines

### Unit Tests
- Test domain entities in isolation
- Test adapters with mocks
- Aim for 100% coverage of critical paths

### Integration Tests
- Test adapters working together
- Test with real dependencies
- Verify end-to-end workflows

### Acceptance Tests
- Test complete features via CLI
- Follow Given-When-Then format
- Should remain stable as implementation changes

## Key Principles

1. **Write Tests First** - Define behavior before implementation
2. **Single Responsibility** - Each component does one thing
3. **Depend on Abstractions** - Use interfaces, not concrete types
4. **Domain Models First** - Implement domain entities before adapters
5. **Keep It Simple** - Don't over-engineer for hypothetical future needs
6. **Commit Frequently** - Small, focused commits are easier to review
7. **Run Tests Often** - Catch problems early

## Related Documentation

- [Architecture Overview](./architecture/OVERVIEW.md) - How everything fits together
- [Architectural Decisions](./architecture/DECISIONS.md) - Why we made certain choices
- [Agent Guide](./AGENT_GUIDE.md) - Tips for AI agents working on the codebase
