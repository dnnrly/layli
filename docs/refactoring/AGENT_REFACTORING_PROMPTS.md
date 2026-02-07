# Layli Refactoring Agent Prompts

This document contains a series of prompts designed to guide an AI coding agent (like Windsurf, Cursor, Aider, or Claude Code) through the complete refactoring of the layli project from a monolithic structure to a clean, layered architecture with BDD-first design.

## 🎯 Pre-Flight Checklist

Before starting, ensure:
- [ ] All current tests pass
- [ ] You have a clean git working directory
- [ ] You're on a new branch: `git checkout -b refactor/layered-architecture`

---

## Phase 0: Establish Baseline

### Prompt 0.1: Understand Current State

```
I'm refactoring the layli project (https://github.com/dnnrly/layli) from a monolithic structure to a clean layered architecture. This is a Go CLI tool that generates diagram layouts from configuration files.

First, help me understand the current state:

1. Run the acceptance tests and capture the baseline:
   - Execute: `make acceptance-test`
   - Save the output to `docs/baseline_acceptance.txt`
   - Count and report: how many scenarios pass/fail

2. Analyze the test structure:
   - List all files in `test/features/`
   - Show me the structure of one feature file as an example
   - List all step definition files in `test/steps/`

3. Analyze the code structure:
   - List all `.go` files in the root directory (exclude _test.go)
   - Identify which files contain:
     * Configuration parsing logic
     * Layout algorithms
     * Pathfinding logic
     * Rendering logic

4. Create a document `docs/CURRENT_STATE.md` summarizing:
   - Number of passing acceptance tests
   - List of feature files and their purpose
   - Current code organization (which files contain which logic)
   - Key dependencies between components

Do NOT make any code changes yet. Just analyze and document.
```

### Prompt 0.2: Create Feature-to-Code Map

```
Now create a detailed mapping between Gherkin features and current code:

1. For each feature file in `test/features/`:
   - Read the feature file
   - Identify which current `.go` files implement that feature
   - Document the key functions/types involved

2. Create `docs/FEATURE_MAP.md` with this structure:

```markdown
# Feature to Code Mapping

## Layout Features
**Feature File:** `test/features/layout.feature`

**Current Implementation:**
- File: `arrangements.go`
  - Functions: [list key functions]
  - Types: [list key types]
- File: `layout.go`
  - Functions: [list key functions]

**Scenarios:**
- Scenario 1: [name and brief description]
- Scenario 2: [name and brief description]

[Repeat for each feature]
```

3. Identify which step definitions in `test/steps/` map to which features

This mapping will guide our refactoring decisions.
```

### Prompt 0.3: Tag Baseline

```
Create a baseline for our refactoring:

1. Ensure all acceptance tests pass: `make acceptance-test`

2. If any tests fail, report them and STOP. We need all tests green before starting.

3. If all tests pass:
   - Create git tag: `git tag v0.0.14-pre-refactor`
   - Push tag: `git push origin v0.0.14-pre-refactor`

4. Create a commit with the documentation:
   ```bash
   git add docs/
   git commit -m "docs: establish refactoring baseline"
   ```

5. Confirm we're ready to start by reporting:
   - ✅ All acceptance tests passing
   - ✅ Baseline tagged
   - ✅ Documentation created
   - ✅ Current branch: refactor/layered-architecture
```

---

## Phase 1: Extract Domain Layer (Week 1)

### Prompt 1.1: Create Domain Package Structure

```
Create the foundation for the domain layer:

1. Create directory structure:
   ```bash
   mkdir -p internal/domain
   ```

2. Create `internal/domain/doc.go` with package documentation:
   ```go
   // Package domain contains pure business entities for layli.
   //
   // This package defines the core domain model with NO external dependencies
   // (only Go stdlib). These types represent the "ubiquitous language" from
   // our Gherkin feature files.
   //
   // Key principles:
   //   - No dependencies on external packages (except Go stdlib)
   //   - Value objects are immutable where possible
   //   - Business logic lives here, not in adapters
   //   - Types map directly to concepts in feature files
   //
   // Domain entities:
   //   - Diagram: Complete diagram specification
   //   - Node: A box/component in the diagram  
   //   - Edge: A connection between nodes
   //   - Position: X/Y coordinates on the grid
   //   - Path: A series of positions forming an edge route
   package domain
   ```

3. Run acceptance tests: `make acceptance-test`
   - They should still pass (we haven't changed behavior)

4. Commit:
   ```bash
   git add internal/domain/
   git commit -m "refactor(domain): create domain package structure"
   ```
```

### Prompt 1.2: Extract Diagram Entity

```
Extract the core Diagram entity to the domain layer:

1. Analyze the current `config.go` file and find the main configuration struct

2. Create `internal/domain/diagram.go`:
   - Extract the diagram/config struct and rename it to `Diagram`
   - Add a `DiagramConfig` struct for configuration settings
   - Create a `LayoutType` type with constants for each layout algorithm
   - Add a `Validate()` method that checks diagram invariants

3. The file should follow this pattern:
   ```go
   package domain
   
   // Diagram represents the complete diagram specification.
   // Maps to: "Given a diagram with..." in feature files.
   type Diagram struct {
       Nodes  []Node
       Edges  []Edge
       Config DiagramConfig
   }
   
   // DiagramConfig holds diagram-level settings.
   type DiagramConfig struct {
       Width   int
       Height  int
       Layout  LayoutType
       Border  int
       Margin  int
   }
   
   // LayoutType enumerates available layout algorithms.
   type LayoutType string
   
   const (
       LayoutFlowSquare LayoutType = "flow-square"
       LayoutTopoSort   LayoutType = "topo-sort"
       LayoutTarjan     LayoutType = "tarjan"
       LayoutAbsolute   LayoutType = "absolute"
   )
   
   // Validate ensures diagram invariants are met.
   func (d *Diagram) Validate() error {
       // Add validation logic
   }
   ```

4. Do NOT change any existing files yet - just create the new domain file

5. Run acceptance tests: `make acceptance-test` (should still pass)

6. Commit:
   ```bash
   git add internal/domain/diagram.go
   git commit -m "refactor(domain): extract Diagram entity"
   ```
```

### Prompt 1.3: Extract Node Entity

```
Extract the Node entity to the domain layer:

1. Analyze current code to find how nodes are represented

2. Create `internal/domain/node.go`:
   ```go
   package domain
   
   // Node represents a box/component in the diagram.
   // Maps to: "And a node 'A'" in feature files.
   type Node struct {
       ID       string
       Contents string
       Position Position
       Width    int
       Height   int
   }
   
   // Validate ensures node invariants.
   func (n *Node) Validate() error {
       if n.ID == "" {
           return fmt.Errorf("node ID cannot be empty")
       }
       if n.Width < 0 || n.Height < 0 {
           return fmt.Errorf("node dimensions must be non-negative")
       }
       return nil
   }
   
   // Bounds returns the rectangular bounds of the node.
   func (n *Node) Bounds() Bounds {
       return Bounds{
           Min: n.Position,
           Max: Position{
               X: n.Position.X + n.Width,
               Y: n.Position.Y + n.Height,
           },
       }
   }
   ```

3. Run acceptance tests: `make acceptance-test`

4. Commit:
   ```bash
   git add internal/domain/node.go
   git commit -m "refactor(domain): extract Node entity"
   ```
```

### Prompt 1.4: Extract Edge and Position Entities

