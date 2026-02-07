# Layli Architecture Overview

## Architecture Style

Layli uses **Clean Architecture** with **Ports and Adapters** (also known as Hexagonal Architecture).

This design ensures that:
- Business logic is independent of external frameworks and libraries
- Easy to test in isolation
- Easy to swap implementations (e.g., different layout algorithms)
- Future-proof against technology changes

## Core Principles

1. **Domain First**: Business logic lives in `internal/domain/`
2. **Use Cases**: Application logic lives in `internal/usecases/`
3. **Adapters**: Implementation details live in `internal/adapters/`
4. **Dependency Rule**: Dependencies point inward (Domain ← UseCases ← Adapters ← External)

## Layer Responsibilities

### Domain Layer (`internal/domain/`)

Pure business entities and value objects with no external dependencies (stdlib only).

**Key Entities:**
- `Diagram` - Complete diagram specification with nodes, edges, and config
- `DiagramConfig` - Layout settings, dimensions, algorithms
- `Node` - Individual diagram node with position and dimensions
- `Edge` - Connection between two nodes
- `Position` - Spatial coordinates and bounds for layout calculations
- `Path` - Edge routing information to prevent line crossings

**Characteristics:**
- No imports outside Go stdlib
- Defines invariants via `Validate()` methods
- Embodies the "ubiquitous language" from Gherkin features
- 100% unit test coverage

### Use Case Layer (`internal/usecases/`)

Application-specific business rules that orchestrate the domain entities.

**Key Components:**
- `GenerateDiagram` - Main use case coordinating the diagram generation workflow
- `ConfigParser` - Port for reading configuration
- `LayoutEngine` - Port for arranging nodes
- `Pathfinder` - Port for routing edges
- `Renderer` - Port for generating output

**Characteristics:**
- Defines port interfaces (contracts for adapters)
- Orchestrates domain entities through use cases
- Maps to Gherkin scenarios: Given → When → Then
- 100% unit test coverage
- Depends only on domain layer

### Adapter Layer (`internal/adapters/`)

Concrete implementations of port interfaces that interact with external systems.

**Packages:**
- `config/` - YAML configuration parsing
- `layout/` - Layout algorithms (FlowSquare, TopoSort, Tarjan, Absolute, RandomShortest)
- `pathfinding/` - Dijkstra pathfinding implementation
- `rendering/` - SVG output generation
- `filesystem/` - File I/O operations

**Characteristics:**
- Implements port interfaces from use cases
- Can be replaced without affecting domain/usecases
- Depends on external libraries as needed
- 88.9% - 100% test coverage

## Data Flow

```
CLI Entry
   ↓
Composition Root (Wire Up All Adapters)
   ↓
Use Case (GenerateDiagram.Execute)
   ├─→ ConfigParser.Parse() → Domain.Diagram
   ├─→ LayoutEngine.Arrange() → Positions
   ├─→ Pathfinder.FindPaths() → Routes
   └─→ Renderer.Render() → SVG Output
```

## Dependency Diagram

```
┌─────────────────────────────────────┐
│           External (CLI)            │
└──────────────────┬──────────────────┘
                   │
                   ↓
         ┌─────────────────────┐
         │  Composition Root   │ (Wiring)
         └──────────┬──────────┘
                    │
                    ↓
         ┌─────────────────────────────────┐
         │   Use Case Layer                │
         │ (GenerateDiagram orchestrates)  │
         └──────────┬──────────────────────┘
                    │
        ┌───────────┼───────────┬───────────┐
        ↓           ↓           ↓           ↓
    ┌────────┐ ┌────────┐ ┌─────────┐ ┌────────┐
    │ Config │ │ Layout │ │ Path    │ │Render  │
    │Adapter │ │Adapter │ │Adapter  │ │Adapter │
    └────┬───┘ └────┬───┘ └────┬────┘ └───┬────┘
         │          │          │          │
         └──────────┼──────────┼──────────┘
                    ↓
         ┌──────────────────────┐
         │    Domain Layer      │
         │  (Pure Entities)     │
         └──────────────────────┘
```

## Key Interfaces (Ports)

All adapters implement these interfaces defined in `internal/usecases/`:

### ConfigParser
```go
type ConfigParser interface {
    Parse(path string) (*domain.Diagram, error)
}
```
Implemented by: `config.YAMLParser`

### LayoutEngine
```go
type LayoutEngine interface {
    Arrange(diagram *domain.Diagram) error
}
```
Implemented by: `layout.LayoutAdapter` (which selects the appropriate algorithm)

### Pathfinder
```go
type Pathfinder interface {
    FindPaths(diagram *domain.Diagram) error
}
```
Implemented by: `pathfinding.DijkstraPathfinder`

### Renderer
```go
type Renderer interface {
    Render(diagram *domain.Diagram, outputPath string) error
}
```
Implemented by: `rendering.SVGRenderer`

## Test Strategy

### Acceptance Tests
- Verify end-to-end behavior using Gherkin scenarios
- 26 scenarios covering all major features
- Located in `test/features/`

### Integration Tests
- Verify multiple layers work together
- Test adapters in composition
- Located in `test/integration/`

### Unit Tests
- Test individual domain entities and adapters in isolation
- 97% overall code coverage
- Located in `*_test.go` files alongside implementation

## Adding New Features

See [ADDING_FEATURES.md](../ADDING_FEATURES.md) for step-by-step guides on:
- Adding new layout algorithms
- Adding new output formats
- Adding new configuration formats

## See Also

- [DECISIONS.md](./DECISIONS.md) - Architectural decision records
- [DIAGRAMS.md](./DIAGRAMS.md) - Visual diagrams of the architecture
- [../ADDING_FEATURES.md](../ADDING_FEATURES.md) - Feature development guide
- [../AGENT_GUIDE.md](../AGENT_GUIDE.md) - Guide for AI agents working on this codebase
