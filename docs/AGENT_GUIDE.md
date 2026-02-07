# AI Agent Guide for Layli

This guide is specifically for AI coding agents (Amp, Cursor, Windsurf, Claude Code, etc.) working on the layli project.

## Quick Start

### 1. Understand the Codebase (30 minutes)

Read these in order:
1. Feature files (what the system does): `cat test/features/*.feature`
2. Architecture overview: `cat docs/architecture/OVERVIEW.md`
3. Current state: `cat docs/CURRENT_STATE.md`

### 2. Set Up Your Environment

```bash
# Clone and navigate
cd layli
git checkout -b feature/your-feature

# Verify all tests pass
make acceptance-test

# View test coverage
go test ./... -cover
```

### 3. Explore the Code Structure

```
Domain Layer (Pure Business Logic)
  └─ internal/domain/*.go
     - No external dependencies
     - Defines invariants via Validate()
     - 100% test coverage

Use Cases Layer (Application Logic)
  └─ internal/usecases/*.go
     - Port interfaces (contracts for adapters)
     - GenerateDiagram orchestrator
     - 100% test coverage

Adapters Layer (Implementation Details)
  └─ internal/adapters/
     ├─ config/      (YAML parser)
     ├─ layout/      (Layout algorithms)
     ├─ pathfinding/ (Dijkstra)
     ├─ rendering/   (SVG output)
     └─ filesystem/  (File I/O)

Tests
  └─ test/
     ├─ features/     (Gherkin scenarios)
     ├─ integration/  (Multi-layer tests)
     ├─ support/     (Test helpers)
     └─ fixtures/    (Test data)
```

## Working on Features

### BDD Workflow

1. **RED**: Write Gherkin scenario in `test/features/`
   ```gherkin
   Scenario: New feature does something
     Given some initial state
     When I perform an action
     Then I should see a result
   ```

2. **RED**: Run tests to see undefined steps
   ```bash
   make acceptance-test
   # Output: undefined step: "some initial state"
   ```

3. **IMPLEMENT**: Create step definitions in `test/steps/`
   ```go
   func (ctx *testContext) someInitialState() error {
       // Setup test data
       return nil
   }
   ```

4. **IMPLEMENT**: Create domain entities in `internal/domain/`
   ```go
   type MyEntity struct {
       // Fields
   }
   
   func (m *MyEntity) Validate() error {
       // Invariants
   }
   ```

5. **IMPLEMENT**: Create use case in `internal/usecases/`
   ```go
   type MyUseCase struct {
       adapter SomeAdapter // Interface!
   }
   
   func (u *MyUseCase) Execute(...) error {
       // Orchestrate adapters
   }
   ```

6. **IMPLEMENT**: Create adapters in `internal/adapters/`
   ```go
   type MyAdapter struct {}
   
   func (a *MyAdapter) SomeMethod() error {
       // Implementation
   }
   ```

7. **GREEN**: Run tests until passing
   ```bash
   make acceptance-test
   # Iterate until 26/26 scenarios pass
   ```

8. **REFACTOR**: Improve code without changing behavior
   ```bash
   make acceptance-test
   # Ensure tests still pass
   ```

9. **COMMIT**: Create focused commits
   ```bash
   git add .
   git commit -m "feat(scope): description
   
   Detailed explanation of changes.
   
   Tests: 26/26 acceptance passing"
   ```

### Refactoring Workflow

When refactoring without adding features:

1. Ensure tests are green: `make acceptance-test`
2. Make ONE small change
3. Run tests: `make acceptance-test`
4. If red: revert with `git checkout .`
5. If green: commit immediately
6. Repeat until done

**Key Rule**: Never have a failing test. If you break something, fix it or revert.

## Common Tasks

### Add a New Layout Algorithm

