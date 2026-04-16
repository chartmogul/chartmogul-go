#!/bin/bash
# PostToolUse (Bash, scoped to git commit): lint files in the just-created
# commit and report failures as context. Must not run tests.
cd "$CLAUDE_PROJECT_DIR" || exit 0

INPUT=$(cat)
CMD=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

# Only fire after git commit commands
case "$CMD" in
  git\ commit*) ;;
  *) exit 0 ;;
esac

# Get .go files changed in the last commit
FILES=$(git diff --name-only HEAD~1 HEAD 2>/dev/null | grep '\.go$' || true)
[[ -z "$FILES" ]] && exit 0

# Run golangci-lint on committed files only
if [ -x ./bin/golangci-lint ]; then
  LINT_OUTPUT=$(./bin/golangci-lint run $FILES 2>&1)
  LINT_EXIT=$?
  if [ $LINT_EXIT -ne 0 ]; then
    CTX=$(echo "$LINT_OUTPUT" | head -20)
    jq -n --arg ctx "golangci-lint found issues in committed files:\n$CTX" \
      '{"hookSpecificOutput": {"hookEventName": "PostToolUse", "additionalContext": $ctx}}' || true
  fi
fi

exit 0
