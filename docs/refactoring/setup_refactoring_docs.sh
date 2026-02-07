#!/bin/bash
# Setup script for layli refactoring documents
# Usage: ./setup_refactoring_docs.sh

set -e

echo "🚀 Setting up layli refactoring documents..."

# Create directory structure
mkdir -p docs/refactoring
mkdir -p docs/architecture

echo "📁 Created directory structure"

# Check if we're in the layli project
if [ ! -f "go.mod" ] || ! grep -q "github.com/dnnrly/layli" go.mod 2>/dev/null; then
    echo "⚠️  Warning: This doesn't look like the layli project directory"
    echo "   Please run this script from the layli project root"
    read -p "   Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Create placeholder files (you'll need to copy the actual content)
cat > docs/refactoring/README.md << 'EOF'
# Layli Refactoring Documentation

This directory contains the complete refactoring plan for transforming layli
from a monolithic structure to clean architecture with BDD-first design.

## Quick Start

1. Start with `DAY_1_KICKOFF.md` - Copy this to your AI agent (Amp/Windsurf/Cursor)
2. Keep `CHEAT_SHEET.md` open as a reference
3. Follow `AGENT_REFACTORING_PROMPTS.md` phase by phase
4. Track progress with `QUICK_REFERENCE_CHECKLIST.md`

## Files

- **DAY_1_KICKOFF.md** - Start here! Phase 0 prompts to establish baseline
- **AGENT_REFACTORING_PROMPTS.md** - Complete phase-by-phase prompts (main playbook)
- **QUICK_REFERENCE_CHECKLIST.md** - Condensed checklist for tracking progress
- **CHEAT_SHEET.md** - One-page quick reference (keep visible while working)

## Usage with Amp

See `AMP_USAGE_GUIDE.md` for specific instructions on using these prompts with Amp.

## Progress Tracking

As you complete phases, update the checklist and create notes:
- `docs/CURRENT_STATE.md` - Created in Phase 0
- `docs/FEATURE_MAP.md` - Created in Phase 0
- `docs/REFACTORING_COMPLETE.md` - Created in Phase 7
EOF

echo "✅ Created README.md"

# Create a .gitignore for temporary refactoring files
cat > docs/refactoring/.gitignore << 'EOF'
# Temporary files during refactoring
*.tmp
*_backup.md
.DS_Store
EOF

echo "✅ Created .gitignore"

# Create instructions for copying the actual prompt files
cat > docs/refactoring/COPY_FILES_HERE.txt << 'EOF'
📥 Copy These Files to This Directory:

You should have downloaded 4 markdown files from Claude:

1. DAY_1_KICKOFF.md
2. AGENT_REFACTORING_PROMPTS.md
3. QUICK_REFERENCE_CHECKLIST.md
4. CHEAT_SHEET.md

Copy them to this directory (docs/refactoring/), then delete this file.

The files should be at the same level as this file:

docs/refactoring/
  ├── DAY_1_KICKOFF.md
  ├── AGENT_REFACTORING_PROMPTS.md
  ├── QUICK_REFERENCE_CHECKLIST.md
  ├── CHEAT_SHEET.md
  ├── README.md
  └── COPY_FILES_HERE.txt (delete me after copying files)
EOF

echo "✅ Created file placement instructions"

# Create a commit template for refactoring commits
cat > docs/refactoring/commit_template.txt << 'EOF'
# Refactoring Commit Template
# 
# Usage: git commit -t docs/refactoring/commit_template.txt
#
# Format:
# <type>(<scope>): <subject>
#
# <optional body>
#
# Tests: <test results>

# --- COMMIT MESSAGE STARTS BELOW THIS LINE ---


# Tests: 

# --- Template Guide ---
# Type: refactor, test, docs, feat, fix
# Scope: domain, usecases, adapters, test, cli, config
#
# Example:
# refactor(domain): extract Diagram entity
#
# Move Diagram struct from config.go to internal/domain/diagram.go.
# Add validation methods and comprehensive documentation.
#
# Tests: All acceptance tests passing (45/45)
EOF

echo "✅ Created commit template"

echo ""
echo "✨ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Copy the 4 markdown files to docs/refactoring/"
echo "2. Read docs/refactoring/README.md"
echo "3. Start with docs/refactoring/DAY_1_KICKOFF.md"
echo ""
echo "Optional: Set git commit template:"
echo "  git config commit.template docs/refactoring/commit_template.txt"
