# AI Agent Guide for Layli

This guide provides comprehensive information for AI coding agents working on the **layli** project - a Go CLI tool that generates SVG diagrams from YAML configuration files using Clean Architecture principles.

## Project Overview

**Layli** is a diagram generation tool that:
- Parses YAML configuration files defining nodes and edges
- Arranges nodes using various layout algorithms (flow-square, topological, Tarjan, absolute, random-shortest)
- Routes edges between nodes using pathfinding algorithms
- Renders SVG output with automatic layout optimization

### Key Characteristics
- **Language**: Go 1.21+
- **Architecture**: Clean Architecture with Ports and Adapters (Hexagonal)
- **Testing**: BDD with Gherkin scenarios + unit tests + integration tests
- **Coverage**: 97% overall code coverage
- **CLI**: Built with Cobra framework
- **Output**: SVG diagrams with optional grid overlay

## Quick Start for Agents

### 1. Project Understanding (15 minutes)

**Essential Reading Order:**
1. `README.md` - Project purpose, usage examples, principles
2. `docs/architecture/OVERVIEW.md` - Clean Architecture design
3. `docs/FEATURE_MAP.md` - Feature-to-code mapping
4. `docs/AGENT_GUIDE.md` - Detailed development workflow

**Key Files to Examine:**
- `main.go` - CLI entry point and command structure
- `internal/domain/diagram.go` - Core domain entities
- `test/features/layouts.feature` - BDD scenarios for layout features
- `go.mod` - Dependencies and project structure

### 2. Environment Setup

```bash
# Verify all tests pass (critical before making changes)
make acceptance-test

# Run specific test suites
make test                           # All unit tests with coverage
make acceptance-test                # Acceptance tests (BDD)
make ci-test                        # CI test suite

# Build the application
make build
./layli examples/hello-world.layli
```

### 3. Architecture Deep Dive

**Clean Architecture Layers:**

```
┌─────────────────────────────────────┐
│           External (CLI)            │ ← main.go, Cobra commands
└──────────────────┬──────────────────┘
                   │
                   ↓
         ┌─────────────────────┐
         │  Composition Root   │ ← internal/composition/
         │    (Wiring)         │   NewGenerateDiagram()
         └──────────┬──────────┘
                    │
                    ↓
         ┌─────────────────────────────────┐
         │   Use Case Layer                │ ← internal/usecases/
         │ (GenerateDiagram orchestrates)  │   Port interfaces
         └──────────┬──────────────────────┘
                    │
        ┌───────────┼──────────┬──────────┐
        ↓           ↓          ↓          ↓
    ┌────────┐ ┌────────┐ ┌─────────┐ ┌────────┐
    │ Config │ │ Layout │ │ Path    │ │Render  │ ← internal/adapters/
    │Adapter │ │Adapter │ │Adapter  │ │Adapter │   Implementations
    └────┬───┘ └────┬───┘ └────┬────┘ └───┬────┘
         │          │          │          │
         └──────────┼──────────┼──────────┘
                    ↓
         ┌──────────────────────┐
         │    Domain Layer      │ ← internal/domain/
         │  (Pure Entities)     │   Diagram, Node, Edge
         └──────────────────────┘
```

**Key Domain Entities:**
- `Diagram` - Complete diagram with nodes, edges, config
- `Node` - Individual diagram node with position/dimensions
- `Edge` - Connection between nodes with path routing
- `DiagramConfig` - Layout settings, dimensions, algorithms
- `Position` - Spatial coordinates for layout calculations

**Port Interfaces (in `internal/usecases/`):**
```go
type ConfigParser interface { Parse(path string) (*domain.Diagram, error) }
type LayoutEngine interface { Arrange(diagram *domain.Diagram) error }
type Pathfinder interface { FindPaths(diagram *domain.Diagram) error }
type Renderer interface { Render(diagram *domain.Diagram, outputPath string) error }
```

## Working on Features

### BDD-First Development Workflow

**1. RED Phase - Write Failing Test**
```bash
# Add scenario to test/features/layouts.feature
Scenario: Circular layout arranges nodes in a circle
  Given I have a diagram with 6 nodes
  When I arrange using "circular" layout
  Then the nodes should be positioned in a circle

# Run tests - they will fail
make acceptance-test
# Output: undefined step: "the nodes should be positioned in a circle"
```

