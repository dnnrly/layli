# Day 1: Layli Refactoring Kickoff

**Copy this entire prompt to your AI agent (Windsurf/Cursor/Aider) to begin**

---

## 🚀 Mission Briefing

I'm refactoring the layli project from a monolithic Go application to a clean, layered architecture with BDD-first design. This is a complex, multi-week effort that requires careful, incremental changes while keeping all tests green.

**Repository:** https://github.com/dnnrly/layli  
**Goal:** Clean Architecture with Domain → Use Cases → Adapters  
**Constraint:** All Gherkin acceptance tests must pass after EVERY change  
**Duration:** 3-4 weeks (7 phases)

---

## 📋 Your First Tasks (Phase 0: Establish Baseline)

### Task 1: Understand the Project

Clone the repository and familiarize yourself with its structure:

```bash
# If not already cloned
git clone https://github.com/dnnrly/layli.git
cd layli

# Create our refactoring branch
git checkout -b refactor/layered-architecture
```

Now analyze the project:

1. **Run the acceptance tests** to establish our baseline:
   ```bash
   make acceptance-test
   ```
   
   Save the output and tell me:
   - How many scenarios are there?
   - How many pass/fail?
   - What test framework is being used? (I expect godog/Gherkin)

2. **Examine the test structure**:
   ```bash
   # List feature files
   ls -la test/features/
   
   # Show me one feature file
   cat test/features/layout.feature  # or whatever exists
   
   # List step definitions
   ls -la test/steps/
   ```
   
   Summarize:
   - What features exist (layout, pathfinding, config, etc.)?
   - How are step definitions organized?

3. **Analyze current code organization**:
   ```bash
   # List all Go files in root
   ls -la *.go
   
   # Check internal structure
   ls -la internal/ 2>/dev/null || echo "No internal/ directory yet"
   
   # Check what packages exist
   go list ./...
   ```
   
   Create a quick map:
   - Which files contain configuration parsing?
   - Which files contain layout algorithms?
   - Which files contain pathfinding logic?
   - Which files contain rendering (SVG) logic?

### Task 2: Create Initial Documentation

Create `docs/CURRENT_STATE.md` with the following structure:

```markdown
# Current State Analysis

**Date:** [Today's date]
**Branch:** refactor/layered-architecture
**Baseline Tag:** v0.0.14-pre-refactor

## Test Baseline

- **Acceptance Tests:** XXX scenarios
  - Passing: XXX
  - Failing: XXX
- **Test Framework:** godog/Gherkin
- **Feature Files:** 
  - test/features/layout.feature
  - test/features/... [list all]

## Current Code Organization

### Root Directory Files
- `config.go` - [brief description of what it does]
- `arrangements.go` - [brief description]
- `layout.go` - [brief description]
- `path.go` - [brief description]
- `layli.go` - [brief description]
- [... list all .go files in root]

### Current Package Structure
```
layli/
├── cmd/layli/
│   └── main.go
├── algorithms/
├── pathfinder/dijkstra/
├── test/features/
├── test/steps/
└── [other directories]
```

### Key Dependencies
- YAML parsing: [which library?]
- SVG generation: [how is it done?]
- Pathfinding: [Dijkstra implementation location]

## Code Complexity
- Total lines of code: ~XXXX
- Largest file: XXX.go (XXX lines)
- Average file size: ~XXX lines
- Root directory .go files: XX files

## Observations
- [Note any immediate architectural issues you see]
- [Note test organization quality]
- [Note any obvious refactoring challenges]
```

### Task 3: Create Feature-to-Code Mapping

Create `docs/FEATURE_MAP.md`:

For each feature file in `test/features/`, map it to the current implementation. Example structure:

```markdown
# Feature to Code Mapping

## Layout Features

**Feature File:** `test/features/layout.feature`

**Scenarios:**
1. Scenario: Flow-square layout arranges nodes in grid
   - Description: [what it tests]
2. Scenario: ... [list all]

**Current Implementation:**
- Primary files:
  - `arrangements.go` - contains [list key functions]
  - `layout.go` - contains [list key functions]
- Step definitions: `test/steps/layout_steps.go`
- Key types: [list important structs/types]

**Refactoring Target:**
- Domain: `internal/domain/diagram.go`, `internal/domain/node.go`
- Use Case: `internal/usecases/generate_diagram.go`
- Adapter: `internal/adapters/layout/flow_square.go`, etc.

---

## Pathfinding Features

**Feature File:** `test/features/pathfinding.feature` (if exists)

[Similar structure for each feature]

---

## Configuration Features

[Similar structure]
```

