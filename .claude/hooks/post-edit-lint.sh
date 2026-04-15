#!/bin/bash
INPUT=$(cat)
FILE=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

if [[ "$FILE" == *.go ]]; then
  gofmt -w "$FILE" 2>/dev/null
  if [ -x ./bin/golangci-lint ]; then
    ./bin/golangci-lint run --fix "$FILE" 2>&1
  fi
fi

exit 0
