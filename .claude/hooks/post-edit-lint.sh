#!/bin/bash
# PostToolUse (Edit|Write|MultiEdit): record edited file path for batch
# processing in Stop hook. Must be <50ms - no formatting, no linting, no git.
cd "$CLAUDE_PROJECT_DIR" || exit 0

INPUT=$(cat)
FILE=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "default"')

[[ -z "$FILE" ]] && exit 0

# Skip generated/vendored paths
case "$FILE" in
  */vendor/*|*/mock_chartmogul/*|*/coverage/*) exit 0 ;;
esac

# Only track .go files
[[ "$FILE" != *.go ]] && exit 0

# Append to session-scoped tracker, respecting TMPDIR
TRACKER="${TMPDIR:-/tmp}/claude-edited-go-files-${CLAUDE_HOOK_SESSION_ID:-$SESSION_ID}"
echo "$FILE" >> "$TRACKER"
sort -u "$TRACKER" -o "$TRACKER"

exit 0
