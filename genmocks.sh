#!/bin/bash
mkdir -p mock_chartmogul ; mockgen github.com/chartmogul/chartmogul-go IApi > mock_chartmogul/chartmogul.go
