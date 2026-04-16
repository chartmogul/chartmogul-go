#!/bin/bash
# SessionStart: idempotent workspace check (<1s, no network).
# Detects stale lockfiles, uncommitted changes, current branch.
cd "$CLAUDE_PROJECT_DIR" || exit 0

# Generate a stable session ID and persist via CLAUDE_ENV_FILE
# so the edit tracker and stop hook share the same file path
SESSION_ID="$(date +%s)-$$"
if [[ -n "$CLAUDE_ENV_FILE" ]]; then
  echo "export CLAUDE_HOOK_SESSION_ID='$SESSION_ID'" >> "$CLAUDE_ENV_FILE"
fi

ctx=""

# Check if go.sum is missing or stale relative to go.mod
if [[ -f go.mod ]]; then
  if [[ ! -f go.sum ]]; then
    ctx+="go.sum missing - run go mod tidy.\n"
  elif [[ go.mod -nt go.sum ]]; then
    ctx+="go.mod is newer than go.sum - run go mod tidy.\n"
  fi
fi

# Warn about uncommitted changes
dirty=$(git status --porcelain 2>/dev/null | head -5)
if [[ -n "$dirty" ]]; then
  ctx+="Uncommitted changes:\n$dirty\n"
fi

# Report current branch
branch=$(git branch --show-current 2>/dev/null)
if [[ -n "$branch" ]]; then
  ctx+="Branch: $branch\n"
fi

if [[ -n "$ctx" ]]; then
  jq -n --arg ctx "$ctx" '{
    "hookSpecificOutput": {
      "hookEventName": "SessionStart",
      "additionalContext": $ctx
    }
  }' || true
fi

exit 0
