# Layli Refactoring Cheat Sheet for AI Agents

**Context Window Optimized - Keep This Open While Working**

---

## 🎯 Mission
Transform layli from monolithic → Clean Architecture while keeping ALL tests green.

---

## 🏗️ Target Architecture

```
Domain (entities)     ← Pure business logic, no deps
   ↑
Use Cases (ports)     ← Application logic, defines interfaces  
   ↑
Adapters (impl)       ← Concrete implementations
   ↑
CLI/Tests            ← Entry points
```

**Rule:** Dependencies point INWARD. Domain knows nothing about adapters.

---

## 📁 File Map: Old → New

| Current File | New Location | Layer |
|--------------|--------------|-------|
| `config.go` | `internal/adapters/config/yaml_parser.go` | Adapter |
| `arrangements.go` | `internal/adapters/layout/*.go` | Adapter |
| `path.go` | `internal/adapters/pathfinding/dijkstra.go` | Adapter |
| `layli.go` (SVG) | `internal/adapters/rendering/svg.go` | Adapter |
| `position.go` | `internal/domain/position.go` | Domain |
| N/A | `internal/usecases/generate_diagram.go` | Use Case |

---

## ✅ Safety Checklist (After EVERY Change)

```bash
make acceptance-test  # MUST pass before continuing
git status            # Review changes
git add <files>
git commit -m "<type>(<scope>): <description>"
```

**If tests fail:** STOP. Debug. Don't proceed.

---

## 🧪 Test Strategy

**Acceptance (Gherkin)** = SPECIFICATION (don't change during refactor)
**Integration** = Multiple layers working together  
**Unit** = Individual components in isolation

```bash
# After each code change
make acceptance-test      # Must pass!

# Before committing
go test ./... -v          # All unit tests
make integration-test     # If available

# Check coverage
go test ./... -cover
```

---

## 📝 Key Interfaces (Ports)

```go
// internal/usecases/ports.go

type ConfigParser interface {
    Parse(path string) (*domain.Diagram, error)
}

type LayoutEngine interface {
    Arrange(diagram *domain.Diagram) error
}

type Pathfinder interface {
    FindPaths(diagram *domain.Diagram) error
}

type Renderer interface {
    Render(diagram *domain.Diagram, outputPath string) error
}
```

**Adapters implement these. Use cases depend on these.**

---

## 🔄 Phase Quick Reference

| Phase | What | Where | Tests Must Pass |
|-------|------|-------|-----------------|
| 0 | Baseline | `docs/` | ✅ Start |
| 1 | Domain | `internal/domain/` | ✅ After each entity |
| 2 | Use Cases | `internal/usecases/` | ✅ After use case |
| 3 | Adapters | `internal/adapters/` | ✅ After EACH adapter |
| 4 | Wire Up | CLI + tests | ✅ After wiring |
| 5 | Testing | `test/` | ✅ Continuously |
| 6 | Docs | `docs/` | ✅ Final |
| 7 | Ship | PR/Tag | ✅ Final |

---

## 💡 Common Patterns

### Domain Entity
```go
// internal/domain/diagram.go
package domain

// Diagram represents [description from Gherkin].
type Diagram struct {
    Nodes  []Node
    Edges  []Edge
    Config DiagramConfig
}

// Validate ensures invariants.
func (d *Diagram) Validate() error {
    // Business rules
}
```

### Use Case
```go
// internal/usecases/generate_diagram.go
package usecases

type GenerateDiagram struct {
    parser     ConfigParser  // Port interface
    layout     LayoutEngine  // Port interface
    pathfinder Pathfinder    // Port interface
    renderer   Renderer      // Port interface
}

func (uc *GenerateDiagram) Execute(...) error {
    // Orchestrate: parse → layout → paths → render
}
```

### Adapter
```go
// internal/adapters/layout/flow_square.go
package layout

import "github.com/dnnrly/layli/internal/domain"

// FlowSquare implements usecases.LayoutEngine.
type FlowSquare struct {}

// Arrange implements the interface.
func (l *FlowSquare) Arrange(diagram *domain.Diagram) error {
    // Algorithm implementation
}
```

---

## 🚨 Red Flags - STOP If You See:

- ❌ Acceptance test fails
- ❌ Domain imports from adapters (dependency violation)
- ❌ Use case imports concrete adapters (should use interfaces)
- ❌ Changing a `.feature` file during refactoring
- ❌ >300 lines in a single file (probably too big)
- ❌ Stuck on same step >2 hours

**Action:** Report the issue, don't push forward blindly.

---

## 🎯 Commit Message Format

```
<type>(<scope>): <subject>

<optional body>

Tests: <result>
```

**Types:** `refactor`, `test`, `docs`, `feat`, `fix`  
**Scopes:** `domain`, `usecases`, `adapters`, `test`, `cli`

**Examples:**
```
refactor(domain): extract Diagram entity

Tests: 45/45 acceptance tests passing

---

refactor(adapters): extract FlowSquare layout adapter

Moved logic from arrangements.go to internal/adapters/layout/flow_square.go
Implements usecases.LayoutEngine interface.

Tests: 45/45 acceptance tests passing
```

---

## 🔍 Debugging Commands

```bash
# Run specific feature
go test -v ./test -godog.tags="@layout"

# Show step output  
go test -v ./test -godog.format=pretty

# Find where something is defined
rg "type Diagram struct"

# See what depends on a package
go mod graph | grep layli

# Check imports
go list -f '{{.Imports}}' ./internal/domain
```

---

## 📊 Progress Tracking Template

```
✅ Phase X.Y Complete: [Title]

Changes:
- Created: internal/xxx/yyy.go
- Modified: test/steps/zzz_steps.go

Tests: XXX/XXX acceptance passing
Coverage: XX%

Next: Phase X.Y+1
```

---

## 🎓 Key BDD Concepts

**Given** = Set up state (→ Domain entities)  
**When** = Perform action (→ Use Case)  
**Then** = Verify outcome (→ Assertions)

```gherkin
Scenario: Generate diagram with flow-square layout
  Given I have a diagram config "test.layli"    # Parse config
  When I generate the diagram                    # Use case
  Then the diagram should be generated           # Assert file exists
  And nodes should be arranged in a grid         # Assert layout
```

**Maps to:**
```go
func (uc *GenerateDiagram) Execute(...) error {
    diagram := uc.parser.Parse(...)      // Given
    uc.layout.Arrange(diagram)            // When
    uc.renderer.Render(diagram, ...)      // When
    return nil                            // Then (if no error)
}
```

---

## 💾 Emergency Rollback

```bash
# Save current work
git stash save "WIP: describe state"

# Return to baseline
git checkout v0.0.14-pre-refactor
git checkout -b refactor/attempt-2

# Or just revert last commit
git revert HEAD
```

---

## ✨ Success = ALL of these ✅

- [ ] All acceptance tests passing (same count as start)
- [ ] Unit test coverage >80%
- [ ] Integration tests added
- [ ] Clean architecture (domain → usecases → adapters)
- [ ] Documentation complete
- [ ] CLI works identically
- [ ] PR merged, v0.1.0 tagged

---

**Remember:** Tests are your safety net. Keep them green! 🟢
