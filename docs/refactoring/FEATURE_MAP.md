# Feature to Code Mapping

This document maps each feature file to the current implementation code, identifying which modules must be refactored.

---

## Feature 1: Layout Behavior

**Feature File:** `test/features/layouts.feature`  
**Acceptance Scenarios:** 14 (all passing)

### Feature Coverage

1. **Generates a single node diagram** - Tests basic single-node layout
2. **Generates a 2 node diagram** - Tests 2-node layout with path
3. **Generates image with smallest area** - Tests auto-layout for 14 nodes
4. **Generates image with square number of nodes** - Tests 3x3 grid (9 nodes)
5. **Generates image with topological sorted nodes** - Tests dependency ordering
6. **Generates image with random shortest square nodes** - Tests randomized layout
7. **Arranges paths to prevent blockages** - Tests path routing
8. **Sets output file correctly** - Tests output parameter handling
9. **Shows path grid positions** - Tests debug grid visualization
10. **Corrects crossing lines without moving nodes** - Tests path uncrossing
11. **Generates layout with absolute positions** - Tests absolute layout mode
12. **Generates layout with styles** - Tests CSS styling in output
13. **Embeds layout details in output** - Tests SVG data attributes
14. **[14 scenarios total, all passing]**

### Current Implementation

**Primary Files:**
- `layout.go` (344 lines) - Main layout orchestration
  - `NewLayoutFromConfig()` - Creates layout from configuration
  - `Layout.Draw()` - Renders nodes and paths to SVG
  - `Layout.ShowGrid()` - Renders debug grid
  - `NewArrangement()` - Factory for arrangement algorithms

- `arrangements.go` (220 lines) - Node placement algorithms
  - `FlowSquareArrangement` - Default grid-based layout
  - `AbsoluteArrangement` - Fixed coordinate layout
  - `RandomArrangement` - Randomized layout
  - Overlap detection and validation

- `path.go` (178 lines) - Path finding between nodes
  - `FindPath()` - Core pathfinding logic
  - `PathFinder` interface (from `layli.go`)
  - Dijkstra integration

- `layli.go` (158 lines) - Diagram rendering
  - `Diagram.Draw()` - SVG generation using svgo library
  - SVG canvas setup and styling

**Supporting Files:**
- `config.go` - Configuration validation for layout parameters
- `position.go` - Position value object
- `vertext_map.go` - Node/vertex graph representation

**Step Definitions:**
- `test/steps_test.go` - 14 scenarios with step implementations
  - Node counting and validation
  - Path counting and non-crossing verification
  - Text boundary checking
  - SVG file assertions

### Refactoring Target

```
internal/
├── domain/
│   ├── diagram.go          # Diagram, Node, Edge entities
│   ├── layout.go           # Layout entity and types
│   ├── position.go         # Position value object
│   └── config.go           # Config value objects
│
├── usecases/
│   ├── generate_diagram.go # Orchestrate layout generation
│   └── arrange_nodes.go    # Node arrangement logic
│
└── adapters/
    ├── layout/
    │   ├── flow_square.go       # FlowSquareArrangement
    │   ├── absolute.go          # AbsoluteArrangement
    │   ├── random.go            # RandomArrangement
    │   └── arrangement_factory.go
    │
    ├── svg/
    │   ├── renderer.go          # SVG drawing logic
    │   └── grid_renderer.go     # Debug grid rendering
    │
    └── pathfinder/
        ├── dijkstra_adapter.go  # Dijkstra pathfinder
        └── path_router.go       # Path finding orchestration
```

### Key Types to Extract

From `layout.go`:
- `Layout` struct → Domain entity
- `Arrangement` interface → Domain interface
- Arrangement implementations → Adapters

From `arrangements.go`:
- `ArrangementFunc` → Use case function
- Overlap detection → Domain business rules
- Arrangement algorithms → Adapter implementations

From `path.go`:
- `PathFinder` interface → Already good, move to domain
- Path finding algorithm → Use case

### Test Strategy
- ✅ Step definitions remain unchanged
- ✅ Feature file remains unchanged
- ✅ Test fixtures remain unchanged
- ✅ All 14 scenarios must pass after refactoring

---

## Feature 2: CLI Configuration

**Feature File:** `test/features/cli.feature`  
**Acceptance Scenarios:** 5 (all passing)

### Feature Coverage