```
Extract Edge and Position value objects:

1. Create `internal/domain/position.go`:
   - Move/recreate the Position struct from the root `position.go`
   - Add a Bounds struct (Min/Max positions)
   - Add helper methods (Distance, Add, etc.)

2. Create `internal/domain/edge.go`:
   ```go
   package domain
   
   // Edge represents a connection between two nodes.
   // Maps to: "And an edge from 'A' to 'B'" in feature files.
   type Edge struct {
       From string    // Node ID
       To   string    // Node ID
       Path Path      // Calculated path (may be empty before pathfinding)
   }
   
   // Validate ensures edge invariants.
   func (e *Edge) Validate() error {
       if e.From == "" || e.To == "" {
           return fmt.Errorf("edge must have from and to nodes")
       }
       if e.From == e.To {
           return fmt.Errorf("edge cannot connect node to itself")
       }
       return nil
   }
   ```

3. Create `internal/domain/path.go`:
   ```go
   package domain
   
   // Path represents a series of positions forming an edge route.
   type Path struct {
       Points []Position
   }
   
   // Length returns the total length of the path.
   func (p *Path) Length() int {
       // Calculate total distance
   }
   
   // Corners returns the number of direction changes in the path.
   func (p *Path) Corners() int {
       // Count corners
   }
   ```

4. Run acceptance tests: `make acceptance-test`

5. Commit:
   ```bash
   git add internal/domain/
   git commit -m "refactor(domain): extract Edge, Position, and Path entities"
   ```
```

### Prompt 1.5: Add Domain Unit Tests

```
Add comprehensive unit tests for the domain layer:

1. Create `internal/domain/diagram_test.go`:
   - Test Diagram.Validate() with various valid/invalid configs
   - Test edge cases (zero width, negative values, etc.)

2. Create `internal/domain/node_test.go`:
   - Test Node.Validate()
   - Test Node.Bounds()

3. Create `internal/domain/edge_test.go`:
   - Test Edge.Validate()

4. Create `internal/domain/position_test.go`:
   - Test Position operations
   - Test Bounds calculations

5. Run all tests:
   ```bash
   go test ./internal/domain/... -v
   make acceptance-test
   ```

6. Commit:
   ```bash
   git add internal/domain/
   git commit -m "test(domain): add comprehensive unit tests"
   ```

Report:
- Number of unit tests added
- Test coverage percentage for internal/domain/
- All tests passing (unit + acceptance)
```

---

## Phase 2: Use Case Layer (Week 1-2)

### Prompt 2.1: Create Use Case Package

```
Create the use case layer that orchestrates the domain:

1. Create directory and package doc:
   ```bash
   mkdir -p internal/usecases
   ```

2. Create `internal/usecases/doc.go`:
   ```go
   // Package usecases contains application use cases.
   //
   // Each use case corresponds to a complete Gherkin scenario.
   // Use cases orchestrate domain entities and depend on ports
   // (interfaces) rather than concrete implementations.
   //
   // Pattern:
   //   - Use case defines WHAT needs to happen
   //   - Ports define HOW to interact with external systems
   //   - Adapters (in internal/adapters/) implement ports
   //
   // Use cases map to Gherkin "When" steps.
   package usecases
   ```

3. Run acceptance tests: `make acceptance-test`

4. Commit:
   ```bash
   git add internal/usecases/
   git commit -m "refactor(usecases): create use case package"
   ```
```

### Prompt 2.2: Define Port Interfaces

```
Define the port interfaces that use cases will depend on:

1. Create `internal/usecases/ports.go`:
   ```go
   package usecases
   
   import "github.com/dnnrly/layli/internal/domain"
   
   // ConfigParser reads configuration files and returns domain diagrams.
   // Implementations: YAML parser, JSON parser, etc.
   type ConfigParser interface {
       // Parse reads a config file and returns a validated Diagram.
       // Maps to: "Given I have a diagram config 'file.layli'"
       Parse(path string) (*domain.Diagram, error)
   }
   
   // LayoutEngine arranges nodes within a diagram.
   // Implementations: FlowSquare, TopoSort, Tarjan, Absolute
   type LayoutEngine interface {
       // Arrange positions all nodes in the diagram.
       // Maps to: "When I arrange using 'flow-square' layout"
       Arrange(diagram *domain.Diagram) error
   }
   
   // Pathfinder calculates edge paths between nodes.
   // Implementations: Dijkstra, A*, etc.
   type Pathfinder interface {
       // FindPaths calculates paths for all edges in the diagram.
       // Maps to: "And calculate paths for all edges"
       FindPaths(diagram *domain.Diagram) error
   }
   
   // Renderer generates output from a positioned diagram.
   // Implementations: SVG, PNG, PDF
   type Renderer interface {
       // Render writes the diagram to the output path.
       // Maps to: "Then the diagram should be generated"
       Render(diagram *domain.Diagram, outputPath string) error
   }
   
   // FileReader abstracts file system reads (for testing).
   type FileReader interface {
       Read(path string) ([]byte, error)
   }
   
   // FileWriter abstracts file system writes (for testing).
   type FileWriter interface {
       Write(path string, data []byte) error
   }
   ```

2. Run acceptance tests: `make acceptance-test`

3. Commit:
   ```bash
   git add internal/usecases/ports.go
   git commit -m "refactor(usecases): define port interfaces"
   ```
```

### Prompt 2.3: Create GenerateDiagram Use Case

```
Create the main use case that orchestrates diagram generation:

1. Create `internal/usecases/generate_diagram.go`:
   ```go
   package usecases
   
   import (
       "fmt"
       "github.com/dnnrly/layli/internal/domain"
   )
   
   // GenerateDiagram orchestrates the complete diagram generation workflow.
   // Maps to a complete Gherkin scenario: Given → When → Then
   type GenerateDiagram struct {
       configParser ConfigParser
       layoutEngine LayoutEngine
       pathfinder   Pathfinder
       renderer     Renderer
   }
   
   // NewGenerateDiagram creates a new GenerateDiagram use case.
   func NewGenerateDiagram(
       parser ConfigParser,
       layout LayoutEngine,
       pathfinder Pathfinder,
       renderer Renderer,
   ) *GenerateDiagram {
       return &GenerateDiagram{
           configParser: parser,
           layoutEngine: layout,
           pathfinder:   pathfinder,
           renderer:     renderer,
       }
   }
   
   // Execute runs the complete diagram generation pipeline.
   //
   // Steps:
   //   1. Parse configuration (Given)
   //   2. Validate diagram (Given)
   //   3. Arrange layout (When)
   //   4. Calculate paths (When)
   //   5. Render output (Then)
   func (uc *GenerateDiagram) Execute(configPath, outputPath string) error {
       // Parse configuration
       diagram, err := uc.configParser.Parse(configPath)
       if err != nil {
           return fmt.Errorf("parse config: %w", err)
       }
       
       // Validate diagram
       if err := diagram.Validate(); err != nil {
           return fmt.Errorf("validate diagram: %w", err)
       }
       
       // Arrange layout
       if err := uc.layoutEngine.Arrange(diagram); err != nil {
           return fmt.Errorf("arrange layout: %w", err)
       }
       
       // Calculate paths
       if err := uc.pathfinder.FindPaths(diagram); err != nil {
           return fmt.Errorf("find paths: %w", err)
       }
       
       // Render output
       if err := uc.renderer.Render(diagram, outputPath); err != nil {
           return fmt.Errorf("render diagram: %w", err)
       }
       
       return nil
   }
   ```

2. Run acceptance tests: `make acceptance-test`

3. Commit:
   ```bash
   git add internal/usecases/generate_diagram.go
   git commit -m "refactor(usecases): create GenerateDiagram use case"
   ```
```

### Prompt 2.4: Add Use Case Tests

