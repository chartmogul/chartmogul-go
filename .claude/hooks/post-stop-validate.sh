#!/bin/bash
# Stop: batch-format files touched this turn, run lint + tests, report failures.
# Tests run unconditionally (~2s) since edits to test or lib code both matter.
cd "$CLAUDE_PROJECT_DIR" || exit 0

TRACKER="/tmp/claude-edited-go-files-${CLAUDE_HOOK_SESSION_ID:-default}"

ctx=""

# Batch gofmt + golangci-lint on tracked files (if any were edited)
if [[ -f "$TRACKER" ]]; then
  files=$(cat "$TRACKER")
  rm -f "$TRACKER"

  if [[ -n "$files" ]]; then
    echo "$files" | xargs gofmt -w 2>/dev/null || true

    if [ -x ./bin/golangci-lint ]; then
      lint_out=$(echo "$files" | xargs ./bin/golangci-lint run --fix 2>&1)
      lint_exit=$?
      if [ $lint_exit -ne 0 ]; then
        offenses=$(echo "$lint_out" | head -20)
        ctx+="golangci-lint offenses remaining after autocorrect:\n$offenses\n"
      fi
    fi
  fi
fi

# Always run tests - edits to test or lib code both matter (~2s)
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