**2. IMPLEMENT Phase - Make Tests Pass**
```bash
# 1. Implement step definitions in test/steps/
# 2. Add domain types if needed
# 3. Create adapter implementing port interface
# 4. Register in composition root
# 5. Run tests until green
make acceptance-test
```

**3. REFACTOR Phase - Improve Code**
```bash
# Improve without changing behavior
make acceptance-test  # Must still pass
```

**4. COMMIT Phase**
```bash
git add .
git commit -m "feat(layout): add circular layout algorithm

- Create Circular layout adapter implementing LayoutEngine
- Add LayoutCircular constant to domain
- Register in LayoutAdapter factory
- Include acceptance tests and example

Tests: 26/26 acceptance tests passing"
```

### Common Feature Types

#### Adding New Layout Algorithms
**Files to Modify:**
- `test/features/layouts.feature` - Add scenario
- `test/steps/layout_steps.go` - Step definitions
- `internal/domain/diagram.go` - Add LayoutType constant
- `internal/adapters/layout/` - Create algorithm implementation
- `internal/adapters/layout/engine.go` - Register in factory

**Example Implementation:**
```go
// internal/adapters/layout/circular.go
type Circular struct{}

func (c *Circular) Arrange(diagram *domain.Diagram) error {
    // Position nodes in circle pattern
    return nil
}
```

#### Adding New Output Formats
**Files to Modify:**
- `test/features/rendering.feature` - Add scenario
- `internal/adapters/rendering/` - Create renderer
- `internal/composition/generate_diagram.go` - Wire up
- `main.go` - Add CLI flags

#### Adding New Configuration Formats
**Files to Modify:**
- `internal/adapters/config/` - Create parser
- `internal/composition/generate_diagram.go` - Wire up
- Test with various config formats

### Debugging and Troubleshooting

#### When Tests Fail
```bash
# 1. Check specific test output
make acceptance-test 2>&1 | grep -A 10 "FAIL\|Error"

# 2. Run specific scenario
cd test && go test -v -timeout 20s -tags acceptance -godog.tags="@layout"

# 3. Check test coverage
make ci-test
make coverage-report

# 4. Run with verbose output
cd test && go test -v -timeout 20s -tags acceptance -godog.format=pretty
```

#### Common Issues and Solutions

**Issue**: "undefined step" error
**Solution**: Implement step definition in appropriate `test/steps/*.go` file

**Issue**: Domain imports from adapters (violates Clean Architecture)
**Solution**: Reverse dependency using interfaces in use cases layer

**Issue**: Low test coverage
**Solution**: Add tests for error paths and edge cases in `*_test.go` files

**Issue**: Tests pass locally but fail in CI
**Solution**: Check relative paths, use `t.TempDir()`, avoid hardcoded file paths

## Code Style and Patterns

### Naming Conventions
- **Packages**: lowercase, short, no underscores (`config`, `layout`, `pathfinding`)
- **Types**: PascalCase (`DiagramConfig`, `LayoutEngine`, `SVGRenderer`)
- **Functions**: camelCase (`NewDiagram()`, `Validate()`, `Arrange()`)
- **Constants**: UPPER_SNAKE_CASE for enums (`LayoutFlowSquare`, `LayoutTopoSort`)
- **Interfaces**: End with noun/verb (`Reader`, `Writer`, `Parser`, `Engine`)

### File Structure Pattern
```go
package layout

import (
    "fmt"
    "github.com/dnnrly/layli/internal/domain"
)

// Type definitions
type FlowSquare struct {
    // fields
}

// Constructor
func NewFlowSquare() *FlowSquare {
    return &FlowSquare{}
}

// Interface implementation
func (fs *FlowSquare) Arrange(diagram *domain.Diagram) error {
    // implementation
}

// Helper functions
func calculateGrid(count int) (rows, cols int) {
    // ...
}
```

### Error Handling Pattern
```go
// Good - wrap with context
if err := parser.Parse(path); err != nil {
    return fmt.Errorf("parsing config: %w", err)
}

// Avoid - bare error returns
if err != nil {
    return err
}
```

### Testing Pattern
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
            // test logic using support.NewDiagram() builder
        })
    }
}
```

## Key Commands Reference

```bash
# Essential commands (run frequently)
make acceptance-test          # Must always pass before committing
make test                     # All unit tests with coverage
make ci-test                  # CI test suite with coverage
make coverage-report          # Generate merged coverage report