```
Add tests for the GenerateDiagram use case using mocks:

1. Install mockery if not already available:
   ```bash
   go install github.com/vektra/mockery/v2@latest
   ```

2. Generate mocks for the port interfaces:
   ```bash
   mockery --name=ConfigParser --dir=internal/usecases --output=internal/usecases/mocks
   mockery --name=LayoutEngine --dir=internal/usecases --output=internal/usecases/mocks
   mockery --name=Pathfinder --dir=internal/usecases --output=internal/usecases/mocks
   mockery --name=Renderer --dir=internal/usecases --output=internal/usecases/mocks
   ```

3. Create `internal/usecases/generate_diagram_test.go`:
   - Test successful execution (happy path)
   - Test parse failure
   - Test validation failure
   - Test layout failure
   - Test pathfinding failure
   - Test rendering failure

4. Example test structure:
   ```go
   func TestGenerateDiagram_Execute_Success(t *testing.T) {
       // Arrange
       mockParser := new(mocks.ConfigParser)
       mockLayout := new(mocks.LayoutEngine)
       mockPathfinder := new(mocks.Pathfinder)
       mockRenderer := new(mocks.Renderer)
       
       diagram := &domain.Diagram{
           Config: domain.DiagramConfig{Width: 10, Height: 10},
       }
       
       mockParser.On("Parse", "test.layli").Return(diagram, nil)
       mockLayout.On("Arrange", diagram).Return(nil)
       mockPathfinder.On("FindPaths", diagram).Return(nil)
       mockRenderer.On("Render", diagram, "output.svg").Return(nil)
       
       uc := NewGenerateDiagram(mockParser, mockLayout, mockPathfinder, mockRenderer)
       
       // Act
       err := uc.Execute("test.layli", "output.svg")
       
       // Assert
       assert.NoError(t, err)
       mockParser.AssertExpectations(t)
       mockLayout.AssertExpectations(t)
       mockPathfinder.AssertExpectations(t)
       mockRenderer.AssertExpectations(t)
   }
   ```

5. Run tests:
   ```bash
   go test ./internal/usecases/... -v
   make acceptance-test
   ```

6. Commit:
   ```bash
   git add internal/usecases/
   git commit -m "test(usecases): add comprehensive use case tests with mocks"
   ```
```

---

## Phase 3: Adapter Layer (Week 2-3)

### Prompt 3.1: Create Adapter Package Structure

```
Create the adapter layer structure:

1. Create directories:
   ```bash
   mkdir -p internal/adapters/config
   mkdir -p internal/adapters/layout
   mkdir -p internal/adapters/pathfinding
   mkdir -p internal/adapters/rendering
   ```

2. Create `internal/adapters/doc.go`:
   ```go
   // Package adapters contains implementations of port interfaces.
   //
   // Adapters are the concrete implementations that interact with
   // external systems or provide specific algorithms. They implement
   // the port interfaces defined in the usecases package.
   //
   // Structure:
   //   - config/     : Configuration file parsers (YAML, JSON, etc.)
   //   - layout/     : Layout algorithms (FlowSquare, TopoSort, etc.)
   //   - pathfinding/: Pathfinding algorithms (Dijkstra, A*, etc.)
   //   - rendering/  : Output renderers (SVG, PNG, etc.)
   //
   // Each adapter is independent and can be tested in isolation.
   package adapters
   ```

3. Run acceptance tests: `make acceptance-test`

4. Commit:
   ```bash
   git add internal/adapters/
   git commit -m "refactor(adapters): create adapter package structure"
   ```
```

### Prompt 3.2: Extract Config Adapter

```
Extract the configuration parsing logic to an adapter:

1. Analyze the current `config.go` file to understand YAML parsing logic

2. Create `internal/adapters/config/yaml_parser.go`:
   ```go
   package config
   
   import (
       "fmt"
       "os"
       
       "gopkg.in/yaml.v3"
       "github.com/dnnrly/layli/internal/domain"
       "github.com/dnnrly/layli/internal/usecases"
   )
   
   // YAMLParser parses YAML configuration files.
   // Implements usecases.ConfigParser interface.
   type YAMLParser struct {
       validator *Validator
   }
   
   // NewYAMLParser creates a new YAML parser.
   func NewYAMLParser() *YAMLParser {
       return &YAMLParser{
           validator: NewValidator(),
       }
   }
   
   // Parse implements usecases.ConfigParser.
   func (p *YAMLParser) Parse(path string) (*domain.Diagram, error) {
       // Read file
       data, err := os.ReadFile(path)
       if err != nil {
           return nil, fmt.Errorf("read file: %w", err)
       }
       
       // Parse YAML
       var raw rawConfig
       if err := yaml.Unmarshal(data, &raw); err != nil {
           return nil, fmt.Errorf("unmarshal yaml: %w", err)
       }
       
       // Convert to domain model
       diagram := p.toDomain(raw)
       
       // Validate
       if err := p.validator.Validate(diagram); err != nil {
           return nil, fmt.Errorf("validate: %w", err)
       }
       
       return diagram, nil
   }
   
   // rawConfig is the intermediate YAML structure
   type rawConfig struct {
       Layout string `yaml:"layout"`
       Width  int    `yaml:"width"`
       Height int    `yaml:"height"`
       Border int    `yaml:"border"`
       Margin int    `yaml:"margin"`
       Nodes  []struct {
           ID       string `yaml:"id"`
           Contents string `yaml:"contents"`
       } `yaml:"nodes"`
       Edges []struct {
           From string `yaml:"from"`
           To   string `yaml:"to"`
       } `yaml:"edges"`
   }
   
   func (p *YAMLParser) toDomain(raw rawConfig) *domain.Diagram {
       // Convert raw config to domain model
       // Move logic from current config.go here
   }
   ```

3. Create `internal/adapters/config/validator.go`:
   - Extract validation logic from current code
   - Validate diagram constraints

4. Update current `config.go` to use the new adapter (temporary bridge)

5. Run acceptance tests: `make acceptance-test` (MUST PASS)

6. Commit:
   ```bash
   git add internal/adapters/config/
   git add config.go  # if modified
   git commit -m "refactor(adapters): extract config parser adapter"
   ```
```

### Prompt 3.3: Extract Layout Adapters - Part 1 (FlowSquare)

```
Extract the first layout algorithm to establish the pattern:

1. Create `internal/adapters/layout/flow_square.go`:
   ```go
   package layout
   
   import (
       "fmt"
       "github.com/dnnrly/layli/internal/domain"
       "github.com/dnnrly/layli/internal/usecases"
   )
   
   // FlowSquare implements the flow-square layout algorithm.
   // Arranges nodes in rows and columns, like reading text.
   // Implements usecases.LayoutEngine interface.
   type FlowSquare struct {
       gridSize int
   }
   
   // NewFlowSquare creates a new FlowSquare layout engine.
   func NewFlowSquare() *FlowSquare {
       return &FlowSquare{
           gridSize: 1, // Default grid size
       }
   }
   
   // Arrange implements usecases.LayoutEngine.
   // Maps to: "When I arrange using 'flow-square' layout"
   func (l *FlowSquare) Arrange(diagram *domain.Diagram) error {
       // Move logic from current arrangements.go here
       // Algorithm:
       //   1. Calculate rows/cols based on diagram width/height
       //   2. Position nodes left-to-right, top-to-bottom
       //   3. Ensure nodes fit within bounds
       //   4. Respect border and margin settings
       
       return nil
   }
   ```

2. Find the FlowSquare logic in current codebase (likely in `arrangements.go`)

3. Move that logic to the new adapter, converting it to work with domain types

4. Update any existing code that calls FlowSquare to use the adapter

5. Run acceptance tests: `make acceptance-test` (MUST PASS)

6. Commit:
   ```bash
   git add internal/adapters/layout/
   git commit -m "refactor(adapters): extract FlowSquare layout adapter"
   ```

This establishes the pattern. Report if successful before moving to other layouts.
```

### Prompt 3.4: Extract Layout Adapters - Part 2 (Remaining Layouts)

