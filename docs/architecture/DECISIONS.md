# Architectural Decision Records (ADRs)

This document captures key decisions about layli's architecture and the rationale behind them.

## ADR-001: Use Clean Architecture with Ports and Adapters

**Date:** 2024  
**Status:** Accepted  
**Context:** Layli needed to support multiple layout algorithms, output formats, and configuration sources. The original monolithic design made it difficult to swap implementations and test in isolation.

**Decision:** Adopt Clean Architecture with explicit port interfaces (Hexagonal Architecture).

**Rationale:**
- **Testability**: Each layer can be tested independently with mocks
- **Flexibility**: Easy to add new layout algorithms without modifying core logic
- **Maintainability**: Clear separation of concerns makes code easier to understand
- **Framework Independence**: Core business logic doesn't depend on external frameworks
- **Scalability**: New adapters can be added without affecting existing code

**Consequences:**
- More files and packages than monolithic approach
- Requires discipline in enforcing dependency rules
- Additional indirection through interfaces (minor performance cost)

**Alternatives Considered:**
- Monolithic design: Simpler initially but harder to extend
- Layered architecture without ports: Less flexible for swapping implementations
- Plugin architecture: Too complex for current project scope

---

## ADR-002: Use BDD with Gherkin for Acceptance Tests

**Date:** 2024  
**Status:** Accepted  
**Context:** Layli is a visual tool where behavior is defined by how diagrams look and how algorithms work. We needed a way to specify behavior that non-programmers could understand.

**Decision:** Use Gherkin (feature files) with godog for acceptance tests as the single source of truth for behavior.

**Rationale:**
- **Shared Language**: Gherkin is readable by both developers and domain experts
- **Living Documentation**: Feature files describe what the system does in human terms
- **Executable Specs**: Tests are runnable specifications that prevent regression
- **Change History**: Feature files version control captures requirement evolution
- **Visual Regression**: Step definitions can verify SVG output properties

**Consequences:**
- Need to maintain step definitions alongside feature files
- Slower test execution than pure unit tests
- Requires discipline in keeping features and code synchronized

**Alternatives Considered:**
- Pure unit tests: Fast but doesn't capture end-to-end behavior
- Manual testing: Can't guarantee consistency or prevent regression
- API contracts: Not appropriate for a CLI tool

---

## ADR-003: Implement Layout Algorithms as Adapters

**Date:** 2024  
**Status:** Accepted  
**Context:** Different layout algorithms have different characteristics (speed, output quality, node spacing). The system should support multiple algorithms selectable by the user.

**Decision:** Implement each layout algorithm as a separate adapter implementing `LayoutEngine` interface.

**Rationale:**
- **Algorithm Independence**: Each algorithm is isolated from others
- **Easy Testing**: Each algorithm can be tested independently
- **Easy Addition**: New algorithms don't require changes to core logic
- **Runtime Selection**: User can choose algorithm via configuration
- **Performance**: Algorithms can be optimized independently

**Consequences:**
- Layout selection logic needed in `LayoutAdapter` factory
- Each algorithm must handle the same domain types
- Requires consistent error handling across implementations

**Current Implementations:**
- **FlowSquare**: Arranges nodes in a rectangular grid
- **TopoSort**: Topological sort for DAG-based layouts
- **Tarjan**: Strongly connected component decomposition
- **Absolute**: User-specified positions (no arrangement)
- **RandomShortest**: Shortest edges with randomization for tie-breaking

---

## ADR-004: Use Dijkstra for Path Finding

**Date:** 2024  
**Status:** Accepted  
**Context:** Need to find non-intersecting paths between connected nodes to avoid visual clutter and improve readability.

**Decision:** Implement Dijkstra's shortest path algorithm to find optimal routing avoiding node collisions.

**Rationale:**
- **Proven Algorithm**: Well-understood, efficient pathfinding
- **Optimal Paths**: Finds shortest routes which look cleaner visually
- **Deterministic**: Same input always produces same output
- **Extensible**: Easy to swap for different pathfinding algorithms (A*, RRT, etc.)

**Consequences:**
- Paths are recalculated for each render (not cached)
- Some diagrams may have no valid non-crossing paths (rare)
- Algorithm has O(E log V) complexity

**Alternatives Considered:**
- Straight lines: Simpler but looks cluttered with many connections
- A* pathfinding: Overkill for our constraint-free grid
- Manual path specification: Not user-friendly

---

## ADR-005: Use SVG as Primary Output Format

**Date:** 2024  
**Status:** Accepted  
**Context:** Diagrams need to be rendered visually, and the output should be:
- Web-friendly
- Scalable to any size
- Embeddable in documents
- Human-readable source

**Decision:** Use SVG (Scalable Vector Graphics) as the primary output format.

**Rationale:**
- **Vector Format**: Scales without quality loss
- **Web Native**: Works in all modern browsers
- **Human Readable**: SVG is XML-based and can be viewed in any text editor
- **Embeddable**: Can be embedded in HTML and styled with CSS
- **Tooling**: Easy to post-process with other tools