See: [Adding a New Layout Algorithm](./ADDING_FEATURES.md#adding-a-new-layout-algorithm)

Quick checklist:
- [ ] Add scenario to `test/features/layouts.feature`
- [ ] Add step definitions to `test/steps/layout_steps.go`
- [ ] Add `LayoutType` constant to `internal/domain/diagram.go`
- [ ] Create adapter in `internal/adapters/layout/new_algorithm.go`
- [ ] Update factory in `internal/adapters/layout/engine.go`
- [ ] All tests passing: `make acceptance-test`

### Add a New Renderer

See: [Adding a New Output Format](./ADDING_FEATURES.md#adding-a-new-output-format)

### Add a New Config Format

See: [Adding a New Configuration Format](./ADDING_FEATURES.md#adding-a-new-configuration-format)

### Fix a Bug

1. Find the failing test
2. Understand what's expected
3. Find the code causing the bug
4. Fix the code
5. Verify test passes
6. Commit with message: `fix(scope): description`

### Improve Test Coverage

```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Find low-coverage areas
go tool cover -func=coverage.out | grep -v "100.0%"

# Add tests to cover the gaps
# Run tests to verify
make acceptance-test
```

## Design Patterns Used

### Ports and Adapters

All external interactions go through interfaces:

```go
// Port (defined in usecases)
type ConfigParser interface {
    Parse(path string) (*domain.Diagram, error)
}

// Adapter (defined in adapters)
type YAMLParser struct {
    reader FileReader
}

func (p *YAMLParser) Parse(path string) (*domain.Diagram, error) {
    // Implementation
}

// Usage (in use case)
func (u *GenerateDiagram) Execute(...) error {
    diagram, err := u.parser.Parse(inputPath)
    // ...
}
```

This makes it easy to:
- Swap implementations (YAMLParser → JSONParser)
- Test with mocks (MockConfigParser)
- Extend without modifying existing code

### Fluent Builder

Test helpers use the builder pattern:

```go
diagram := support.NewDiagram().
    WithLayout(domain.LayoutFlowSquare).
    AddNode("a", "Node A").
    AddNode("b", "Node B").
    AddEdge("a", "b").
    Build()
```

### Dependency Injection

The composition root wires everything:

```go
func NewGenerateDiagram(showGrid bool) *usecases.GenerateDiagram {
    reader := filesystem.NewOSFileReader()
    writer := filesystem.NewOSFileWriter()
    parser := config.NewYAMLParser(reader)
    layout := layout.NewLayoutAdapter()
    pathfinder := pathfinding.NewDijkstraPathfinder()
    renderer := rendering.NewSVGRenderer(writer, showGrid)
    return usecases.NewGenerateDiagram(parser, layout, pathfinder, renderer)
}
```

## Key Commands

```bash
# Run acceptance tests (must always pass!)
make acceptance-test

# Run all tests with coverage
make test

# Run integration tests only
go test ./test/integration -v

# Generate coverage HTML report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linter
go vet ./...

# Build binary
go build -o layli ./cmd/layli

# Test a specific feature
go test -v ./test -godog.tags="@layout"

# See detailed test output
go test -v ./test -godog.format=pretty

# Check what depends on a package
go list -f '{{.Imports}}' ./internal/domain

# Find all TODOs
grep -r "TODO" internal/
```

## Code Style Guidelines

### Naming

- **Packages**: lowercase, short, no underscores: `config`, `layout`, `pathfinding`
- **Types**: PascalCase: `DiagramConfig`, `LayoutEngine`, `SVGRenderer`
- **Functions**: camelCase: `NewDiagram()`, `Validate()`, `Arrange()`
- **Constants**: UPPER_SNAKE_CASE (for enums): `LayoutFlowSquare`, `LayoutTopoSort`
- **Interfaces**: Usually end with a noun/verb: `Reader`, `Writer`, `Parser`, `Engine`

### Structure

1. Package declaration
2. Imports (stdlib, third-party, local)
3. Type definitions
4. Constructor functions
5. Interface implementations
6. Helper functions

Example:
```go
package layout

import (
    "fmt"
    "github.com/dnnrly/layli/internal/domain"
)

type FlowSquare struct {
    // fields
}

func NewFlowSquare() *FlowSquare {
    return &FlowSquare{}
}

func (fs *FlowSquare) Arrange(diagram *domain.Diagram) error {
    // implementation
}

// helper functions
func calculateGrid(count int) (rows, cols int) {
    // ...
}
```

### Error Handling

Always wrap errors with context:

```go
// Good
if err := parser.Parse(path); err != nil {
    return fmt.Errorf("parsing config: %w", err)
}

// Avoid
if err != nil {
    return err
}
```

### Testing

Use table-driven tests for multiple scenarios:

```go
func TestArrange(t *testing.T) {
    tests := []struct {
        name    string
        nodes   int
        wantErr bool
    }{
        {"single node", 1, false},
        {"multiple nodes", 5, false},
        {"zero nodes", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## When You're Stuck

### Check These In Order

1. **Failing test output** - Read the error message carefully
   ```bash
   make acceptance-test 2>&1 | grep -A 10 "FAIL\|Error"
   ```

2. **Code that's being tested** - Find the actual implementation
   ```bash
   # Find where X is defined
   grep -r "func X\|type X" internal/
   ```

3. **Related tests** - See how similar features are tested
   ```bash
   grep -r "Similar Feature" test/
   ```

4. **Architecture documentation** - Understand the design
   ```bash
   cat docs/architecture/OVERVIEW.md
   ```

5. **Git history** - See how similar features were added
   ```bash
   git log --oneline | grep "similar"
   git show <commit>
   ```

### Common Issues

**Issue**: Tests fail with "undefined step"
**Solution**: Implement step definition in `test/steps/`

**Issue**: "Domain imports from adapters"
**Solution**: Reverse the dependency - use interfaces

**Issue**: "Cannot implement interface"
**Solution**: Check method signature matches exactly (receiver, params, returns)

**Issue**: Coverage is below 80%
**Solution**: Add tests for error paths, edge cases

**Issue**: Test works locally but not in CI
**Solution**: Check relative paths, use `t.TempDir()`, avoid hardcoded paths

## Performance Tips

For agents working on large changes:

1. **Read the CHEAT_SHEET first** - Quick overview of all phases
2. **Keep tests running** - Don't go long without `make acceptance-test`
3. **Commit small changes** - Easier to debug if something breaks
4. **Use interfaces early** - Saves refactoring later
5. **Follow domain-first approach** - Domain → Use Cases → Adapters

## Resources

### In This Repository
- `docs/architecture/OVERVIEW.md` - System design
- `docs/architecture/DECISIONS.md` - Why we made choices
- `docs/architecture/DIAGRAMS.md` - Visual representations
- `docs/ADDING_FEATURES.md` - Step-by-step feature guides
- `docs/refactoring/CHEAT_SHEET.md` - Quick reference

### External
- [Clean Architecture (Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Gherkin Syntax](https://cucumber.io/docs/gherkin/)

## Chat History Tips

### When Starting Fresh
Include this context:
```
I'm working on the layli project (GitHub: dnnrly/layli).
It's a Go CLI tool that generates diagram layouts using Clean Architecture + BDD.

Key facts:
- Domain: pure business logic (100% coverage)
- Use Cases: application logic with port interfaces (100% coverage)
- Adapters: implementations of ports (88-100% coverage)
- Tests: 26 acceptance tests (Gherkin), integration tests, unit tests
- All tests must pass before committing

I want to [YOUR TASK HERE]
```

### When Asking for Help
Be specific:
```
I'm working on [FEATURE]. I need help with [SPECIFIC PROBLEM].

Current state:
- [WHAT YOU'VE DONE]
- [WHAT'S FAILING]
- [ERROR MESSAGE if any]

Here's the relevant code:
[CODE SNIPPET]
```

## Useful Aliases

Add to your shell config for faster development:

```bash
alias lat="make acceptance-test"        # Run acceptance tests
alias ltest="go test ./..."             # Run all tests
alias lcov="go test ./... -cover"       # Show coverage
alias lbuild="go build -o layli ./cmd/layli"  # Build binary
alias lrun="./layli"                    # Run the binary
```

## Sign Off

When your work is complete:

1. All acceptance tests pass: `make acceptance-test` ✅
2. Code is formatted: `go fmt ./...` ✅
3. Code is linted: `go vet ./...` ✅
4. Tests have good coverage: `go test ./... -cover` ✅
5. Commit message explains the change ✅
6. Related documentation is updated ✅

Good luck! 🚀