```
Now extract the remaining layout algorithms using the same pattern:

1. Create `internal/adapters/layout/topo_sort.go`:
   - Extract TopoSort algorithm from current code
   - Implement usecases.LayoutEngine interface
   - Add appropriate documentation

2. Create `internal/adapters/layout/tarjan.go`:
   - Extract Tarjan algorithm from current code
   - Implement usecases.LayoutEngine interface
   - Add appropriate documentation

3. Create `internal/adapters/layout/absolute.go`:
   - Extract Absolute positioning from current code
   - Implement usecases.LayoutEngine interface
   - Add appropriate documentation

4. Create `internal/adapters/layout/factory.go`:
   ```go
   package layout
   
   import (
       "fmt"
       "github.com/dnnrly/layli/internal/domain"
       "github.com/dnnrly/layli/internal/usecases"
   )
   
   // Factory creates layout engines based on layout type.
   type Factory struct {
       engines map[domain.LayoutType]usecases.LayoutEngine
   }
   
   // NewFactory creates a new layout factory with all available engines.
   func NewFactory() *Factory {
       f := &Factory{
           engines: make(map[domain.LayoutType]usecases.LayoutEngine),
       }
       
       // Register all layout engines
       f.engines[domain.LayoutFlowSquare] = NewFlowSquare()
       f.engines[domain.LayoutTopoSort] = NewTopoSort()
       f.engines[domain.LayoutTarjan] = NewTarjan()
       f.engines[domain.LayoutAbsolute] = NewAbsolute()
       
       return f
   }
   
   // Get returns the layout engine for the given type.
   func (f *Factory) Get(layoutType domain.LayoutType) (usecases.LayoutEngine, error) {
       engine, ok := f.engines[layoutType]
       if !ok {
           return nil, fmt.Errorf("unknown layout type: %s", layoutType)
       }
       return engine, nil
   }
   ```

5. After each layout extraction:
   - Run: `make acceptance-test`
   - Commit if passing

6. Final commit after all layouts:
   ```bash
   git add internal/adapters/layout/
   git commit -m "refactor(adapters): extract all layout adapters and factory"
   ```

Report the status of each layout extraction.
```

### Prompt 3.5: Extract Pathfinding Adapter

```
Extract the pathfinding logic to an adapter:

1. Analyze the current `pathfinder/dijkstra/` directory

2. Create `internal/adapters/pathfinding/dijkstra.go`:
   ```go
   package pathfinding
   
   import (
       "github.com/dnnrly/layli/internal/domain"
       "github.com/dnnrly/layli/internal/usecases"
   )
   
   // Dijkstra implements pathfinding using Dijkstra's algorithm.
   // Implements usecases.Pathfinder interface.
   type Dijkstra struct {
       strategy ConflictStrategy
   }
   
   // NewDijkstra creates a new Dijkstra pathfinder.
   func NewDijkstra() *Dijkstra {
       return &Dijkstra{
           strategy: &AvoidCrossing{},
       }
   }
   
   // FindPaths implements usecases.Pathfinder.
   func (d *Dijkstra) FindPaths(diagram *domain.Diagram) error {
       // Move logic from current pathfinder/dijkstra/ here
       // Algorithm:
       //   1. Create grid from diagram
       //   2. Mark nodes as obstacles
       //   3. For each edge, find shortest path
       //   4. Apply conflict resolution strategy
       
       return nil
   }
   ```

3. Create `internal/adapters/pathfinding/strategy.go`:
   ```go
   package pathfinding
   
   import "github.com/dnnrly/layli/internal/domain"
   
   // ConflictStrategy handles path crossing/conflicts.
   type ConflictStrategy interface {
       Resolve(paths []domain.Path) error
   }
   
   // AvoidCrossing strategy ensures paths don't cross.
   type AvoidCrossing struct{}
   
   func (s *AvoidCrossing) Resolve(paths []domain.Path) error {
       // Implement path conflict resolution
       return nil
   }
   ```

4. Move the existing `pathfinder/dijkstra/` code to the new adapter location

5. Run acceptance tests: `make acceptance-test` (MUST PASS)

6. Commit:
   ```bash
   git add internal/adapters/pathfinding/
   git commit -m "refactor(adapters): extract pathfinding adapter"
   ```
```

### Prompt 3.6: Extract Rendering Adapter

```
Extract the rendering logic to an adapter:

1. Analyze current SVG generation code (likely in `layli.go`)

2. Create `internal/adapters/rendering/svg.go`:
   ```go
   package rendering
   
   import (
       "fmt"
       "os"
       "github.com/dnnrly/layli/internal/domain"
       "github.com/dnnrly/layli/internal/usecases"
   )
   
   // SVGRenderer generates SVG output from diagrams.
   // Implements usecases.Renderer interface.
   type SVGRenderer struct {
       showGrid bool
       theme    string
   }
   
   // Options configures SVG rendering.
   type Options struct {
       ShowGrid bool
       Theme    string
   }
   
   // NewSVGRenderer creates a new SVG renderer.
   func NewSVGRenderer(opts Options) *SVGRenderer {
       return &SVGRenderer{
           showGrid: opts.ShowGrid,
           theme:    opts.Theme,
       }
   }
   
   // Render implements usecases.Renderer.
   func (r *SVGRenderer) Render(diagram *domain.Diagram, outputPath string) error {
       // Move SVG generation logic here
       // Steps:
       //   1. Generate SVG header with viewBox
       //   2. Render grid (if showGrid)
       //   3. Render edges (paths)
       //   4. Render nodes (rectangles with text)
       //   5. Write to file
       
       svg := r.generateSVG(diagram)
       
       if err := os.WriteFile(outputPath, []byte(svg), 0644); err != nil {
           return fmt.Errorf("write file: %w", err)
       }
       
       return nil
   }
   
   func (r *SVGRenderer) generateSVG(diagram *domain.Diagram) string {
       // Move SVG generation logic from current code here
       return ""
   }
   ```

3. Move SVG generation code from current location to the adapter

4. Run acceptance tests: `make acceptance-test` (MUST PASS)

5. Commit:
   ```bash
   git add internal/adapters/rendering/
   git commit -m "refactor(adapters): extract SVG rendering adapter"
   ```
```

---

## Phase 4: Wire Everything Together (Week 3)

### Prompt 4.1: Update Step Definitions to Use Use Cases

```
Update the Gherkin step definitions to use the new architecture:

1. Analyze current step definitions in `test/steps/`

2. For each step definition file, update to use the new use case layer:

Example for `test/steps/diagram_steps.go`:
   ```go
   package steps
   
   import (
       "github.com/cucumber/godog"
       "github.com/dnnrly/layli/internal/usecases"
       "github.com/dnnrly/layli/internal/adapters/config"
       "github.com/dnnrly/layli/internal/adapters/layout"
       "github.com/dnnrly/layli/internal/adapters/pathfinding"
       "github.com/dnnrly/layli/internal/adapters/rendering"
   )
   
   type diagramContext struct {
       useCase    *usecases.GenerateDiagram
       configPath string
       outputPath string
       lastError  error
   }
   
   func newDiagramContext() *diagramContext {
       // Wire up dependencies
       parser := config.NewYAMLParser()
       layoutFactory := layout.NewFactory()
       pathfinder := pathfinding.NewDijkstra()
       renderer := rendering.NewSVGRenderer(rendering.Options{})
       
       // Note: We need to make the factory implement the LayoutEngine interface
       // or create a wrapper that selects the right layout based on diagram config
       
       return &diagramContext{
           useCase: usecases.NewGenerateDiagram(parser, layoutFactory, pathfinder, renderer),
       }
   }
   
   func (ctx *diagramContext) iHaveADiagramConfig(file string) error {
       ctx.configPath = file
       return nil
   }
   
   func (ctx *diagramContext) iGenerateTheDiagram() error {
       ctx.lastError = ctx.useCase.Execute(ctx.configPath, ctx.outputPath)
       return nil
   }
   
   func (ctx *diagramContext) theDiagramShouldBeGenerated() error {
       if ctx.lastError != nil {
           return fmt.Errorf("expected success but got: %w", ctx.lastError)
       }
       // Verify output file exists
       return nil
   }
   
   func InitializeDiagramScenario(ctx *godog.ScenarioContext) {
       dCtx := newDiagramContext()
       
       ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
           dCtx = newDiagramContext()
           return ctx, nil
       })
       
       ctx.Given(`^I have a diagram config "([^"]*)"$`, dCtx.iHaveADiagramConfig)
       ctx.When(`^I generate the diagram$`, dCtx.iGenerateTheDiagram)
       ctx.Then(`^the diagram should be generated$`, dCtx.theDiagramShouldBeGenerated)
   }
   ```

3. Update each step definition file similarly

4. IMPORTANT: You may need to create an adapter that wraps the layout factory
   to implement the LayoutEngine interface and select layouts dynamically

5. Run acceptance tests: `make acceptance-test` (MUST PASS)

6. Commit:
   ```bash
   git add test/steps/
   git commit -m "refactor(test): update step definitions to use new architecture"
   ```

If tests fail, debug and report the specific failures.
```

