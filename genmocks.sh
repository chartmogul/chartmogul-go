#!/bin/bash
MOCKGEN_VERSION=v0.6.0
mkdir -p mock_chartmogul ; go run go.uber.org/mock/mockgen@"${MOCKGEN_VERSION}" github.com/chartmogul/chartmogul-go/v5 IApi > mock_chartmogul/chartmogul.go
