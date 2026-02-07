# Using Refactoring Prompts with Amp

This guide explains how to use the refactoring prompts with Amp (or similar AI coding agents).

## 🎯 What is Amp?

Amp is an AI-powered coding assistant that can read your codebase, execute commands, make changes, and help with complex refactoring tasks.

## 📋 Prerequisites

Before starting:
- [ ] Amp is installed and configured
- [ ] You're in the layli project directory
- [ ] All tests currently pass: `make acceptance-test`
- [ ] You have a clean git working directory
- [ ] You've reviewed the refactoring plan in `AGENT_REFACTORING_PROMPTS.md`

## 🚀 Quick Start Guide

### Step 1: Initial Setup

```bash
# Create refactoring branch
git checkout -b refactor/layered-architecture

# Tag the current state
git tag v0.0.14-pre-refactor

# Verify tests pass
make acceptance-test
# Ensure all tests are green before proceeding!
```

### Step 2: Start Amp

```bash
# Start Amp in your layli project directory
amp

# Or if Amp has a specific command:
amp chat
# (adjust based on your Amp installation)
```

### Step 3: Provide Initial Context to Amp

**First Message to Amp:**

```
I'm refactoring the layli project using a structured, multi-phase plan. 
I have detailed prompts for each phase in docs/refactoring/.

Current status:
- Project: layli (diagram generation tool in Go)
- Testing: BDD with Gherkin/godog
- Goal: Transform to Clean Architecture (Domain → Use Cases → Adapters)
- Constraint: All acceptance tests must pass after every change

Please read the following files to understand the project and plan:
1. docs/refactoring/CHEAT_SHEET.md - Quick reference
2. docs/refactoring/DAY_1_KICKOFF.md - Phase 0 instructions

After reading, confirm you understand the approach and are ready to begin Phase 0.
```

### Step 4: Execute Phase 0 (Baseline)

**Second Message to Amp:**

```
Execute the tasks in docs/refactoring/DAY_1_KICKOFF.md starting with "Task 1: Understand the Project".

Work through each task sequentially:
- Task 1: Understand the Project
- Task 2: Create Initial Documentation
- Task 3: Create Feature-to-Code Mapping
- Task 4: Tag Baseline
- Task 5: Report Back

After each task, report:
- What you completed
- Any findings
- Test status (make acceptance-test results)

CRITICAL: If any tests fail, STOP and report the failure before proceeding.
```

### Step 5: Subsequent Phases

For each subsequent phase, provide Amp with the specific prompt from `AGENT_REFACTORING_PROMPTS.md`.

**Example for Phase 1:**

```
We're ready to begin Phase 1: Extract Domain Layer.

Please read and execute the prompts in docs/refactoring/AGENT_REFACTORING_PROMPTS.md 
starting from "Phase 1: Extract Domain Layer (Week 1)".

Start with Prompt 1.1: Create Domain Package Structure.

After each prompt:
1. Execute the steps
2. Run: make acceptance-test
3. If tests pass: commit the changes
4. If tests fail: STOP and report
5. Provide a status update

Proceed with Prompt 1.1 now.
```

## 💡 Best Practices for Working with Amp

### 1. Break Down Large Prompts

Instead of feeding Amp an entire phase at once, provide one prompt at a time:

```
Execute Prompt 1.1 from AGENT_REFACTORING_PROMPTS.md:
"Create Domain Package Structure"

Follow the steps exactly as written. Report when complete.
```

### 2. Verify After Each Step

After Amp makes changes:

```
Run these verification steps:
1. make acceptance-test
2. go test ./... -v
3. git status

Report the results before proceeding to the next prompt.
```

### 3. Keep Amp Focused

If Amp starts to deviate from the plan:

```
Stop. Please stick to the exact steps in Prompt X.Y.

Do not:
- Skip steps
- Combine multiple prompts
- Make changes beyond what's described

Proceed with only the steps in Prompt X.Y.
```

### 4. Use the Cheat Sheet as Context

Periodically remind Amp of the architecture rules:

```
Before proceeding, re-read docs/refactoring/CHEAT_SHEET.md 
to ensure you understand:
- Target architecture (Domain → Use Cases → Adapters)
- Dependency rules
- Testing requirements

Confirm understanding, then proceed with next prompt.
```

### 5. Track Progress

Update the checklist after each prompt:

```
Update docs/refactoring/QUICK_REFERENCE_CHECKLIST.md:
- Mark Prompt X.Y as complete ✅
- Note current test status
- Estimate time remaining

Then commit the updated checklist.
```

## 🔄 Typical Workflow Pattern

Here's the rhythm you'll establish:

```
1. You: "Execute Prompt X.Y from AGENT_REFACTORING_PROMPTS.md"
   
2. Amp: [Reads prompt, executes steps, makes changes]
   
3. Amp: "Changes made. Running tests..."
   
4. Amp: "Tests passing: XXX/XXX. Ready to commit."
   
5. You: "Commit the changes with the message from the prompt"
   
6. Amp: [Creates commit]
   
7. You: "Proceed to Prompt X.Y+1"
   
[Repeat]
```

