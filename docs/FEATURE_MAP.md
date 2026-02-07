# Feature to Code Mapping

## CLI Features

**Feature File:** `test/features/cli.feature`

**Scenarios:**
1. Prints help correctly - Tests `-h` flag output
2. Prints on bad config - Tests error handling for malformed .layli files
3. Sets output file correctly - Tests `--output` flag
4. Non-existent file returns error - Tests missing input file handling
5. Errors when cannot write output - Tests write permission failures

**Current Implementation:**
- Primary files:
  - `cmd/layli/main.go` - CLI entry point using Cobra
  - `config.go` - Configuration parsing and validation
  - `layli.go` - Main orchestration and error handling
- Step definitions: `test/steps_test.go` (Background, When, Then steps for CLI)
- Key types: Command struct, Config struct from YAML

**Current Flow:**
```
CLI args → Cobra command handler → GenerateDiagram use case → adapters → output
```

**Implementation:**
- Domain: `internal/domain/config.go` (Config validation rules)
- Use Case: `internal/usecases/generate_diagram.go` (Orchestration)
- Adapter: `internal/adapters/config/yaml_parser.go` (YAML parsing)
- Adapter: `internal/adapters/cli/cobra_commands.go` (CLI handler)

---

## Layout Features

**Feature File:** `test/features/layouts.feature`

**Scenarios (14 total):**
1. Generates a single node diagram - Single node with no edges
2. Generates a 2 node diagram - Two nodes with one edge
3. Generates an image with smallest area - 14 nodes (non-square), space optimization
4. Generates an image with a square number of nodes - 9 nodes in 3x3 grid
5. Generates an image with topological sorted nodes - Respects dependency order
6. Generates an image with random shortest square nodes - FAILING (known issue)
7. Arranges paths to prevent blockages - Path routing around obstacles
8. Sets output file correctly - Custom output filename
9. Shows path grid positions - `--show-grid` debug flag
10. Corrects crossing lines without moving nodes - Path optimization
11. Generates layout with absolute positions - Node positioning mode
12. Generates layout with styles - CSS style application
13. Embeds layout details in output - SVG data attributes
14. Complex multi-node arrangements - Various test configurations

**Current Implementation:**
- Primary files:
  - `arrangements.go` (150 lines) - Grid arrangement algorithms:
    - `ArrangeFlowSquare()` - Main grid layout algorithm
    - `CalculateGridDimensions()` - Compute rows/columns
    - Margin/overlap calculation
  - `layout.go` (200 lines) - Layout type/strategy selection:
    - `LayoutStrategies` map (flow-square, topological, random-shortest)
    - `ApplyLayout()` - Dispatcher
    - `SortNodesTopologically()` - DAG ordering
  - `path.go` (140 lines) - Edge routing:
    - `FindPath()` - Route between nodes
    - `CheckAndFixCrossings()` - Detect/fix line crossings
    - `UpdateNodePositions()` - Apply final layout
- Step definitions: `test/steps_test.go`
  - SVG assertions (node count, position verification)
  - Overlap detection, path routing validation
- Key types: Node, Edge, Diagram, Position, LayoutStrategy

**Current Flow:**
```
Parsed config → ApplyLayout(strategy) 
  → ArrangeFlowSquare()/other
  → FindPath(edges)
  → CheckAndFixCrossings()
  → SVG rendering
```

**Implementation:**
- Domain: 
  - `internal/domain/diagram.go` - Diagram entity
  - `internal/domain/node.go` - Node entity with position
  - `internal/domain/edge.go` - Edge entity
- Use Case: 
  - `internal/usecases/arrange_diagram.go` - Layout orchestration
  - `internal/usecases/ports.go` - LayoutEngine interface
- Adapters:
  - `internal/adapters/layout/flow_square.go` - Grid algorithm
  - `internal/adapters/layout/topological.go` - Topological sort
  - `internal/adapters/pathfinding/dijkstra_router.go` - Pathfinding
  - `internal/adapters/rendering/svg.go` - SVG generation

---

## Error Handling Features

**Feature File:** `test/features/errors.feature`

**Scenarios (7 total):**
1. Errors when cannot find paths without crossing - Graph with no non-crossing solution
2. Errors when overlapping absolute nodes defined - Node position collision detection
3. Errors when not specifying output for to-absolute - Command validation
4. Errors when cannot find input for to-absolute - File not found
5. Errors when not specifying invalid input for to-absolute - File type validation
6. Errors when cannot parse input for to-absolute - Format error handling

**Current Implementation:**
- Primary files:
  - `arrangements.go` - Overlap detection logic
  - `path.go` - No-path-found error handling
  - `layli.go` - Error wrapping and reporting
  - `config.go` - Config validation
- Step definitions: `test/steps_test.go`
  - Error output assertions
  - Exit code checking (0 for success, non-zero for errors)
- Error patterns: "no path found", "margins overlap", "creating config:"

**Implementation:**
- Domain: `internal/domain/errors.go` - Domain-specific error types
- Use Cases: Validation at use case boundaries
- Adapters: Parse errors, file I/O errors

---

## Image Reversal Features (to-absolute)

**Feature File:** `test/features/reverse.feature`

**Scenarios (2):**
1. Flow generated image can be reversed into absolute layout file
   - Parse SVG output → Extract node positions → Regenerate .layli file
2. Can consume to-absolute converted layli file
   - Convert flow → absolute → regenerate from absolute config

**Current Implementation:**
- Primary files:
  - `layli.go` - `ToAbsolute()` command handler
  - Uses `github.com/antchfx/xmlquery` for SVG parsing
  - `config.go` - Config generation and writing
- Step definitions: `test/steps_test.go`
  - File content assertions
  - YAML parsing validation
- Key operation: SVG element extraction → Node position mapping

**Implementation:**
- Use Case: `internal/usecases/reverse_svg_to_config.go`
- Adapter: `internal/adapters/svg/parser.go` - SVG parsing
- Adapter: `internal/adapters/config/writer.go` - YAML writing

---

## Summary of Code-to-Layer Mapping

| Code Location | Current Purpose | Target Layer | Target Location |
|---|---|---|---|
| `position.go` | Position struct | Domain | `internal/domain/position.go` |
| `config.go` | YAML parsing & validation | Adapter (+ Domain validation) | Adapter: `internal/adapters/config/yaml_parser.go` + Domain: validation rules |
| `arrangements.go` | Grid layout algorithm | Adapter | `internal/adapters/layout/flow_square.go` |
| `layout.go` | Strategy selection & topological sort | Adapter | `internal/adapters/layout/strategies.go` |
| `path.go` | Pathfinding & edge routing | Adapter | `internal/adapters/pathfinding/router.go` |
| `layli.go` | Main orchestration & SVG rendering | Use Case + Adapter | Use Case: `internal/usecases/generate_diagram.go` + Adapter: `internal/adapters/rendering/svg.go` |
| `cmd/layli/main.go` | CLI handler | Adapter | `internal/adapters/cli/` |
| `vertext_map.go` | Grid vertex management | Adapter | `internal/adapters/pathfinding/vertex_map.go` |
| Test files | BDD steps | Test Layer | `test/` (steps reference domain+usecases) |

---

## Test Step Distribution

- **CLI steps:** 10 steps across 5 scenarios
- **Layout steps:** 100+ steps across 14 scenarios (heavy validation)
- **Error steps:** 20+ steps across 7 scenarios
- **Reverse steps:** 15+ steps across 2 scenarios

**Total: 141 steps in `test/steps_test.go`**

The acceptance tests are **feature-centric**, not code-centric. They test end-to-end behavior, which means they work regardless of the internal implementation.
