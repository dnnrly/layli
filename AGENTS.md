# AI Agent Guide for Layli

**Layli** is a Go CLI tool that generates SVG diagrams from YAML files using Clean Architecture.

## Quick Start

# 1. Verify tests pass (critical before changes)
```shell
make acceptance-test
```

# 2. Build and test
```shell
make build
./layli examples/hello-world.layli
```

# 3. Key files to understand
- `main.go` - CLI entry
- `internal/domain/diagram.go` - Core entities
- `test/features/layouts.feature` - BDD scenarios

## Architecture

Clean Architecture layers:
- **Domain** (`internal/domain/`) - Pure entities (Diagram, Node, Edge)
- **Use Cases** (`internal/usecases/`) - Port interfaces
- **Adapters** (`internal/adapters/`) - Implementations (YAML parser, layout engines, SVG renderer)

**Key Port Interfaces:**
```go
type ConfigParser interface { Parse(path string) (*domain.Diagram, error) }
type LayoutEngine interface { Arrange(diagram *domain.Diagram) error }
type Pathfinder interface { FindPaths(diagram *domain.Diagram) error }
type Renderer interface { Render(diagram *domain.Diagram, outputPath string) error }
```

## Working on Features

### BDD Workflow

1. **RED**: Add scenario to `test/features/layouts.feature`
2. **RED**: Run `make acceptance-test` (will fail with undefined steps)
3. **IMPLEMENT**: Add step definitions, domain types, adapters
4. **GREEN**: Run `make acceptance-test` until passing
5. **COMMIT**: `git commit -m "feat(scope): description"`

### Adding Layout Algorithms

**Files to modify:**
- `test/features/layouts.feature` - Add scenario
- `test/steps/layout_steps.go` - Step definitions  
- `internal/domain/diagram.go` - Add LayoutType constant
- `internal/adapters/layout/` - Create algorithm
- `internal/adapters/layout/engine.go` - Register in factory

### Common Issues

- **"undefined step"**: Implement in `test/steps/*.go`
- **Domain imports adapters**: Use interfaces in use cases layer
- **Tests fail locally but pass in CI**: Check relative paths, use `t.TempDir()`

## Essential Commands

```bash
# Critical (run frequently)
make acceptance-test          # Must pass before committing
make test                     # Unit tests with coverage
make build                    # Build binary
make lint                     # Linting

# Testing
cd test && go test -v -timeout 20s -tags acceptance -godog.tags="@layout"  # Specific feature
cd test && go test -v -timeout 20s -tags acceptance -godog.format=pretty   # Verbose output

# Debug
make acceptance-test 2>&1 | grep -A 10 "FAIL\|Error"
grep -r "TODO" internal/      # Find todos
```

## Common Patterns

### Test Builder
```go
diagram := support.NewDiagram().
    WithLayout(domain.LayoutFlowSquare).
    AddNode("a", "Node A").
    AddNode("b", "Node B").
    AddEdge("a", "b").
    Build()
```

### Layout Algorithm
```go
type MyLayout struct{}
func (m *MyLayout) Arrange(diagram *domain.Diagram) error {
    // Implementation
    return nil
}
```

## Git Commits

**When to commit:**
- After each working feature
- When tests pass: `make acceptance-test`
- For small, logical changes

**Before committing:**
```bash
make acceptance-test  # Must pass
make lint            # No errors
git diff --cached    # Review changes
```

**Commit format:**
```bash
feat(layout): add circular algorithm
fix(pathfinding): resolve edge crossing
docs: update README examples
```

**Golden rule: If `make acceptance-test` fails, fix it or revert before committing.**

## Getting Help

- Read failing test output carefully
- Check existing layout algorithms for patterns
- Use `--show-grid` flag for visual debugging
- Never commit failing tests

---

**Rule**: Always run `make acceptance-test` before committing.
