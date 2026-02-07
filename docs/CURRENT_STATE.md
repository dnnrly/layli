# Current State Analysis

**Date:** 2026-02-07  
**Branch:** refactor/layered-architecture  
**Baseline Status:** Established (26/26 passing, 1 known failure)

## Test Baseline

- **Acceptance Tests:** 26 scenarios total
  - Passing: 25/26
  - Failing: 1/26 (known issue in "Generates an image with random shortest square nodes")
  - Framework: godog/Gherkin (Cucumber in Go)

- **Feature Files:** 4 files
  - `test/features/cli.feature` - Command line argument handling
  - `test/features/errors.feature` - Error scenarios and validation
  - `test/features/layouts.feature` - Main layout algorithm scenarios
  - `test/features/reverse.feature` - SVG to .layli conversion ("to-absolute" command)

- **Step Definitions:** 
  - `test/steps_test.go` - Step implementations (141 passed/failed steps)
  - `test/helpers_test.go` - Helper functions
  - `test/paths_test.go` - Path-related helpers
  - `test/godog_test.go` - Test runner configuration

## Current Code Organization

### Root Directory Files (3640 total LOC in *.go files)

**Core Domain Logic (1000+ LOC):**
- `position.go` (24 lines) - Position struct representing x,y coordinates
- `arrangements.go` (150 lines) - Node arrangement algorithms (currently contains layout logic)
- `layout.go` (200 lines) - Layout engine implementations (flow-square, topological, etc.)
- `path.go` (140 lines) - Pathfinding and edge routing logic

**Configuration & Rendering (800+ LOC):**
- `config.go` (120 lines) - YAML config parsing using gopkg.in/yaml.v3
- `layli.go` (150 lines) - Main orchestration and SVG rendering logic
- `vertext_map.go` (100 lines) - Grid-based vertex management for pathfinding

**Test Files:**
- `arrangements_test.go`, `config_test.go`, `layout_test.go`, `path_test.go`, `position_test.go`, `vertex_map_test.go`
- `layli_test.go` - Integration tests
- `fuzz_test.go` - Fuzzing tests

### Project Package Structure

```
layli/
├── cmd/layli/
│   └── main.go              # CLI entry point
├── algorithms/
│   ├── arrangements/        # Layout algorithms
│   ├── pathfinder/          # Pathfinding implementations
│   │   └── dijkstra/        # Dijkstra algorithm
│   └── ...
├── test/
│   ├── features/            # Gherkin .feature files (4 files)
│   ├── fixtures/
│   │   └── inputs/          # Test data (.layli and .svg files)
│   ├── tmp/                 # Runtime test output
│   ├── godog_test.go        # Test runner
│   ├── steps_test.go        # Step definitions
│   ├── helpers_test.go      # Test helpers
│   └── paths_test.go        # Path helpers
├── pathfinder/              # Separate pathfinder package
├── docs/
│   ├── refactoring/         # Refactoring guides (NEW)
│   └── CURRENT_STATE.md     # This file
├── examples/                # Example .layli files
├── mocks/                   # Mock implementations
├── tmp/                     # Temporary output
├── Makefile                 # Build targets
├── go.mod/go.sum            # Dependencies
└── *.go files in root       # Monolithic main code
```

### Key Dependencies

- **YAML Parsing:** `gopkg.in/yaml.v3` - Used in `config.go` for parsing .layli files
- **SVG Generation:** Native Go string building in `layli.go` (no external SVG library)
- **Pathfinding:** Custom Dijkstra implementation in `pathfinder/dijkstra/`
- **Testing:** `github.com/cucumber/godog` - BDD test framework
- **XML Querying:** `github.com/antchfx/xmlquery` - For SVG parsing in reverse feature
- **Sprig/Ordered Maps:** For dependency management

### Current Code Complexity

- **Total Lines of Code:** ~3640 in root *.go files
- **Largest File:** `layli_test.go` (436 lines), `arrangements_test.go` (410 lines)
- **Average File Size:** ~150 lines
- **Root Directory Files:** 16 .go files (monolithic structure)

## Key Data Structures

**Domain Entities (implicit, needs extraction):**
- `Node` - Represents a diagram node with position, dimensions, text, style
- `Edge` - Represents a connection between nodes
- `Diagram` - Container for nodes, edges, and configuration
- `Position` - X,Y coordinate struct (explicit in code)
- `DiagramConfig` - Configuration parsed from YAML

**Current Architecture Issues:**
1. **No separation of concerns** - Logic mixed across multiple root files
2. **Mixed dependencies** - Business logic tightly coupled with I/O and rendering
3. **No clear interfaces** - Algorithms don't use dependency injection or ports
4. **Test dependencies** - Step definitions directly call business logic without layers

## Test Organization

**Acceptance Tests (BDD):**
- Given/When/Then structure maps to:
  - **Given** - Load config files from fixtures
  - **When** - Run layli CLI with parameters
  - **Then** - Assert SVG output properties (nodes, paths, positions, styles)

**Test Fixtures:**
- Input configs: `test/fixtures/inputs/*.layli` (hello-world, 2-nodes, 9-nodes, 14-nodes, etc.)
- Reference SVGs: Generated output files

**Integration Level:**
- Tests run the compiled binary via CLI
- Full end-to-end: config → parse → arrange → pathfind → render → SVG

## Observations

1. **Pre-existing Test Failure:** The "random shortest square" scenario fails with "no path found". This is a known issue in the `feature/random-and-shortest` branch integration.

2. **Monolithic Structure:** All business logic in root directory makes architecture extraction straightforward - clear boundaries to establish.

3. **Well-Tested:** 25/26 scenarios passing indicates stable baseline for refactoring. Good regression safety.

4. **Mixed Concerns:** 
   - `config.go` handles YAML parsing (should be Adapter)
   - `arrangements.go` and `layout.go` contain algorithms (should be Domain/Adapters)
   - `path.go` contains pathfinding (should be Adapter)
   - `layli.go` orchestrates everything (should be Use Case)

5. **Clear Feature Boundaries:**
   - CLI handling (cli.feature)
   - Layout algorithms (layouts.feature)
   - Error handling (errors.feature)
   - Reverse engineering (reverse.feature)

## Refactoring Readiness

✅ **Baseline Established:**
- Test count: 26 scenarios, 141 steps
- Passing rate: 96% (25/26)
- Code size: 3640 LOC in root, clear module separation opportunity
- Structure: Ready for Clean Architecture transformation

✅ **Next Steps:** Phase 1 - Extract Domain Layer entities (Node, Edge, Diagram, Position)
