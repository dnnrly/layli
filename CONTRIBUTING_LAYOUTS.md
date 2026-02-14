# Contributing a New Layout Algorithm

This guide explains how to add a new layout algorithm to Layli.

## Overview

Layout algorithms in Layli determine how nodes are positioned on the diagram. The architecture uses a layered approach:

- **Domain Layer** (`internal/domain/`): Defines available layout types as constants
- **Layout Layer** (`layout/`): Implements the arrangement logic
- **Adapter Layer** (`internal/adapters/layout/`): Maps domain types to layout functions

## Step-by-Step Guide

### Step 1: Define the Layout Type

Add a new constant to `internal/domain/diagram.go` in the `LayoutType` constants block:

```go
const (
    LayoutFlowSquare       LayoutType = "flow-square"
    LayoutTopoSort         LayoutType = "topo-sort"
    LayoutTarjan           LayoutType = "tarjan"
    LayoutAbsolute         LayoutType = "absolute"
    LayoutRandomShortest   LayoutType = "random-shortest-square"
    LayoutMyNewLayout      LayoutType = "my-new-layout"  // Add here
)
```

The string value is what users specify in their `.layli` files:
```yaml
layout: my-new-layout
```

### Step 2: Implement the Algorithm

Create a function in `layout/arrangements.go` following the signature:

```go
func LayoutMyNewLayout(c *Config) (LayoutNodes, error) {
    // c.Nodes: slice of ConfigNode to arrange
    // c.NodeWidth, c.NodeHeight: dimensions of each node
    // c.Margin: spacing between nodes
    // c.Border: spacing from canvas edge
    
    nodes := make(LayoutNodes, len(c.Nodes))
    
    // Implement your positioning logic
    for i, n := range c.Nodes {
        nodes[i] = NewLayoutNode(
            n.Id,
            n.Contents,
            x, y,  // calculated positions
            c.NodeWidth, c.NodeHeight,
            n.Class,  // preserve CSS class
            n.Style,  // preserve inline styles
        )
    }
    
    return nodes, nil
}
```

**Important**: Always preserve node class and style attributes by passing them to `NewLayoutNode`.

### Step 3: Register in selectArrangement()

Add a case to the `selectArrangement()` function in `layout/arrangements.go`:

```go
func selectArrangement(c *Config) (LayoutArrangementFunc, error) {
    switch c.Layout {
    case "":
        return LayoutFlowSquare, nil
    // ... other cases ...
    case "my-new-layout":
        return LayoutMyNewLayout, nil
    }
    
    return nil, errors.New("do not understand layout " + c.Layout)
}
```

### Step 4: Register in the Adapter

Add a case to `selectArranger()` in `internal/adapters/layout/engine.go`:

```go
func selectArranger(lt domain.LayoutType) (layout.LayoutArrangementFunc, error) {
    switch lt {
    // ... other cases ...
    case domain.LayoutMyNewLayout:
        return layout.LayoutMyNewLayout, nil
    default:
        return nil, fmt.Errorf("unknown layout type: %s", lt)
    }
}
```

This bridges domain types with layout implementation.

### Step 5: Register for Discovery

Add the constant to `layout/options.go` so it's discoverable via the public API:

```go
// GetLayoutOptions returns all available layout algorithm names
func GetLayoutOptions() []string {
	return []string{
		string(domain.LayoutFlowSquare),
		string(domain.LayoutTopoSort),
		string(domain.LayoutTarjan),
		string(domain.LayoutAbsolute),
		string(domain.LayoutRandomShortest),
		string(domain.LayoutMyNewLayout),  // Add here
	}
}
```

This ensures the new layout appears in:
- `layli config` CLI output
- `layout.GetLayoutOptions()` programmatic API
- Any UI that uses these functions to populate dropdowns

### Step 6: Add Tests

Create tests in `layout/arrangements_test.go`:

```go
func TestLayoutMyNewLayout(t *testing.T) {
    t.Run("basic layout", func(t *testing.T) {
        c := newConfig(4, 5, 3, 2, 1)
        l, err := LayoutMyNewLayout(c)
        
        assert.NoError(t, err)
        require.Len(t, l, 4)
        
        // Test positioning invariants
        assertLeftOf(t, *l.ByID("1"), *l.ByID("2"))
        // Add more assertions...
    })
    
    t.Run("preserves class and style", func(t *testing.T) {
        c := &Config{
            Nodes: ConfigNodes{
                ConfigNode{Id: "test", Class: "my-class", Style: "color: red;"},
            },
            NodeWidth: 5, NodeHeight: 3, Margin: 1, Border: 1,
        }
        
        l, err := LayoutMyNewLayout(c)
        require.NoError(t, err)
        
        node := l.ByID("test")
        assert.Equal(t, "my-class", node.class)
        assert.Equal(t, "color: red;", node.style)
    })
}
```

Also add your layout to the `TestSelectArrangement` and `TestArrangementsPassClassAndStyle` tests.

### Step 7: Test It

```bash
# Run layout-specific tests
go test -v ./layout -run TestLayoutMyNewLayout

# Run all layout tests
go test ./layout

# Build and test with a real file
make build
./layli my-diagram.layli --output output.svg

# Verify it appears in config discovery
./layli config | grep my-new-layout
```

## Common Patterns

### Using LayoutAttempts for Configuration

The `c.LayoutAttempts` field has dual purposes depending on the layout algorithm:
- For `LayoutRandomShortestSquare`: number of random permutations to try
- For other algorithms: available for custom configuration (e.g., number of columns)

Document this clearly if your algorithm uses it differently.

### Positioning Calculation

Most algorithms follow this pattern for position calculation:
```go
xPos := c.Border + c.Margin + (col * c.NodeWidth) + (col * (c.Margin * 2))
yPos := c.Border + c.Margin + (row * c.NodeHeight) + (row * (c.Margin * 2))
```

This ensures nodes respect borders and margins consistently.

### Helper Assertions for Tests

Useful assertion functions already exist:
- `assertLeftOf(t, leftNode, rightNode)`
- `assertAbove(t, upperNode, lowerNode)`
- `assertSameRow(t, node1, node2)`
- `assertSameColumn(t, node1, node2)`

## Checklist

- [ ] Added `LayoutMyNewLayout` constant to `internal/domain/diagram.go`
- [ ] Implemented `LayoutMyNewLayout()` function in `layout/arrangements.go`
- [ ] Added case in `selectArrangement()` in `layout/arrangements.go`
- [ ] Added case in `selectArranger()` in `internal/adapters/layout/engine.go`
- [ ] **Added constant to `layout/options.go`** ← Makes it discoverable
- [ ] Added tests in `layout/arrangements_test.go`
- [ ] Updated `TestSelectArrangement` to include new layout
- [ ] Updated `TestArrangementsPassClassAndStyle` to include new layout
- [ ] Verified all tests pass: `go test ./...`
- [ ] Verified option appears in `layli config` output
- [ ] Tested with a real `.layli` file
- [ ] Preserved node class and style attributes

## Questions?

Refer to existing layouts like `LayoutFlowSquare` or `LayoutTarjan` as examples of simpler vs. more complex implementations.