1. **Prints help correctly** - Tests `-h` flag handling
2. **Prints on bad config** - Tests error handling for invalid config
3. **Sets output file correctly** - Tests `--output` parameter
4. **Non-existent file returns error** - Tests file not found handling
5. **Errors when cannot write output** - Tests write permission errors

### Current Implementation

**Primary Files:**
- `cmd/layli/main.go` (110 lines) - CLI interface
  - Cobra command setup
  - Flag parsing (--output, --layout, --show-grid)
  - Main flow orchestration
  - Subcommand: `to-absolute` (SVG reversal)
  - Error handling and messaging

- `config.go` - Configuration file parsing and validation
  - YAML deserialization
  - Default value assignment
  - Validation rules

**Step Definitions:**
- `test/steps_test.go` - CLI-specific steps
  - App execution with parameters
  - Exit code checking
  - Output message validation
  - File existence checking

### Refactoring Target

```
internal/
├── adapters/
│   └── cli/
│       ├── main_command.go      # Main generate-diagram command
│       ├── to_absolute_command.go # Reverse command
│       └── cli_executor.go      # Command execution wrapper
│
└── usecases/
    └── cli/                      # CLI-specific use cases
```

### Notes
- CLI framework (Cobra) stays in `cmd/layli/main.go`
- Main command handler can be abstracted to use case
- Error messages propagate from domain/use cases
- Flag parsing stays at adapter layer

### Test Strategy
- ✅ All 5 scenarios remain unchanged
- ✅ Error messages from domain logic propagate correctly
- ✅ Exit codes maintained

---

## Feature 3: Error Handling

**Feature File:** `test/features/errors.feature`  
**Acceptance Scenarios:** 7 (all passing)

### Feature Coverage

1. **Cannot find paths without crossing** - Pathfinding fails gracefully
2. **Overlapping absolute nodes defined** - Node arrangement validation
3. **Not specifying output for to-absolute** - CLI validation
4. **Cannot find input for to-absolute** - File not found handling
5. **Invalid input for to-absolute** - File type validation
6. **Cannot parse input for to-absolute** - Parsing error handling
7. **[7 scenarios total, all passing]**

### Current Implementation

**Error Sources:**

From `config.go`:
- Margin validation (max 10)
- Layout attempts validation (max 10,000)
- Path attempts validation (max 10,000)
- Node validation (at least 1, all have IDs)
- Edge validation (from/to references exist)

From `layout.go`:
- Node arrangement errors (overlap detection)
- Path finding errors (no path exists)

From `path.go`:
- No path found between nodes

From `cmd/layli/main.go`:
- File opening errors
- File writing errors
- Output parameter validation

**Step Definitions:**
- `test/steps_test.go`
  - Error exit code checking
  - Error message validation
  - Partial message matching

### Refactoring Target

```
internal/
├── domain/
│   └── errors.go              # Domain error types
│
├── usecases/
│   └── errors.go              # Use case error handling
│
└── adapters/
    ├── config/
    │   └── validation_error.go
    │
    ├── layout/
    │   └── arrangement_error.go
    │
    └── pathfinder/
        └── pathfinding_error.go
```

### Error Categories
- **Configuration Errors** → domain/usecases
- **Arrangement Errors** → adapters/layout
- **Pathfinding Errors** → adapters/pathfinder
- **File I/O Errors** → adapters/cli
- **Format Errors** → adapters/config

### Test Strategy
- ✅ All error messages preserved
- ✅ Error propagation chain documented
- ✅ Error types centralized in domain

---

## Feature 4: Image Reversal (SVG to Layli)

**Feature File:** `test/features/reverse.feature`  
**Acceptance Scenarios:** 2 (all passing)

### Feature Coverage

1. **Flow generated image can be reversed into absolute** - SVG → YAML conversion
2. **Can consume to-absolute converted layli file** - Round-trip test

### Current Implementation

**Primary Files:**
- `layli.go` (158 lines)
  - `AbsoluteFromSVG()` - Core conversion function
  - SVG parsing using xmlquery
  - Node position extraction
  - Edge extraction
  - Style preservation

- `cmd/layli/main.go`
  - `to-absolute` subcommand
  - Output file handling

**Dependencies:**
- `xmlquery` library for SVG DOM parsing
- Node/edge data embedded in SVG (data-* attributes)
- Style block parsing

