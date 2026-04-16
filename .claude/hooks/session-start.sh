#!/bin/bash
# SessionStart: idempotent workspace check (<1s, no network).
# Detects stale lockfiles, uncommitted changes, current branch.
cd "$CLAUDE_PROJECT_DIR" || exit 0

ctx=""

# Check if go.mod is newer than go.sum (deps out of date)
if [[ -f go.mod && -f go.sum ]]; then
  if [[ go.mod -nt go.sum ]]; then
    ctx+="go.mod is newer than go.sum - run go mod tidy\n"
  fi
fi

# Check for uncommitted changes
dirty=$(git diff --name-only 2>/dev/null | head -5)
if [[ -n "$dirty" ]]; then
  ctx+="Uncommitted changes:\n$dirty\n"
fi

# Report current branch
branch=$(git branch --show-current 2>/dev/null)
if [[ -n "$branch" ]]; then
  ctx+="Branch: $branch\n"
fi

# Clear stale file tracker from previous session
rm -f "$CLAUDE_PROJECT_DIR/.claude/.edited_files" 2>/dev/null

if [[ -n "$ctx" ]]; then
  jq -n --arg ctx "$ctx" \
    '{"hookSpecificOutput": {"hookEventName": "SessionStart", "additionalContext": $ctx}}' || true
fi

exit 0