### Prompt 4.2: Update CLI to Use Architecture

```
Update the CLI entry point to use the new architecture:

1. Simplify `cmd/layli/main.go`:
   ```go
   package main
   
   import (
       "flag"
       "fmt"
       "log"
       "os"
       
       "github.com/dnnrly/layli/internal/usecases"
       "github.com/dnnrly/layli/internal/adapters/config"
       "github.com/dnnrly/layli/internal/adapters/layout"
       "github.com/dnnrly/layli/internal/adapters/pathfinding"
       "github.com/dnnrly/layli/internal/adapters/rendering"
   )
   
   func main() {
       // Parse CLI flags
       var (
           configPath = flag.String("config", "", "path to layli config file")
           outputPath = flag.String("output", "output.svg", "output file path")
           showGrid   = flag.Bool("show-grid", false, "show grid in output")
       )
       flag.Parse()
       
       if *configPath == "" {
           if len(flag.Args()) > 0 {
               *configPath = flag.Args()[0]
           } else {
               fmt.Fprintln(os.Stderr, "Error: config file required")
               flag.Usage()
               os.Exit(1)
           }
       }
       
       // Build dependencies
       parser := config.NewYAMLParser()
       layoutFactory := layout.NewFactory()
       pathfinder := pathfinding.NewDijkstra()
       renderer := rendering.NewSVGRenderer(rendering.Options{
           ShowGrid: *showGrid,
       })
       
       // TODO: Create an adapter that selects layout from factory based on diagram config
       // For now, this might need a wrapper
       
       // Create use case
       // useCase := usecases.NewGenerateDiagram(parser, layoutFactory, pathfinder, renderer)
       
       // Execute
       // if err := useCase.Execute(*configPath, *outputPath); err != nil {
       //     log.Fatalf("Error generating diagram: %v", err)
       // }
       
       // fmt.Printf("Diagram generated: %s\n", *outputPath)
   }
   ```

2. You'll likely need to create a "layout selector" adapter that:
   - Implements usecases.LayoutEngine
   - Accepts a layout.Factory
   - Reads the diagram's layout type
   - Delegates to the appropriate engine from the factory

3. Create `internal/adapters/layout/selector.go`:
   ```go
   package layout
   
   import (
       "github.com/dnnrly/layli/internal/domain"
       "github.com/dnnrly/layli/internal/usecases"
   )
   
   // Selector selects the appropriate layout engine based on diagram config.
   type Selector struct {
       factory *Factory
   }
   
   // NewSelector creates a layout selector.
   func NewSelector() *Selector {
       return &Selector{
           factory: NewFactory(),
       }
   }
   
   // Arrange implements usecases.LayoutEngine.
   func (s *Selector) Arrange(diagram *domain.Diagram) error {
       engine, err := s.factory.Get(diagram.Config.Layout)
       if err != nil {
           return err
       }
       return engine.Arrange(diagram)
   }
   ```

4. Update main.go to use the selector

5. Test the CLI manually:
   ```bash
   go run cmd/layli/main.go demo.layli --show-grid
   ```

6. Run acceptance tests: `make acceptance-test` (MUST PASS)

7. Commit:
   ```bash
   git add cmd/layli/main.go
   git add internal/adapters/layout/selector.go  # if created
   git commit -m "refactor(cli): update CLI to use new architecture"
   ```
```

### Prompt 4.3: Remove Old Code

```
Now that everything is working through the new architecture, remove old code:

1. Identify files in the root directory that have been replaced by adapters:
   - `config.go` → internal/adapters/config/
   - `arrangements.go` → internal/adapters/layout/
   - `layout.go` → internal/adapters/layout/
   - `path.go` → internal/adapters/pathfinding/
   - Parts of `layli.go` → internal/adapters/rendering/
   - `position.go` → internal/domain/position.go

2. For each file, verify:
   - Is all logic now in adapters/domain/usecases?
   - Are there any remaining references to this file?
   - Can it be safely deleted?

3. Remove old files one at a time:
   ```bash
   git rm <filename>
   make acceptance-test  # Verify still passing
   git commit -m "refactor: remove old <filename> (replaced by new architecture)"
   ```

4. Files to keep in root:
   - `*.layli` example files
   - `demo.svg`
   - `README.md`, `LICENSE`, etc.
   - `go.mod`, `go.sum`
   - `Makefile`

5. After removing each file, verify:
   - Acceptance tests pass: `make acceptance-test`
   - CLI still works: `go run cmd/layli/main.go demo.layli`

6. Report which files were removed and confirm all tests pass
```

---

## Phase 5: Test Infrastructure (Week 3-4)

### Prompt 5.1: Create Test Helpers

```
Create test helper utilities to make writing tests easier:

1. Create `test/support/diagram_builder.go`:
   ```go
   package support
   
   import "github.com/dnnrly/layli/internal/domain"
   
   // DiagramBuilder provides a fluent API for creating test diagrams.
   type DiagramBuilder struct {
       diagram *domain.Diagram
   }
   
   // NewDiagram creates a new diagram builder with defaults.
   func NewDiagram() *DiagramBuilder {
       return &DiagramBuilder{
           diagram: &domain.Diagram{
               Nodes: []domain.Node{},
               Edges: []domain.Edge{},
               Config: domain.DiagramConfig{
                   Width:  10,
                   Height: 10,
                   Layout: domain.LayoutFlowSquare,
                   Border: 2,
                   Margin: 1,
               },
           },
       }
   }
   
   // WithWidth sets the diagram width.
   func (b *DiagramBuilder) WithWidth(w int) *DiagramBuilder {
       b.diagram.Config.Width = w
       return b
   }
   
   // WithHeight sets the diagram height.
   func (b *DiagramBuilder) WithHeight(h int) *DiagramBuilder {
       b.diagram.Config.Height = h
       return b
   }
   
   // WithLayout sets the layout type.
   func (b *DiagramBuilder) WithLayout(layout domain.LayoutType) *DiagramBuilder {
       b.diagram.Config.Layout = layout
       return b
   }
   
   // AddNode adds a node to the diagram.
   func (b *DiagramBuilder) AddNode(id, contents string) *DiagramBuilder {
       b.diagram.Nodes = append(b.diagram.Nodes, domain.Node{
           ID:       id,
           Contents: contents,
           Width:    5,
           Height:   3,
       })
       return b
   }
   
   // AddEdge adds an edge to the diagram.
   func (b *DiagramBuilder) AddEdge(from, to string) *DiagramBuilder {
       b.diagram.Edges = append(b.diagram.Edges, domain.Edge{
           From: from,
           To:   to,
       })
       return b
   }
   
   // Build returns the constructed diagram.
   func (b *DiagramBuilder) Build() *domain.Diagram {
       return b.diagram
   }
   ```

2. Create example usage in step definitions to show how it simplifies tests

3. Run tests: `make acceptance-test`

4. Commit:
   ```bash
   git add test/support/
   git commit -m "test: add diagram builder test helper"
   ```
```

### Prompt 5.2: Add Integration Tests

