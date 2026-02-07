# Layli Refactoring Quick Reference Checklist

**Use this alongside the detailed AGENT_REFACTORING_PROMPTS.md**

## Setup
- [ ] Create branch: `git checkout -b refactor/layered-architecture`
- [ ] Run baseline tests: `make acceptance-test` (save output)
- [ ] Tag baseline: `git tag v0.0.14-pre-refactor`

---

## Phase 0: Baseline (1-2 hours)
- [ ] 0.1: Understand current state, document in `docs/CURRENT_STATE.md`
- [ ] 0.2: Create feature-to-code map in `docs/FEATURE_MAP.md`
- [ ] 0.3: Verify all tests pass, tag baseline

---

## Phase 1: Domain Layer (2-3 days)
- [ ] 1.1: Create `internal/domain/` structure with doc.go
- [ ] 1.2: Extract `Diagram` entity → `internal/domain/diagram.go`
- [ ] 1.3: Extract `Node` entity → `internal/domain/node.go`
- [ ] 1.4: Extract `Edge`, `Position`, `Path` → `internal/domain/*.go`
- [ ] 1.5: Add comprehensive unit tests for domain layer

**Checkpoint:** ✅ All acceptance tests still passing

---

## Phase 2: Use Case Layer (3-4 days)
- [ ] 2.1: Create `internal/usecases/` structure with doc.go
- [ ] 2.2: Define port interfaces → `internal/usecases/ports.go`
- [ ] 2.3: Create `GenerateDiagram` use case → `internal/usecases/generate_diagram.go`
- [ ] 2.4: Add use case tests with mocks

**Checkpoint:** ✅ All acceptance tests still passing

---

## Phase 3: Adapter Layer (5-7 days)
- [ ] 3.1: Create `internal/adapters/` structure
- [ ] 3.2: Extract config parser → `internal/adapters/config/yaml_parser.go`
- [ ] 3.3: Extract FlowSquare layout → `internal/adapters/layout/flow_square.go`
- [ ] 3.4: Extract remaining layouts (TopoSort, Tarjan, Absolute) + factory
- [ ] 3.5: Extract pathfinding → `internal/adapters/pathfinding/dijkstra.go`
- [ ] 3.6: Extract SVG renderer → `internal/adapters/rendering/svg.go`

**Checkpoint after EACH adapter:** ✅ Acceptance tests passing

---

## Phase 4: Integration (2-3 days)
- [ ] 4.1: Update step definitions to use new architecture
- [ ] 4.2: Update CLI (`cmd/layli/main.go`) + create layout selector
- [ ] 4.3: Remove old code from root directory (config.go, arrangements.go, etc.)

**Checkpoint:** ✅ All tests passing, CLI works

---

## Phase 5: Testing (2-3 days)
- [ ] 5.1: Create test helpers → `test/support/diagram_builder.go`
- [ ] 5.2: Add integration tests → `test/integration/generation_test.go`
- [ ] 5.3: Improve coverage to >80% across all packages

**Checkpoint:** ✅ Coverage >80%, all tests passing

---

## Phase 6: Documentation (2 days)
- [ ] 6.1: Create architecture docs → `docs/architecture/OVERVIEW.md`
- [ ] 6.2: Create developer guides → `docs/ADDING_FEATURES.md`, `docs/AGENT_GUIDE.md`
- [ ] 6.3: Update README with architecture section

---

## Phase 7: Final Validation (1 day)
- [ ] 7.1: Run comprehensive test suite, create completion report
- [ ] 7.2: Create PR with detailed description
- [ ] 7.3: Tag v0.1.0, update CHANGELOG, create release

---

## Test Commands (Run After Each Phase)

```bash
# Acceptance tests (MUST pass after every change)
make acceptance-test

# Unit tests
go test ./... -v

# Integration tests  
make integration-test

# Coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Lint
make lint

# Manual CLI test
go run cmd/layli/main.go demo.layli --show-grid
```

---

## Commit Template

```
<type>(<scope>): <subject>

<body>

Tests: <test results>
```

**Types:** refactor, test, docs, feat, fix
**Scopes:** domain, usecases, adapters, test, docs

**Example:**
```
refactor(domain): extract Diagram entity

Move Diagram and DiagramConfig from config.go to internal/domain/diagram.go.
Adds validation method and comprehensive documentation.

Tests: All acceptance tests passing (45/45)
```

---

## Red Flags 🚨

**STOP and ask for help if:**
- ❌ Acceptance tests fail after a change
- ❌ You're unsure which file to modify
- ❌ You've been on same step for >2 hours
- ❌ Test coverage drops below 70%
- ❌ You need to change a feature file during refactoring

---

## Success Metrics

| Metric | Before | Target |
|--------|--------|--------|
| Acceptance Tests | XX/XX | XX/XX (same) |
| Unit Test Coverage | ~XX% | >80% |
| Integration Tests | 0 | >5 |
| Root Directory Files | ~15 | <5 |
| Avg File Size | >300 LOC | <200 LOC |
| Time to Add Feature | 2-3 days | 2-4 hours |

---

## Emergency Commands

```bash
# Rollback last commit
git revert HEAD

# See what changed
git diff HEAD~1

# Return to baseline
git checkout v0.0.14-pre-refactor

# Stash current work
git stash save "WIP: describe what you were doing"

# See stashed changes
git stash show -p stash@{0}
```

---

## Phase Estimates

| Phase | Duration | Complexity |
|-------|----------|------------|
| 0: Baseline | 1-2 hours | Easy |
| 1: Domain | 2-3 days | Medium |
| 2: Use Cases | 3-4 days | Medium |
| 3: Adapters | 5-7 days | Hard |
| 4: Integration | 2-3 days | Medium |
| 5: Testing | 2-3 days | Easy |
| 6: Documentation | 2 days | Easy |
| 7: Final | 1 day | Easy |
| **Total** | **3-4 weeks** | - |

---

## Daily Standup Template

```
Yesterday:
- Completed: Phase X.Y - [description]
- Committed: X commits
- Tests: XXX/XXX passing

Today:
- Working on: Phase X.Y+1 - [description]
- Goal: [specific deliverable]

Blockers:
- [None / describe issue]
```

---

## File Organization Quick Reference

```
layli/
├── cmd/layli/main.go          ← Thin CLI entry point
├── internal/
│   ├── domain/                ← Pure entities (NO deps)
│   │   ├── diagram.go
│   │   ├── node.go
│   │   ├── edge.go
│   │   ├── position.go
│   │   └── path.go
│   ├── usecases/              ← Application logic + ports
│   │   ├── ports.go
│   │   └── generate_diagram.go
│   └── adapters/              ← Concrete implementations
│       ├── config/            ← YAML/JSON parsers
│       ├── layout/            ← Layout algorithms
│       ├── pathfinding/       ← Path algorithms
│       └── rendering/         ← SVG/PNG/etc
├── test/
│   ├── features/              ← Gherkin (DON'T CHANGE)
│   ├── steps/                 ← Step definitions
│   ├── support/               ← Test helpers
│   └── integration/           ← Integration tests
└── docs/
    ├── architecture/          ← Architecture docs
    ├── ADDING_FEATURES.md
    └── AGENT_GUIDE.md
```

---

## Key Rules

1. **Always run tests after changes**
2. **Commit after each successful step**
3. **Never modify feature files during refactoring**
4. **One logical change per commit**
5. **Ask before making breaking changes**
6. **Document as you go, not at the end**
7. **If tests fail, STOP and debug immediately**

---

**Remember:** This is a marathon, not a sprint. Take breaks, commit frequently, and keep tests green! 🚀
