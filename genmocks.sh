#!/bin/bash
mkdir -p mock_chartmogul ; mockgen github.com/chartmogul/chartmogul-go/v2 IApi > mock_chartmogul/chartmogul.go
