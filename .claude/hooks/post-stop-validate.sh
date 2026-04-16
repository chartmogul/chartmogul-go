#!/bin/bash
# Stop: batch-format files touched this turn, run lint + tests, report failures.
# Tests and lint only run when files were tracked this turn.
cd "$CLAUDE_PROJECT_DIR" || exit 0

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "default"')
TRACKER="${TMPDIR:-/tmp}/claude-edited-go-files-${CLAUDE_HOOK_SESSION_ID:-$SESSION_ID}"

# Nothing tracked this turn - skip everything
[[ ! -f "$TRACKER" ]] && exit 0

files=$(cat "$TRACKER")
rm -f "$TRACKER"

[[ -z "$files" ]] && exit 0

ctx=""

# 1. Autoformat touched files
echo "$files" | xargs gofmt -w 2>/dev/null || true

# 2. Lint touched files and report remaining offenses
if [ -x ./bin/golangci-lint ]; then
  lint_out=$(echo "$files" | xargs ./bin/golangci-lint run --fix 2>&1)
  lint_exit=$?
  if [ $lint_exit -ne 0 ]; then
    offenses=$(echo "$lint_out" | head -20)
    ctx+="golangci-lint offenses remaining after autocorrect:\n$offenses\n"
  fi
fi

# 3. Run test suite (~2s)
test_out=$(go test -timeout=30s ./... 2>&1)
test_exit=$?
if [[ $test_exit -ne 0 ]]; then
  summary=$(echo "$test_out" | grep -E 'FAIL|Error|panic' | head -20)
  ctx+="go test failed:\n$summary\n"
fi

if [[ -n "$ctx" ]]; then
  jq -n --arg ctx "$ctx" '{
    "hookSpecificOutput": {
      "hookEventName": "Stop",
      "additionalContext": $ctx
    }
  }' || true
fi

exit 0
