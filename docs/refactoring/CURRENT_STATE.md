# Current State Analysis

**Date:** February 7, 2026  
**Branch:** refactor/layered-architecture  
**Baseline Tag:** v0.0.14-pre-refactor

## Test Baseline

- **Acceptance Tests:** 26 scenarios
  - Passing: 26 ✅
  - Failing: 0
- **Test Framework:** godog/Gherkin + Go standard test package
- **Feature Files:** 4 files
  - `test/features/cli.feature` - CLI argument handling and validation
  - `test/features/errors.feature` - Error handling scenarios
  - `test/features/layouts.feature` - Layout algorithm behavior (primary feature)
  - `test/features/reverse.feature` - SVG-to-Layli conversion
- **Test Baseline Test Steps:** 146 steps (all passing)
- **Test Execution Time:** ~8.3 seconds

## Current Code Organization

### Root Directory Go Files

| File | Lines | Purpose |
|------|-------|---------|
| `layli.go` | 158 | Core domain: Diagram struct, Draw() method, AbsoluteFromSVG() |
| `layout.go` | 344 | Layout algorithms container, node arrangement, path routing |
| `arrangements.go` | 220 | Node arrangement implementations (flow-square, absolute, random) |
| `path.go` | 178 | Path finding and routing between nodes |
| `config.go` | 157 | YAML configuration parsing and validation |
| `position.go` | 57 | Simple Position struct (X, Y coordinates) |
| `vertext_map.go` | 200 | Vertex/node mapping, graph operations |
| **Total** | **1,314** | Core implementation (~3,640 LOC with tests) |

### Current Package Structure

```
layli/
├── cmd/layli/
│   └── main.go (110 lines) - CLI entry point using Cobra
├── algorithms/
│   ├── tarjan/ - Strongly connected components
│   └── topological/ - Topological sort
├── pathfinder/
│   └── dijkstra/ - Dijkstra shortest path implementation
├── mocks/ - Mock implementations for testing
├── test/
│   ├── features/ - Gherkin acceptance test scenarios
│   ├── steps_test.go - Step definitions for Gherkin scenarios
│   ├── godog_test.go - Gherkin test runner
│   └── fixtures/ - Test input files (YAML config files)
└── docs/
    └── refactoring/ - Refactoring documentation
```

### Key Code Flow

**Main Application Flow:**
1. `cmd/layli/main.go` - CLI entry point (Cobra framework)
2. `config.go` - Parse YAML configuration file
3. `layout.go` - Create layout from config
4. `arrangements.go` - Arrange nodes (flow-square, absolute, random)
5. `path.go` - Find paths between nodes using Dijkstra
6. `layli.go` - Draw diagram to SVG output

**SVG Reverse (to-absolute command):**
1. Parse SVG file
2. Extract node positions, edges, styles
3. Generate YAML config with absolute layout

### Test Organization

- **Acceptance Tests:** Defined in `test/features/*.feature` (Gherkin format)
- **Step Definitions:** `test/steps_test.go` (integrated with main test file)
- **Test Runner:** `test/godog_test.go` - godog framework initialization
- **Test Fixtures:** `test/fixtures/inputs/*.layli` - Sample configuration files
- **Unit Tests:** Mixed with main code (`*_test.go` files in root)

### Key Dependencies

- **YAML Parsing:** `gopkg.in/yaml.v3`
- **SVG Generation:** `github.com/ajstarks/svgo`
- **SVG Parsing:** `github.com/antchfx/xmlquery`
- **CLI Framework:** `github.com/spf13/cobra`
- **Graph Algorithms:** Custom implementations + standard library
- **Testing:** `github.com/cucumber/godog` (Gherkin)

## Code Complexity Analysis

| Metric | Value |
|--------|-------|
| Total Lines (*.go in root) | 3,640 |
| Total Lines (*.go all) | ~5,000+ (with tests) |
| Largest File | `layout_test.go` (438 lines) |
| Average File Size | ~200 lines |
| Root Directory .go Files | 16 files (9 production + 7 tests) |
| Packages | 7 (layli, cmd/layli, algorithms/*, pathfinder/*, test, mocks) |

## Architectural Observations

### Current Issues

1. **Monolithic Root Package**
   - All core logic (`Diagram`, `Layout`, `Config`) lives in root `layli/` package
   - No clear separation between domain logic, use cases, and adapters
   - Mix of concerns: config parsing, layout algorithms, SVG rendering all mixed

2. **No Clear Layering**
   - No `internal/domain/` for entities
   - No `internal/usecases/` for business logic
   - No `internal/adapters/` for external integrations (SVG, CLI, file I/O)
   - Configuration, layout, and pathfinding tightly coupled

3. **Testing Mixed with Code**
   - Unit tests live alongside code files (`*_test.go`)
   - Step definitions integrated into test package
   - Good test coverage but lacks separation

4. **Tight Coupling**
   - `Diagram` struct handles output directly
   - `Layout` tightly couples node arrangement with path finding
   - `Config` mixed concerns: parsing AND validation AND defaults

### Strengths

- ✅ Comprehensive acceptance tests (26 scenarios, 146 steps)
- ✅ Clear git history and feature branches
- ✅ Well-organized test fixtures
- ✅ Good use of interfaces (PathFinder, OutputFunc)
- ✅ Separation of concerns in algorithms (Tarjan, Topological, Dijkstra)

## Refactoring Target Architecture

```
layli/
├── internal/
│   ├── domain/              # Pure business logic, no dependencies
│   │   ├── diagram.go       # Diagram, Node, Edge entities
│   │   ├── config.go        # Configuration value objects
│   │   ├── layout.go        # Layout entities
│   │   └── position.go      # Position value object
│   │
│   ├── usecases/            # Application logic, orchestration
│   │   ├── generate_diagram.go    # Main use case
│   │   ├── arrange_nodes.go       # Node arrangement use case
│   │   ├── find_paths.go          # Path finding use case
│   │   └── svg_to_absolute.go     # Reverse conversion use case
│   │
│   ├── adapters/            # External system integration
│   │   ├── config/          # YAML parsing
│   │   ├── svg/             # SVG generation and parsing
│   │   ├── pathfinder/      # Dijkstra implementation
│   │   ├── layout/          # Arrangement implementations
│   │   └── cli/             # Command-line interface
│   │
│   └── repositories/        # Data access abstractions
│       └── (if needed)
│
├── cmd/layli/
│   └── main.go              # CLI entry point
│
└── test/
    ├── features/            # Gherkin scenarios (unchanged)
    └── fixtures/            # Test data (unchanged)
```

## Next Steps

**Phase 1: Extract Domain Layer**
- Create `internal/domain/` package
- Extract value objects: Position, Config
- Extract entities: Diagram, Node, Edge, Layout
- Define clear domain boundaries
- All acceptance tests should still pass

**Phase 2: Extract Use Cases**
- Create `internal/usecases/` package
- Move business logic into use cases
- Define use case interfaces
- Orchestrate domain entities in use cases

**Phase 3: Extract Adapters**
- Create `internal/adapters/` package
- Move implementation details (YAML, SVG, Dijkstra) into adapters
- Implement adapter interfaces for dependency injection

**Phase 4-7: Refine and Consolidate**
- Clean up dependencies
- Add more comprehensive tests
- Final cleanup and optimization

## Success Criteria for Phase 0

✅ All 26 acceptance tests passing  
✅ Understanding of code organization achieved  
✅ Baseline documentation created  
✅ Architecture decision recorded  
✅ Ready to begin Phase 1: Extract Domain Layer
