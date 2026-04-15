#!/bin/bash
if [ -x ./bin/golangci-lint ]; then
  ./bin/golangci-lint run --fix ./... 2>&1
fi
go test -v -timeout=10m ./... 2>&1
go build ./... 2>&1

exit 0