```
Add integration tests that verify multiple layers working together:

1. Create `test/integration/` directory

2. Create `test/integration/generation_test.go`:
   ```go
   // +build integration
   
   package integration
   
   import (
       "os"
       "path/filepath"
       "testing"
       
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
       
       "github.com/dnnrly/layli/internal/usecases"
       "github.com/dnnrly/layli/internal/adapters/config"
       "github.com/dnnrly/layli/internal/adapters/layout"
       "github.com/dnnrly/layli/internal/adapters/pathfinding"
       "github.com/dnnrly/layli/internal/adapters/rendering"
   )
   
   func TestGenerateDiagram_EndToEnd(t *testing.T) {
       // Setup
       tmpDir := t.TempDir()
       outputPath := filepath.Join(tmpDir, "output.svg")
       
       // Wire up real implementations
       parser := config.NewYAMLParser()
       layoutSelector := layout.NewSelector()
       pathfinder := pathfinding.NewDijkstra()
       renderer := rendering.NewSVGRenderer(rendering.Options{})
       
       useCase := usecases.NewGenerateDiagram(parser, layoutSelector, pathfinder, renderer)
       
       // Execute
       err := useCase.Execute("../fixtures/simple.layli", outputPath)
       
       // Assert
       require.NoError(t, err)
       
       // Verify output file exists
       _, err = os.Stat(outputPath)
       assert.NoError(t, err, "output file should exist")
       
       // Verify output file has content
       content, err := os.ReadFile(outputPath)
       require.NoError(t, err)
       assert.Contains(t, string(content), "<svg", "output should be valid SVG")
   }
   
   func TestGenerateDiagram_AllLayouts(t *testing.T) {
       layouts := []struct {
           name   string
           config string
       }{
           {"FlowSquare", "../fixtures/flow-square.layli"},
           {"TopoSort", "../fixtures/topo-sort.layli"},
           {"Tarjan", "../fixtures/tarjan.layli"},
           {"Absolute", "../fixtures/absolute.layli"},
       }
       
       for _, tc := range layouts {
           t.Run(tc.name, func(t *testing.T) {
               tmpDir := t.TempDir()
               outputPath := filepath.Join(tmpDir, "output.svg")
               
               parser := config.NewYAMLParser()
               layoutSelector := layout.NewSelector()
               pathfinder := pathfinding.NewDijkstra()
               renderer := rendering.NewSVGRenderer(rendering.Options{})
               
               useCase := usecases.NewGenerateDiagram(parser, layoutSelector, pathfinder, renderer)
               
               err := useCase.Execute(tc.config, outputPath)
               
               assert.NoError(t, err, "layout %s should succeed", tc.name)
           })
       }
   }
   ```

3. Create a Makefile target for integration tests:
   ```makefile
   .PHONY: integration-test
   integration-test:
   	go test -v -tags=integration ./test/integration/...
   ```

4. Run integration tests:
   ```bash
   make integration-test
   make acceptance-test
   ```

5. Commit:
   ```bash
   git add test/integration/
   git add Makefile  # if modified
   git commit -m "test: add integration tests for diagram generation"
   ```
```

### Prompt 5.3: Improve Test Coverage

```
Analyze and improve test coverage across the codebase:

1. Generate coverage report:
   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out -o coverage.html
   ```

2. Identify packages with low coverage (<80%):
   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -func=coverage.out | grep -v "100.0%" | sort -k3 -n
   ```

3. For each package with <80% coverage:
   - Identify uncovered code paths
   - Add unit tests to cover them
   - Focus on error paths and edge cases

4. Prioritize coverage for:
   - internal/domain/ (should be near 100%)
   - internal/usecases/ (should be >90%)
   - internal/adapters/ (should be >80%)

5. Add tests for error scenarios:
   - Invalid configurations
   - File I/O errors
   - Validation failures
   - Layout constraint violations
   - Path crossing scenarios

6. Run all tests after adding coverage:
   ```bash
   go test ./... -v
   make acceptance-test
   ```

7. Commit improvements:
   ```bash
   git add internal/
   git commit -m "test: improve test coverage to >80% across all packages"
   ```

Report the before/after coverage percentages for each package.
```

---

## Phase 6: Documentation (Week 4)

### Prompt 6.1: Create Architecture Documentation

```
Document the new architecture for future developers and AI agents:

1. Create `docs/architecture/OVERVIEW.md`:
   ```markdown
   # Layli Architecture Overview
   
   ## Architecture Style
   Layli uses **Clean Architecture** with **Ports and Adapters** (Hexagonal Architecture).
   
   ## Core Principles
   
   1. **Domain First**: Business logic lives in `internal/domain/`
   2. **Use Cases**: Application logic lives in `internal/usecases/`
   3. **Adapters**: Implementation details live in `internal/adapters/`
   4. **Dependency Rule**: Dependencies point inward (domain ← usecases ← adapters)
   
   ## Layer Responsibilities
   
   ### Domain Layer (`internal/domain/`)
   - Pure business entities and value objects
   - No external dependencies (only Go stdlib)
   - Defines the "ubiquitous language" from Gherkin features
   - Examples: Diagram, Node, Edge, Position, Path
   
   ### Use Case Layer (`internal/usecases/`)
   - Application-specific business rules
   - Orchestrates domain entities
   - Defines port interfaces (ConfigParser, LayoutEngine, etc.)
   - Maps to Gherkin scenarios (Given → When → Then)
   
   ### Adapter Layer (`internal/adapters/`)
   - Concrete implementations of port interfaces
   - Interacts with external systems (files, SVG generation)
   - Replaceable without affecting domain/usecases
   - Examples: YAMLParser, FlowSquare, SVGRenderer
   
   ## Data Flow
   
   ```
   CLI → UseCase → Adapters → Domain
                  ↓
              Port Interfaces
   ```
   
   ## Adding New Features
   
   See [ADDING_FEATURES.md](ADDING_FEATURES.md) for detailed guides.
   ```

2. Create `docs/architecture/DECISIONS.md`:
   - Document key architectural decisions (ADRs)
   - Why Clean Architecture?
   - Why BDD with Gherkin?
   - Why port interfaces?

3. Create `docs/architecture/DIAGRAMS.md`:
   - Include ASCII or Mermaid diagrams showing:
     * Layer dependencies
     * Data flow
     * Key interfaces

4. Commit:
   ```bash
   git add docs/architecture/
   git commit -m "docs: add architecture documentation"
   ```
```

### Prompt 6.2: Create Developer Guide