## 🚨 Handling Issues

### If Tests Fail

```
STOP. Do not proceed with the next prompt.

Analyze the test failure:
1. Run: make acceptance-test 2>&1 | grep -A 10 "FAILED"
2. Show me the failing scenarios
3. Show the relevant git diff

Options:
a) Debug and fix the issue
b) Revert the last commit: git revert HEAD
c) Ask me for guidance

What do you recommend?
```

### If Amp Gets Confused

```
Let's reset context.

Current state:
- Phase: X
- Last completed prompt: X.Y
- Last commit: [hash]
- Test status: XXX/XXX passing

Re-read:
1. docs/refactoring/CHEAT_SHEET.md
2. The current prompt (Prompt X.Y+1)

Confirm you understand what to do next.
```

### If You Need to Pause

```
We're stopping for now. 

Please create a status report:
1. Last completed prompt: X.Y
2. Current test status: XXX/XXX passing
3. Files changed since last phase
4. Next step when we resume: Prompt X.Z

Save this to docs/refactoring/PROGRESS.md
```

## 📝 Communication Templates

### Starting a New Phase

```
Beginning Phase X: [Phase Name]

Re-read:
- docs/refactoring/AGENT_REFACTORING_PROMPTS.md (Phase X section)
- docs/refactoring/CHEAT_SHEET.md

Objectives for this phase:
- [List from the prompt]

Start with Prompt X.1. Report when ready to begin.
```

### Completing a Prompt

```
Prompt X.Y complete.

Changes:
- [List of files changed]

Tests: XXX/XXX passing ✅
Commit: [commit hash]

Ready for Prompt X.Y+1? (yes/no)
```

### Daily Summary

```
End of day summary:

Completed today:
- Prompts: X.A through X.Z
- Commits: N commits
- Tests: XXX/XXX passing

Tomorrow:
- Continue with Phase X
- Start with Prompt X.Z+1

Create this summary in docs/refactoring/daily/YYYY-MM-DD.md
```

## 🎯 Amp-Specific Tips

### If Amp Has File Viewing Limits

Some AI agents can only view a limited number of files. Prioritize:

1. Current prompt from AGENT_REFACTORING_PROMPTS.md
2. CHEAT_SHEET.md for reference
3. The specific file you're working on
4. Relevant test files

### If Amp Asks Questions

Be specific and reference the prompts:

```
Good: "Follow the pattern shown in Prompt 3.3 for FlowSquare"
Bad: "Do what you think is best"

Good: "Use the exact structure from CHEAT_SHEET.md section 'Domain Entity'"
Bad: "Create a domain entity however you want"
```

### If Amp Tries to Optimize

Keep Amp focused on the plan:

```
I appreciate the optimization idea, but let's stick to the plan for now.

We're following the structured refactoring in AGENT_REFACTORING_PROMPTS.md 
to ensure:
- Tests stay green
- Changes are incremental
- Rollback is easy if needed

Please execute Prompt X.Y as written, then we can discuss optimizations.
```

## 📊 Progress Tracking

Create a progress file:

```bash
# In docs/refactoring/
touch PROGRESS.md
```

Update it after each session:

```markdown
# Refactoring Progress

## Current Status
- **Phase:** 1 - Domain Layer
- **Last Prompt Completed:** 1.3
- **Next Prompt:** 1.4
- **Tests:** 45/45 passing ✅
- **Date:** 2024-01-15

## Completed Phases
- [x] Phase 0: Baseline (2024-01-14)

## In Progress
- [ ] Phase 1: Domain Layer
  - [x] 1.1: Create domain package structure
  - [x] 1.2: Extract Diagram entity
  - [x] 1.3: Extract Node entity
  - [ ] 1.4: Extract Edge and Position entities

## Notes
- FlowSquare extraction went smoothly
- Need to pay attention to Position dependencies
```

## 🎓 Learning from Each Phase

After completing each major phase, have Amp create a reflection:

```
Phase X complete. Please create a reflection document:

docs/refactoring/reflections/phase-X-reflection.md

Include:
- What went well
- Challenges encountered
- Lessons learned
- Suggestions for next phase

This helps us improve the process.
```

## ✅ Success Checklist

After each prompt:
- [ ] Steps executed as described
- [ ] Tests run: `make acceptance-test`
- [ ] All tests passing
- [ ] Changes committed
- [ ] Checklist updated

After each phase:
- [ ] All prompts in phase complete
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Progress file updated
- [ ] Ready for next phase

## 🚀 You're Ready!

Start with:
1. Run the setup script: `bash docs/refactoring/setup_refactoring_docs.sh`
2. Copy the 4 markdown files to `docs/refactoring/`
3. Start Amp in your project directory
4. Provide the initial context message from "Step 3" above
5. Begin with Phase 0 (DAY_1_KICKOFF.md)

Good luck! 🎯
