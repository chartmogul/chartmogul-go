#!/bin/bash
# Stop: batch-format files touched this turn, run lint + tests, report failures.
# Suite takes ~2s so it's cheap enough to run every turn.
cd "$CLAUDE_PROJECT_DIR" || exit 0

TRACKER="$CLAUDE_PROJECT_DIR/.claude/.edited_files"
ERRORS=""

# Nothing tracked this turn - skip everything
[[ ! -f "$TRACKER" ]] && exit 0

# Dedup tracked files, filter to existing .go files
FILES=$(sort -u "$TRACKER" | while read -r f; do [[ -f "$f" ]] && echo "$f"; done)
rm -f "$TRACKER" 2>/dev/null

[[ -z "$FILES" ]] && exit 0

# 1. Autoformat touched files
echo "$FILES" | xargs gofmt -w 2>/dev/null || true

# 2. Lint touched files and report remaining offenses
if [ -x ./bin/golangci-lint ]; then
  LINT_OUTPUT=$(echo "$FILES" | xargs ./bin/golangci-lint run --fix 2>&1)
  if [ $? -ne 0 ]; then
    ERRORS="$ERRORS\n## golangci-lint\n$(echo "$LINT_OUTPUT" | head -20)"
  fi
fi

# 3. Run full test suite (~2s)
TEST_OUTPUT=$(go test -timeout=30s ./... 2>&1)
if [ $? -ne 0 ]; then
  ERRORS="$ERRORS\n## go test\n$(echo "$TEST_OUTPUT" | grep -E 'FAIL|Error|panic' | head -20)"
fi

if [[ -n "$ERRORS" ]]; then
  jq -n --arg ctx "$(printf "Issues found after this turn:$ERRORS")" \
    '{"hookSpecificOutput": {"hookEventName": "Stop", "additionalContext": $ctx}}' || true
fi

exit 0