```
Create comprehensive guides for developers and AI agents:

1. Create `docs/ADDING_FEATURES.md`:
   ```markdown
   # Adding Features to Layli
   
   ## Adding a New Layout Algorithm
   
   ### 1. Write the Feature First (BDD)
   
   Edit `test/features/layout.feature`:
   ```gherkin
   Scenario: Circular layout arranges nodes in a circle
     Given I have a diagram with 6 nodes
     When I arrange using "circular" layout
     Then the nodes should be positioned in a circle
     And the circle radius should be optimal
   ```
   
   ### 2. Run Tests (They'll be Pending)
   ```bash
   make acceptance-test
   # You'll see: "undefined step definition"
   ```
   
   ### 3. Implement Step Definitions
   
   Edit `test/steps/layout_steps.go`:
   ```go
   func (ctx *layoutContext) theNodesShouldBePositionedInACircle() error {
       // Verification logic
   }
   ```
   
   ### 4. Create the Adapter
   
   Create `internal/adapters/layout/circular.go`:
   ```go
   package layout
   
   import "github.com/dnnrly/layli/internal/domain"
   
   type Circular struct {
       // Configuration
   }
   
   func NewCircular() *Circular {
       return &Circular{}
   }
   
   func (c *Circular) Arrange(diagram *domain.Diagram) error {
       // Algorithm implementation
       return nil
   }
   ```
   
   ### 5. Register in Factory
   
   Edit `internal/adapters/layout/factory.go`:
   ```go
   f.engines[domain.LayoutCircular] = NewCircular()
   ```
   
   ### 6. Update Domain
   
   Edit `internal/domain/diagram.go`:
   ```go
   const (
       // ... existing
       LayoutCircular LayoutType = "circular"
   )
   ```
   
   ### 7. Test Until Green
   ```bash
   make acceptance-test
   # Iterate until all scenarios pass
   ```
   
   ### 8. Add Example
   
   Create `examples/circular.layli`:
   ```yaml
   layout: circular
   nodes:
     - id: a
       contents: Node A
     # ... more nodes
   ```
   
   ## Adding a New Output Format
   
   [Similar step-by-step guide for renderers]
   
   ## Adding a New Config Format
   
   [Similar step-by-step guide for parsers]
   ```

2. Create `docs/AGENT_GUIDE.md`:
   ```markdown
   # AI Agent Guide for Layli
   
   ## Quick Start
   
   ### Understanding the Codebase
   
   1. Read feature files first:
      ```bash
      cat test/features/*.feature
      ```
      These define WHAT layli does.
   
   2. Read architecture overview:
      ```bash
      cat docs/architecture/OVERVIEW.md
      ```
   
   3. Explore the code:
      - Start with `internal/domain/` (business entities)
      - Then `internal/usecases/` (application logic)
      - Then `internal/adapters/` (implementations)
   
   ### Working on Features
   
   #### Feature Development Workflow
   
   1. Write Gherkin scenario
   2. Run tests (red)
   3. Implement step definitions
   4. Implement domain/use case/adapter code
   5. Run tests (green)
   6. Refactor
   7. Commit
   
   #### Refactoring Workflow
   
   1. Ensure tests green: `make acceptance-test`
   2. Make ONE structural change
   3. Run tests: `make acceptance-test`
   4. If red: revert
   5. If green: commit
   6. Repeat
   
   ### Key Commands
   
   ```bash
   # Run all tests
   make test
   
   # Run acceptance tests
   make acceptance-test
   
   # Run integration tests
   make integration-test
   
   # Check test coverage
   make coverage-report
   
   # Run linting
   make lint
   
   # Generate mocks
   make mocks
   
   # Build CLI
   make build
   ```
   
   ### Architecture Rules
   
   1. **Domain layer has NO external dependencies**
   2. **Use cases depend ONLY on domain + port interfaces**
   3. **Adapters implement port interfaces**
   4. **Tests must pass after EVERY commit**
   5. **Feature files don't change during refactoring**
   
   ### Common Tasks
   
   See [ADDING_FEATURES.md](ADDING_FEATURES.md) for detailed guides.
   
   ### Debugging
   
   #### Acceptance Test Failures
   
   ```bash
   # Run specific feature
   go test -v ./test -godog.tags="@layout"
   
   # Show detailed output
   go test -v ./test -godog.format=pretty
   
   # Run specific scenario
   go test -v ./test -godog.tags="@flow-square"
   ```
   
   #### Finding Code
   
   ```bash
   # Find where a type is defined
   rg "type Diagram struct"
   
   # Find where a function is called
   rg "\.Arrange\("
   
   # Find Gherkin step implementation
   rg "func.*iHaveADiagram"
   ```
   ```

3. Commit:
   ```bash
   git add docs/
   git commit -m "docs: add developer and agent guides"
   ```
```

### Prompt 6.3: Update README

```
Update the main README to reflect the new architecture:

1. Edit `README.md`:
   - Add "Architecture" section after "Using layli"
   - Brief overview of the layered architecture
   - Link to detailed architecture docs
   - Update "Developing layli" section
   - Add link to AGENT_GUIDE.md

2. Add new section:
   ```markdown
   ## Architecture
   
   Layli uses Clean Architecture with Ports and Adapters (Hexagonal Architecture):
   
   - **Domain Layer** (`internal/domain/`): Pure business entities
   - **Use Case Layer** (`internal/usecases/`): Application logic
   - **Adapter Layer** (`internal/adapters/`): Implementations
   
   For detailed architecture documentation, see [docs/architecture/OVERVIEW.md](docs/architecture/OVERVIEW.md).
   
   ### For Developers
   
   - [Adding Features Guide](docs/ADDING_FEATURES.md)
   - [Developer Guide](docs/AGENT_GUIDE.md)
   
   ### Testing
   
   Layli uses Behavior-Driven Development (BDD) with Gherkin feature files:
   
   ```bash
   # Run acceptance tests (Gherkin scenarios)
   make acceptance-test
   
   # Run unit tests
   make test
   
   # Run integration tests
   make integration-test
   
   # Generate coverage report
   make coverage-report
   ```
   ```

3. Commit:
   ```bash
   git add README.md
   git commit -m "docs: update README with new architecture info"
   ```
```

---

## Phase 7: Final Validation (Week 4)

### Prompt 7.1: Comprehensive Test Run

```
Run all tests to ensure everything works:

1. Clean build:
   ```bash
   go clean -cache
   go clean -testcache
   ```

2. Run full test suite:
   ```bash
   # Unit tests
   go test ./... -v -race -coverprofile=coverage.out
   
   # Acceptance tests
   make acceptance-test
   
   # Integration tests
   make integration-test
   
   # Lint
   make lint
   ```

3. Verify CLI works:
   ```bash
   # Build
   go build -o layli cmd/layli/main.go
   
   # Test with examples
   ./layli demo.layli --output=test_output.svg
   ./layli examples/flow-square.layli --output=test_flow_square.svg --show-grid
   ```

4. Generate coverage report:
   ```bash
   go tool cover -html=coverage.out -o coverage.html
   ```

5. Create a summary report in `docs/REFACTORING_COMPLETE.md`:
   ```markdown
   # Refactoring Completion Report
   
   ## Test Results
   
   - Unit Tests: ✅ XXX/XXX passing
   - Acceptance Tests: ✅ XXX/XXX passing
   - Integration Tests: ✅ XXX/XXX passing
   - Lint: ✅ No errors
   
   ## Code Coverage
   
   - internal/domain/: XX%
   - internal/usecases/: XX%
   - internal/adapters/: XX%
   - Overall: XX%
   
   ## Files Changed
   
   - Files added: XXX
   - Files modified: XXX
   - Files deleted: XXX
   - Lines of code: XXX → XXX
   
   ## Architecture Improvements
   
   - ✅ Domain layer extracted
   - ✅ Use case layer created
   - ✅ Adapter layer implemented
   - ✅ All interfaces defined
   - ✅ Tests updated
   - ✅ Documentation complete
   
   ## Performance
   
   - Test execution time: XX seconds
   - Binary size: XX MB
   - No performance regressions detected
   ```

6. If any tests fail, report them and DO NOT PROCEED

7. If all tests pass:
   ```bash
   git add docs/REFACTORING_COMPLETE.md
   git commit -m "docs: add refactoring completion report"
   ```
```

### Prompt 7.2: Create Pull Request

