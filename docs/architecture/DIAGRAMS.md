# Architecture Diagrams

Visual representations of layli's architecture and key flows.

## Dependency Layers

Clean Architecture with dependencies pointing inward:

```
┌──────────────────────────────────────────────────────┐
│                    External World                     │
│              (CLI, Files, SVG Rendering)             │
└──────────────────┬───────────────────────────────────┘
                   │
                   ↓
┌──────────────────────────────────────────────────────┐
│            Adapter Layer                             │
│  ┌──────────────────────────────────────────────┐   │
│  │ Config | Layout | Pathfinding | Rendering    │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────┬───────────────────────────────────┘
                   │
                   ↓
┌──────────────────────────────────────────────────────┐
│         Use Case Layer (GenerateDiagram)             │
│    • Orchestrates domain entities                    │
│    • Defines port interfaces                         │
│    • Implements business workflows                   │
└──────────────────┬───────────────────────────────────┘
                   │
                   ↓
┌──────────────────────────────────────────────────────┐
│          Domain Layer (Pure Entities)                │
│    • Diagram, Node, Edge, Position, Path             │
│    • Validation rules (invariants)                   │
│    • No external dependencies                        │
└──────────────────────────────────────────────────────┘

Direction of Dependencies: ↑ (inward)
External components depend on inner layers.
Inner layers know nothing about outer layers.
```

## Data Flow: Diagram Generation

```
User Input (YAML File)
        ↓
   [CLI Parse Args]
        ↓
[Composition Root]
(Wire all adapters)
        ↓
[GenerateDiagram Use Case]
        │
        ├─→ Parse Config ─→ [YAMLParser Adapter]
        │       ↓
        │   [Domain Model Created]
        │       ↓
        ├─→ Arrange Nodes ─→ [LayoutEngine Adapter]
        │       ↓
        │   [Node Positions Set]
        │       ↓
        ├─→ Find Paths ─→ [Pathfinder Adapter]
        │       ↓
        │   [Edge Routes Calculated]
        │       ↓
        ├─→ Render ─→ [SVGRenderer Adapter]
        │       ↓
        │   [SVG Generated]
        │       ↓
        └─→ Write File ─→ [FileWriter Adapter]
                ↓
          SVG Output File
```

## Package Structure

```
layli/
├── internal/
│   ├── domain/          (Pure business logic)
│   │   ├── diagram.go
│   │   ├── node.go
│   │   ├── edge.go
│   │   └── position.go
│   │
│   ├── usecases/        (Application logic & ports)
│   │   ├── generate_diagram.go  (Main use case)
│   │   ├── ports.go             (Port interfaces)
│   │   └── mocks/               (Test mocks)
│   │
│   ├── adapters/        (Implementations)
│   │   ├── config/      (YAML parser)
│   │   ├── layout/      (Layout algorithms)
│   │   ├── pathfinding/ (Dijkstra)
│   │   ├── rendering/   (SVG generation)
│   │   └── filesystem/  (File I/O)
│   │
│   └── composition/     (Dependency wiring)
│       └── generate_diagram.go
│
├── cmd/
│   └── layli/
│       └── main.go      (CLI entry point)
│
├── test/
│   ├── features/        (Gherkin specs)
│   ├── steps/           (Step implementations)
│   ├── integration/     (Integration tests)
│   ├── support/         (Test helpers)
│   └── fixtures/        (Test data)
│
└── docs/
    ├── architecture/    (This documentation)
    ├── ADDING_FEATURES.md
    └── AGENT_GUIDE.md
```

## Port and Adapter Pattern

```
Use Case Layer (Defines Ports)
    ↓
┌─────────────────────────────────────────────┐
│  type ConfigParser interface {              │
│      Parse(path string) (...Diagram, error) │
│  }                                          │
└──────────────────┬──────────────────────────┘
                   │
    ┌──────────────┴──────────────┐
    ↓                             ↓
[Adapter A]              [Adapter B]
YAML Parser              JSON Parser
(current)                (future)

This makes it easy to:
- Swap implementations
- Add new formats without changing use cases
- Test with mocks
```

## Complete System Flow