### Task 4: Tag Baseline

Only do this AFTER all tests are passing:

```bash
# Verify tests pass
make acceptance-test

# If passing, tag the baseline
git tag v0.0.14-pre-refactor
git push origin v0.0.14-pre-refactor

# Commit our documentation
git add docs/
git commit -m "docs: establish refactoring baseline

- Document current state of codebase
- Create feature-to-code mapping
- Establish test baseline

Tests: All acceptance tests passing"
```

### Task 5: Report Back

Provide me with a summary:

```
📊 Phase 0 Complete - Baseline Established

✅ Acceptance Tests: XXX/XXX passing
✅ Feature Files: XX files identified
✅ Current Architecture: [monolithic/partially structured]
✅ Documentation: Created CURRENT_STATE.md and FEATURE_MAP.md
✅ Baseline Tagged: v0.0.14-pre-refactor

Key Findings:
- [Major observations about code organization]
- [Test coverage assessment]
- [Potential challenges identified]

Current Structure:
- Root directory: XX .go files (~XXXX total LOC)
- Primary packages: [list]
- Test organization: [godog/Gherkin with X feature files]

Ready to begin Phase 1: Extract Domain Layer
```

---

## 📚 Reference Documents

I've provided you with several reference documents. Keep these accessible:

1. **AGENT_REFACTORING_PROMPTS.md** - Detailed prompts for each phase
2. **QUICK_REFERENCE_CHECKLIST.md** - Phase-by-phase checklist
3. **CHEAT_SHEET.md** - One-page quick reference (keep this visible)

**Before starting each new phase, read the corresponding prompt from AGENT_REFACTORING_PROMPTS.md**

---

## ⚠️ Critical Rules

1. **NEVER proceed if acceptance tests fail** - Fix or revert first
2. **Commit after each successful step** - Small, focused commits
3. **DON'T modify .feature files** during refactoring - They are the specification
4. **Run tests after EVERY code change** - `make acceptance-test`
5. **Ask questions if uncertain** - Don't guess on critical decisions

---

## 🎯 Success Criteria for Today (Phase 0)

By end of Day 1, we should have:
- ✅ Understanding of current codebase
- ✅ Baseline test count established
- ✅ Documentation created
- ✅ Baseline tagged
- ✅ Clean git status

**Time Estimate:** 1-2 hours for Phase 0

---

## 🚦 Start Here

Begin with Task 1 above. Work through each task sequentially. After completing Task 5 (Report Back), we'll move to Phase 1: Extract Domain Layer.

**Remember:** Quality over speed. Take your time to understand the codebase before making changes.

---

## 💬 Communication Template

As you work, provide updates using this format:

```
🔄 Working on: Task X - [description]

Progress:
- [x] Step 1
- [x] Step 2
- [ ] Step 3 (current)

Current file: [path/to/file.go]
Tests: XXX/XXX passing
Status: [on track / blocked / question]
```

---

## 🆘 If You Get Stuck

If you encounter issues:

1. **Report the specific error:**
   ```
   Error in: [file/command]
   Error message: [exact error]
   What I was trying to do: [description]
   ```

2. **Show me the context:**
   ```bash
   git status
   git diff [file]
   make acceptance-test  # Output
   ```

3. **Ask specific questions:**
   - "Should X be in domain or adapters?"
   - "This function does Y and Z - how to split it?"
   - "Tests fail with [error] - what's wrong?"

---

## 🎬 Begin Now

Start with **Task 1: Understand the Project**. Take your time, be thorough, and establish a solid foundation. Good luck! 🚀

---

**Next Steps After Phase 0:**
Once Phase 0 is complete, I'll provide you with the Phase 1 prompt from AGENT_REFACTORING_PROMPTS.md to begin extracting the domain layer.