```
Prepare the refactoring for review:

1. Ensure you're on the refactoring branch:
   ```bash
   git branch --show-current  # Should be: refactor/layered-architecture
   ```

2. Review all commits:
   ```bash
   git log --oneline origin/master..HEAD
   ```

3. Create a comprehensive PR description in `PR_DESCRIPTION.md`:
   ```markdown
   # Refactor: Layered Architecture with BDD-First Design
   
   ## Summary
   
   This PR refactors layli from a monolithic structure to a clean layered architecture
   with BDD-first design. All functionality is preserved - this is a pure refactoring
   with NO behavior changes.
   
   ## Motivation
   
   - Make codebase easier to maintain and extend
   - Improve testability with clear boundaries
   - Enable AI-agent-friendly development
   - Prepare for future features (new layouts, formats, etc.)
   
   ## Changes
   
   ### Architecture
   
   - ✅ **Domain Layer**: Pure business entities (`internal/domain/`)
   - ✅ **Use Case Layer**: Application logic with port interfaces (`internal/usecases/`)
   - ✅ **Adapter Layer**: Concrete implementations (`internal/adapters/`)
   
   ### Code Organization
   
   - Moved configuration parsing → `internal/adapters/config/`
   - Moved layout algorithms → `internal/adapters/layout/`
   - Moved pathfinding → `internal/adapters/pathfinding/`
   - Moved rendering → `internal/adapters/rendering/`
   - Extracted domain entities → `internal/domain/`
   - Created orchestration layer → `internal/usecases/`
   
   ### Testing
   
   - ✅ All acceptance tests passing (XXX scenarios)
   - ✅ Added comprehensive unit tests (XX% coverage)
   - ✅ Added integration tests
   - ✅ Updated step definitions to use new architecture
   - ✅ Created test helpers and builders
   
   ### Documentation
   
   - ✅ Architecture overview and diagrams
   - ✅ Developer guides (adding features, agent guide)
   - ✅ Updated README
   - ✅ ADRs (Architecture Decision Records)
   
   ## Testing Done
   
   ```bash
   make test           # All unit tests passing
   make acceptance-test # All XXX scenarios passing
   make integration-test # All integration tests passing
   make lint           # No linting errors
   ```
   
   ## Breaking Changes
   
   None. This is a pure refactoring. The CLI interface and all functionality
   remain identical.
   
   ## Migration Guide
   
   N/A - No user-facing changes.
   
   For developers, see [docs/AGENT_GUIDE.md](docs/AGENT_GUIDE.md) and
   [docs/ADDING_FEATURES.md](docs/ADDING_FEATURES.md).
   
   ## Checklist
   
   - [x] All tests passing
   - [x] No behavior changes
   - [x] Documentation updated
   - [x] Code coverage maintained/improved
   - [x] No performance regressions
   - [x] Clean commit history
   
   ## Screenshots
   
   [If applicable, add screenshots showing CLI still works identically]
   ```

4. Push branch:
   ```bash
   git push origin refactor/layered-architecture
   ```

5. Create the PR on GitHub and paste the description

6. Report the PR URL
```

### Prompt 7.3: Post-Refactoring Cleanup

```
Final cleanup and tagging:

1. After PR is merged, pull latest master:
   ```bash
   git checkout master
   git pull origin master
   ```

2. Tag the new version:
   ```bash
   git tag v0.1.0 -a -m "Major refactoring: Clean Architecture with BDD-first design"
   git push origin v0.1.0
   ```

3. Update CHANGELOG.md:
   ```markdown
   # Changelog
   
   ## [0.1.0] - 2024-XX-XX
   
   ### Changed
   - **BREAKING (Internal)**: Complete architecture refactoring
   - Implemented Clean Architecture with layered design
   - Extracted domain layer with pure business entities
   - Created use case layer with port interfaces
   - Implemented adapter layer for all implementations
   - Improved test infrastructure with builders and helpers
   - Added comprehensive documentation for developers and AI agents
   
   ### Added
   - Integration test suite
   - Architecture documentation
   - Developer guides (ADDING_FEATURES.md, AGENT_GUIDE.md)
   - Test helper utilities
   
   ### Fixed
   - N/A (pure refactoring, no bug fixes)
   
   ### Note
   No user-facing changes. CLI interface and functionality remain identical.
   This is a pure refactoring to improve code maintainability and extensibility.
   
   ## [0.0.14] - [Previous date]
   [Previous changelog entries...]
   ```

4. Create post-refactoring documentation:
   ```bash
   cat > docs/POST_REFACTORING.md << 'EOF'
   # Post-Refactoring Guide
   
   ## What Changed
   
   The codebase has been refactored from a monolithic structure to Clean Architecture.
   
   ## For Existing Contributors
   
   ### Old → New Mappings
   
   - `config.go` → `internal/adapters/config/yaml_parser.go`
   - `arrangements.go` → `internal/adapters/layout/*.go`
   - `path.go` → `internal/adapters/pathfinding/dijkstra.go`
   - `layli.go` (SVG) → `internal/adapters/rendering/svg.go`
   - `position.go` → `internal/domain/position.go`
   
   ### How to Continue Development
   
   See:
   - [docs/ADDING_FEATURES.md](ADDING_FEATURES.md) for adding new features
   - [docs/AGENT_GUIDE.md](AGENT_GUIDE.md) for working with the new structure
   - [docs/architecture/OVERVIEW.md](architecture/OVERVIEW.md) for architecture details
   
   ### Common Tasks
   
   **Adding a layout**: Create adapter in `internal/adapters/layout/`
   **Adding output format**: Create adapter in `internal/adapters/rendering/`
   **Modifying domain logic**: Edit files in `internal/domain/`
   
   ## Benefits
   
   - ✅ Easier to test (dependency injection)
   - ✅ Easier to extend (new features = new adapters)
   - ✅ Clearer boundaries (each layer has one job)
   - ✅ AI-agent friendly (well-documented interfaces)
   EOF
   ```

5. Commit and push:
   ```bash
   git add CHANGELOG.md docs/POST_REFACTORING.md
   git commit -m "docs: add post-refactoring documentation and changelog"
   git push origin master
   ```

6. Create a GitHub release for v0.1.0 with release notes

7. Celebrate! 🎉 Report completion with summary statistics
```

---

## 🚨 Emergency Procedures

### If Acceptance Tests Fail

```
STOP IMMEDIATELY. Do not proceed further.

1. Report which scenarios are failing:
   ```bash
   make acceptance-test 2>&1 | grep -A 10 "FAILED"
   ```

2. Analyze the failure:
   - Is it a step definition issue?
   - Is it a logic error in the adapter?
   - Is it a missing import/dependency?

3. Debug:
   ```bash
   # Run specific failing scenario
   go test -v ./test -godog.tags="@<tag>" -godog.format=pretty
   ```

4. Options:
   - Fix the issue and continue
   - Revert the last commit: `git revert HEAD`
   - Ask for help with specific error details
```

### If You Get Lost

```
If you're unsure about the next step or something seems wrong:

1. Check current state:
   ```bash
   git status
   git log --oneline -10
   make acceptance-test
   ```

2. Review the feature map:
   ```bash
   cat docs/FEATURE_MAP.md
   ```

3. Ask for clarification with specific details:
   - What were you trying to do?
   - What error did you encounter?
   - What does `git diff` show?
   - What do the test failures say?
```

### Rollback Procedure

```
If something goes terribly wrong:

1. Save current work:
   ```bash
   git stash
   ```

2. Return to baseline:
   ```bash
   git checkout v0.0.14-pre-refactor
   git checkout -b refactor/layered-architecture-v2
   ```

3. Review what went wrong in the stashed changes:
   ```bash
   git stash show -p
   ```

4. Start fresh with lessons learned
```

---

## 📊 Progress Tracking

After each major prompt, report:

```
✅ Phase X.Y Complete: [Description]

Status:
- Acceptance Tests: XXX/XXX passing
- Unit Test Coverage: XX%
- Files Changed: +XX -XX
- Commits: XX total

Next: Phase X.Y+1
```

---

## Notes for Agent Execution

1. **Always run acceptance tests after each change**
2. **Commit after each successful phase**
3. **Never change feature files during refactoring**
4. **Report failures immediately and STOP**
5. **Ask for clarification if uncertain**
6. **Keep commits focused (one logical change per commit)**
7. **Write descriptive commit messages**

---

## Success Criteria

The refactoring is complete when:

- ✅ All acceptance tests passing (same count as baseline)
- ✅ Unit test coverage >80% across all layers
- ✅ Integration tests added and passing
- ✅ No code in root directory (except main.go, examples, config)
- ✅ All documentation complete
- ✅ CLI works identically to before
- ✅ PR created and approved
- ✅ Version tagged as v0.1.0

---

Good luck! 🚀