**Consequences:**
- SVG files are larger than raster formats for complex diagrams
- Rendering is done at generation time (not interactive)
- Limited animation support

**Alternatives Considered:**
- PNG/JPEG: Raster formats, not scalable
- PDF: Good for print but less web-friendly
- Canvas: Not human-readable, harder to debug

---

## ADR-006: YAML Configuration Format

**Date:** 2024  
**Status:** Accepted  
**Context:** Users need a human-friendly way to specify diagrams including nodes, edges, and layout options.

**Decision:** Use YAML as the configuration format.

**Rationale:**
- **Human Readable**: Clear structure without excessive syntax
- **Standard Format**: Familiar to DevOps/Kubernetes users
- **Type Support**: Handles strings, numbers, lists, maps natively
- **Comments**: Supports documentation inline
- **Tooling**: Good editor support and validation libraries

**Consequences:**
- YAML indentation can be error-prone
- Less formal than JSON schema validation
- Must handle parsing errors gracefully

**Example Configuration:**
```yaml
layout: flow-square
width: 100
height: 100
nodes:
  - id: a
    contents: "Node A"
edges:
  - from: a
    to: b
```

---

## ADR-007: CLI as Primary Interface

**Date:** 2024  
**Status:** Accepted  
**Context:** Layli targets automation and documentation use cases where CLI integration is valuable.

**Decision:** Implement CLI as the primary user interface using Cobra.

**Rationale:**
- **Automation**: Easy to integrate into scripts and CI/CD pipelines
- **Unix Philosophy**: Follows standard CLI conventions
- **No Dependencies**: No need for web server or GUI framework
- **Portable**: Works on Windows, Mac, Linux identically
- **Simple to Extend**: Subcommands like `to-absolute` are easy to add

**Consequences:**
- No graphical interface for interactive design
- Limited streaming/progressive rendering
- Users must learn CLI flags

**Future Alternatives:**
- Web UI: Could add later without changing core architecture
- API Server: Could expose use cases via REST/gRPC

---

## ADR-008: Go as Implementation Language

**Date:** 2024  
**Status:** Accepted  
**Context:** Need a language that's fast, reliable, and produces single-binary deployments.

**Decision:** Use Go for the entire codebase.

**Rationale:**
- **Performance**: Fast enough for diagram generation at scale
- **Single Binary**: Easy distribution, no runtime dependencies
- **Concurrency**: Good goroutine support if needed for scaling
- **Testing**: Excellent testing and benchmarking support
- **Tooling**: Strong static analysis and code quality tools

**Consequences:**
- Go is different from other languages (unfamiliar to some)
- No GUI frameworks (acceptable given CLI choice)
- Smaller ecosystem than Python/JavaScript for algorithms

---

## ADR-009: Composition Root for Dependency Injection

**Date:** 2024  
**Status:** Accepted  
**Context:** All adapters need to be instantiated and wired together before use. Needed a centralized, consistent way to do this.

**Decision:** Create a composition root (`internal/composition/`) that wires all dependencies.

**Rationale:**
- **Single Responsibility**: All wiring in one place
- **Easy Testing**: Can inject different adapters for integration tests
- **Clear Dependencies**: Easy to see what depends on what
- **Future Extensibility**: Can add dependency configuration/DI framework later

**Consequences:**
- Adds an extra layer of indirection
- Composition root must change when adding new adapters
- Not as flexible as a full DI framework (but simpler)

**Example:**
```go
func NewGenerateDiagram(showGrid bool) *usecases.GenerateDiagram {
    reader := filesystem.NewOSFileReader()
    parser := config.NewYAMLParser(reader)
    layoutEngine := layout.NewLayoutAdapter()
    // ... wire together and return
}
```

---

## Decision Summary

| Decision | Rationale | Trade-offs |
|----------|-----------|-----------|
| Clean Architecture | Testability, flexibility, maintainability | More files, needs discipline |
| BDD with Gherkin | Executable specs, shared language | Slower tests, maintenance |
| Adapter Pattern for Layouts | Easy to add new algorithms | Slight indirection |
| Dijkstra Pathfinding | Proven, optimal | Not cached, rare failures |
| SVG Output | Scalable, web-native, readable | Not interactive |
| YAML Config | Human-friendly, standard | Indentation-sensitive |
| CLI Interface | Automation-friendly, simple | No GUI, no interactivity |
| Go Language | Fast, simple deployment | Smaller ecosystem |
| Composition Root | Clear dependencies, testable | Extra indirection |

## Future Considerations

- Consider adding GraphQL API for programmatic access
- Consider web UI for interactive diagram design
- Consider caching pathfinding results for large diagrams
- Consider streaming SVG generation for huge diagrams
- Consider supporting other config formats (JSON, TOML) if needed
