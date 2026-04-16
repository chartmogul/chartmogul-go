#!/bin/bash
# PostToolUse (Edit|Write|MultiEdit): record edited file path for batch
# processing in Stop hook. Must be <50ms - no formatting, no linting, no git.
cd "$CLAUDE_PROJECT_DIR" || exit 0

INPUT=$(cat)
FILE=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

[[ -z "$FILE" ]] && exit 0

# Skip generated/vendored paths
case "$FILE" in
  */vendor/*|*/mock_chartmogul/*|*/coverage/*|*/dist/*|*/__pycache__/*) exit 0 ;;
esac

# Only track .go files
[[ "$FILE" != *.go ]] && exit 0

# Append to tracker file (dedup happens in Stop hook)
TRACKER="$CLAUDE_PROJECT_DIR/.claude/.edited_files"
echo "$FILE" >> "$TRACKER" 2>/dev/null || true

exit 0