# Development commands
make build                    # Build binary
make lint                     # Run linting
make clean                    # Clean build artifacts

# Testing commands
cd test && go test -v -timeout 20s -tags acceptance -godog.tags="@layout"  # Specific feature
cd test && go test -v -timeout 20s -tags acceptance -godog.format=pretty   # Detailed output
make ci-test && make coverage-report  # Coverage analysis

# Investigation commands
grep -r "TODO" internal/                 # Find todos
go list -f '{{.Imports}}' ./internal/domain  # Check dependencies
git log --oneline | grep "layout"        # Feature history
```

## Project-Specific Patterns

### Test Data Builder Pattern
```go
// Use the fluent builder for test setup
diagram := support.NewDiagram().
    WithLayout(domain.LayoutFlowSquare).
    AddNode("a", "Node A").
    AddNode("b", "Node B").
    AddEdge("a", "b").
    Build()
```

### Layout Algorithm Pattern
```go
// All layout algorithms implement LayoutEngine
type LayoutEngine interface {
    Arrange(diagram *domain.Diagram) error
}

// Registration in factory
engines: map[domain.LayoutType]LayoutEngine{
    domain.LayoutFlowSquare: NewFlowSquare(),
    domain.LayoutCircular:   NewCircular(),
}
```

### SVG Assertion Pattern
```go
// Common test assertions for SVG output
func (ctx *layoutContext) theNumberOfNodesIs(count int) error {
    nodes := ctx.svgDoc.FindElements("rect.node")
    if len(nodes) != count {
        return fmt.Errorf("expected %d nodes, got %d", count, len(nodes))
    }
    return nil
}
```

## Performance and Optimization

### For Large Changes
1. **Read First**: Understand existing patterns before modifying
2. **Test Often**: Run `make acceptance-test` frequently
3. **Commit Small**: Make focused, reviewable commits
4. **Use Interfaces**: Design for extensibility from the start
5. **Domain First**: Implement domain entities before adapters

### Memory and Performance Considerations
- Domain entities are lightweight structs
- Layout algorithms work in-place on diagram
- SVG rendering streams output to avoid large memory usage
- Pathfinding uses Dijkstra's algorithm with grid optimization

## Integration Points

### External Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/ajstarks/svgo` - SVG generation
- `github.com/antchfx/xmlquery` - XML parsing (for to-absolute command)
- `gopkg.in/yaml.v3` - YAML configuration parsing
- `github.com/cucumber/godog` - BDD test framework

### File System Integration
- Input: `.layli` YAML files
- Output: `.svg` files (with optional `.layli` for absolute positioning)
- Test fixtures: `test/fixtures/inputs/`

### CLI Integration Points
- Main command: `layli [flags] [input.layli]`
- Sub-command: `layli to-absolute [svg-file]`
- Flags: `--output`, `--layout`, `--show-grid`

## Git Workflow

### Commit Message Format

Follow conventional commits format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting (no logic change)
- `refactor`: Code improvement without feature change
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependencies, build

**Examples:**
```bash
# Good commit message
feat(layout): add circular layout algorithm

- Create Circular layout adapter implementing LayoutEngine
- Add LayoutCircular constant to domain
- Register in LayoutAdapter factory
- Include acceptance tests and example

Tests: 26/26 acceptance tests passing

# Bug fix commit
fix(pathfinding): resolve edge crossing in dense graphs

- Update Dijkstra algorithm to avoid overlapping paths
- Add validation for minimum node spacing
- Fixes issue where edges would cross in 3x3 grids

Closes #123

# Simple change
fix: correct typo in CLI help text
```

### Commit Workflow

**1. Stage Changes Strategically**
```bash
# Review what you're changing
git status
git diff

# Stage related changes together
git add internal/domain/diagram.go          # Domain changes
git add test/features/layouts.feature       # Feature test
git add test/steps/layout_steps.go         # Step definitions
```

**2. Write Commit Message**
```bash
# Use editor for detailed messages
git commit

# Or inline for simple changes
git commit -m "fix: correct node positioning calculation"
```

**3. Verify Commit**
```bash
# Show last commit
git show --stat HEAD

# Show commit message only
git log -1 --pretty=format:"%s%n%n%b"

# Check recent commits
git log --oneline -5
```