```
                         HUMAN USER
                              ↓
                    [layli command line]
                              ↓
    ┌─────────────────────────┴──────────────────────────┐
    │                                                     │
    ↓                                                     ↓
[layli generate]                              [layli to-absolute]
    │                                                     │
    ├─→ Load YAML config                               │
    │   └─→ ConfigParser.Parse()                        │
    │       └─→ Create Domain.Diagram                   │
    │           └─→ Validate                            │
    │                                                     │
    ├─→ Select & Run Layout Algorithm                  │
    │   ├─→ LayoutEngine.Arrange()                      │
    │   └─→ Update Node.Position                        │
    │                                                     │
    ├─→ Calculate Edge Routes                          │
    │   ├─→ Pathfinder.FindPaths()                      │
    │   └─→ Update Edge.Path                            │
    │                                                     │
    ├─→ Generate SVG                                   │
    │   └─→ Renderer.Render()                           │
    │       └─→ Create SVG elements                     │
    │           └─→ Write to File                       │
    │                                                     │
    └─→ Return: SVG file                               │
                │
                └─→ [SVG output file]
```

## Layout Algorithm Selection Flow

```
DiagramConfig.LayoutType
        ↓
[LayoutAdapter selects algorithm]
        │
        ├─→ LayoutFlowSquare
        │   └─→ Arranges in rectangular grid
        │
        ├─→ LayoutTopoSort
        │   └─→ Topological sort for DAGs
        │
        ├─→ LayoutTarjan
        │   └─→ Uses strongly connected components
        │
        ├─→ LayoutAbsolute
        │   └─→ Uses user-specified positions
        │
        └─→ LayoutRandomShortest
            └─→ Random with edge length optimization
        
        ↓
        [Run selected algorithm]
        ↓
        [Domain.Diagram with updated Node.Position]
```

## Test Pyramid

```
        ▲
        │
        │  ╔════════════════════════════════╗
        │  ║    Unit Tests (Domain)         ║
        │  ║    100% coverage               ║
        │  ╚════════════════════════════════╝
        │     ╔════════════════════════════════════════╗
        │     ║  Unit + Integration Tests (Adapters)   ║
        │     ║  88-100% coverage                      ║
        │     ╚════════════════════════════════════════╝
        │        ╔═══════════════════════════════════════════════╗
        │        ║  Acceptance Tests (End-to-End Gherkin)       ║
        │        ║  26 scenarios, 146 steps                     ║
        │        ╚═══════════════════════════════════════════════╝
        │
        └────────────────────────────────────────────────────────

Fast, many tests          Slow, few tests
Focused, isolated         Full system behavior
```

## Error Handling Flow

```
User Input
    ↓
┌─ Parse Error ─────────────→ [Error: Invalid config format]
├─ Validation Error ────────→ [Error: Missing required field]
├─ Layout Error ────────────→ [Error: Cannot arrange nodes]
├─ Path Error ──────────────→ [Error: No valid path exists]
├─ Rendering Error ─────────→ [Error: Cannot generate SVG]
└─ File I/O Error ──────────→ [Error: Cannot write output]
    ↓
[Error handler in CLI]
    ↓
[Format & Display Error]
    ↓
Exit with Error Code
```

## Component Interaction Sequence

```
CLI
 │
 └─→ Composition.NewGenerateDiagram(showGrid)
      │
      ├─→ filesystem.NewOSFileReader()
      ├─→ filesystem.NewOSFileWriter()
      ├─→ config.NewYAMLParser(reader)
      ├─→ layout.NewLayoutAdapter()
      ├─→ pathfinding.NewDijkstraPathfinder()
      ├─→ rendering.NewSVGRenderer(writer, showGrid)
      │
      └─→ usecases.NewGenerateDiagram(
            parser, layout, pathfinder, renderer)
          │
          └─→ uc.Execute(inputPath, outputPath)
              │
              ├─→ parser.Parse(inputPath)
              │   └─→ Domain.Diagram
              │
              ├─→ layout.Arrange(diagram)
              │   └─→ Updates Node.Position
              │
              ├─→ pathfinder.FindPaths(diagram)
              │   └─→ Updates Edge.Path
              │
              ├─→ renderer.Render(diagram, outputPath)
              │   └─→ writer.Write(svgContent)
              │
              └─→ Return error (or nil)
```

## Concurrency Model

Currently layli is single-threaded:
- One diagram processed at a time
- Sequential operations: Parse → Arrange → Path → Render

Future optimization opportunities:
- Parallel layout calculation for multiple algorithms
- Concurrent edge rendering
- Streaming SVG generation for large diagrams

(Kept simple to avoid complexity during refactoring)