**Step Definitions:**
- `test/steps_test.go`
  - File existence checking
  - YAML file content validation
  - Node position verification

### Refactoring Target

```
internal/
├── adapters/
│   └── svg/
│       ├── svg_parser.go           # SVG parsing logic
│       ├── svg_to_config.go        # Conversion to Config
│       ├── config_serializer.go    # Config → YAML
│       └── metadata_extractor.go   # Data attribute extraction
│
└── usecases/
    └── convert_svg_to_layli.go    # Orchestration
```

### Key Extraction Points
- SVG parsing logic → `adapters/svg/svg_parser.go`
- Config reconstruction → `adapters/svg/svg_to_config.go`
- Style extraction → Part of parsing
- YAML serialization → `adapters/config/config_serializer.go`

### Test Strategy
- ✅ 2 scenarios unchanged
- ✅ Round-trip accuracy maintained
- ✅ Data attribute format preserved

---

## Test Infrastructure

**Test Organization:**
- Test files: `test/steps_test.go` (main), `test/helpers_test.go`, `test/paths_test.go`
- Godog runner: `test/godog_test.go` (ScenarioContext registration)
- Test fixtures: `test/fixtures/inputs/*.layli` (configuration examples)
- Test helpers: Assertion functions for SVG validation

**Key Test Helpers:**
- File existence checking
- SVG DOM querying and validation
- Node/path counting
- Overlap detection
- Path crossing detection
- Text boundary validation

**Step Implementation Pattern:**
```go
func (tc *testContext) theTestFixuresHaveBeenReset(ctx context.Context) error {
    // Reset test data
}

func (tc *testContext) theAppRunsWithParameters(ctx context.Context, params string) error {
    // Execute app with parameters
    // Capture output and exit code
}

func (tc *testContext) theAppExitsWithoutError(ctx context.Context) error {
    // Check exit code is 0
}
```

---

## Refactoring Summary

### Phase-by-Phase Mapping

**Phase 1: Extract Domain Layer**
- Extract `Position`, `Config`, `ConfigNode`, `ConfigEdge` from root
- Create domain entities: `Diagram`, `Node`, `Edge`, `Layout`
- Move interfaces to domain: `PathFinder`, `Arrangement`, `OutputFunc`

**Phase 2: Extract Use Cases**
- `GenerateDiagramUseCase` - Orchestrate main flow
- `ArrangeNodesUseCase` - Node placement logic
- `FindPathsUseCase` - Path finding orchestration
- `ConvertSVGToLayliUseCase` - Reverse conversion

**Phase 3: Extract Adapters**
- **Config Adapter:** YAML parsing, serialization, validation
- **Layout Adapter:** Arrangement algorithms (FlowSquare, Absolute, Random)
- **Pathfinder Adapter:** Dijkstra implementation
- **SVG Adapter:** SVG generation, parsing, rendering
- **CLI Adapter:** Command definitions and parameter handling

**Phase 4-7: Refinement**
- Add repository abstraction if needed
- Consolidate error handling
- Clean up dependencies
- Optimize and document

### Test Preservation Strategy

All acceptance tests will:
1. ✅ Remain **unchanged** (feature files are specification)
2. ✅ Step definitions will adapt to new code structure
3. ✅ All 26 scenarios must pass after each phase
4. ✅ Test fixtures remain unchanged

### Critical Files to Preserve

**Never modify:**
- `test/features/*.feature` - These are the specification
- Test fixtures in `test/fixtures/inputs/`

**Must adapt:**
- `test/steps_test.go` - Will import from new internal packages
- `test/godog_test.go` - May need to import from new locations

---

## Current Test Results

```
26 scenarios (26 passed)
146 steps (146 passed)
8.348770672s
```

### Test Coverage by Feature

| Feature | Scenarios | Steps | Status |
|---------|-----------|-------|--------|
| Layouts | 14 | ~70 | ✅ Passing |
| CLI | 5 | ~15 | ✅ Passing |
| Errors | 7 | ~20 | ✅ Passing |
| Reverse | 2 | ~41 | ✅ Passing |
| **Total** | **26** | **146** | ✅ All Passing |

---

## Success Criteria for Refactoring

After each phase, verify:
1. ✅ All 26 acceptance tests pass
2. ✅ All 146 steps pass
3. ✅ Test execution time remains < 10 seconds
4. ✅ Clean git commit history
5. ✅ Code organization improves without changing behavior