### Branch Strategy

**Feature Branch Workflow:**
```bash
# Create feature branch
git checkout -b feature/circular-layout

# Work on feature (make small, focused commits)
git add .
git commit -m "feat(domain): add LayoutCircular constant"

git add .
git commit -m "feat(layout): implement circular algorithm"

git add .
git commit -m "test: add circular layout scenarios"

# Push branch
git push -u origin feature/circular-layout

# Create pull request
# (Use GitHub UI or gh cli)
gh pr create --title "feat: add circular layout algorithm" --body "Adds circular layout..."
```

### Commit Best Practices

**DO:**
- ✅ Make commits focused and logical units
- ✅ Write descriptive commit messages (50 char max subject)
- ✅ Include body paragraphs explaining "why" not "what"
- ✅ Reference issue numbers when applicable
- ✅ Ensure tests pass before committing
- ✅ Use conventional commit format
- ✅ Break large changes into smaller, reviewable commits

**DON'T:**
- ❌ Commit with failing tests
- ❌ Use vague messages like "fixed stuff" or "wip"
- ❌ Mix unrelated changes in one commit
- ❌ Commit large, unreadable diffs
- ❌ Include temporary/debug code in commits
- ❌ Use force push on shared branches

### Pre-Commit Checklist

Before each commit, verify:

```bash
# 1. Tests pass
make acceptance-test

# 2. Code is formatted
go fmt ./...

# 3. No lint errors
make lint

# 4. Review staged changes
git diff --cached

# 5. Check for sensitive data
git diff --cached | grep -i "password\|key\|secret"
```

### Undoing Changes

**When you make a mistake:**
```bash
# Unstage last added file
git reset HEAD internal/domain/diagram.go

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1

# Amend last commit (add forgotten file)
git add forgotten_file.go
git commit --amend

# Edit last commit message
git commit --amend --no-edit
```

### History Cleanup

**Before creating pull request:**
```bash
# Squash related commits (use with care)
git rebase -i HEAD~3

# Interactive rebase options:
# pick = use commit
# squash = combine with previous commit
# reword = edit commit message
# fixup = combine with previous (discard message)

# Example: squash 3 commits into one
# pick abc123 Initial work
# squash def456 Add tests
# squash ghi789 Fix bugs
# -> Results in one clean commit
```

## Quality Gates

### Before Committing
- [ ] All acceptance tests pass: `make acceptance-test`
- [ ] Code formatted: `go fmt ./...`
- [ ] No lint errors: `make lint`
- [ ] Coverage maintained: `make ci-test && make coverage-report`
- [ ] Staged changes reviewed: `git diff --cached`
- [ ] Commit message follows conventions
- [ ] Examples updated (if applicable)
- [ ] Documentation updated (if needed)

### Code Review Checklist
- [ ] Follows Clean Architecture principles
- [ ] Implements appropriate interfaces
- [ ] Has comprehensive tests
- [ ] Error handling is proper
- [ ] No hardcoded paths or values
- [ ] Domain invariants are validated

## Getting Help

### When Stuck
1. **Read the failing test** - Error messages are descriptive
2. **Check similar implementations** - Look at existing layout algorithms
3. **Review architecture docs** - `docs/architecture/OVERVIEW.md`
4. **Examine git history** - See how similar features were added
5. **Run with debug output** - Use `--show-grid` flag for visual debugging

### Effective Communication
When asking for help, provide:
- What you're trying to accomplish
- What you've already tried
- Specific error messages or test failures
- Relevant code snippets
- Current test results

## Project Evolution

### Current State
- **Stable**: Core layout algorithms (flow-square, topological, Tarjan)
- **Stable**: SVG rendering with path optimization
- **Stable**: CLI interface with Cobra
- **Known Issue**: `random-shortest-square` layout has failing tests
- **Active**: BDD test suite with 26 scenarios

### Future Directions
Potential areas for enhancement:
- Additional layout algorithms (circular, hierarchical, force-directed)
- Alternative output formats (PNG, JSON, DOT/Graphviz)
- Interactive web interface
- Performance optimizations for large diagrams
- Advanced styling and theming

---

**Remember**: The golden rule is **never commit failing tests**. If tests break, either fix the issue or revert the changes. The BDD suite is the project's source of truth for behavior.

Good luck with your development! 🚀
