.DEFAULT_GOAL := test

export GO111MODULE=on

.PHONY: dependencies
dependencies:
	# Get dependencies
	go get ./... || echo "Some dependencies failed"

.PHONY: mocks
mocks: dependencies
	# Build mocks
	./genmocks.sh

.PHONY: test
test: lint
	go test -v -timeout=10m ./...

.git/hooks/pre-commit:
	cp pre-commit .git/hooks/pre-commit
	chmod a+x .git/hooks/pre-commit

.PHONY: pre-commit
pre-commit: .git/hooks/pre-commit lint
	go test -timeout=20s --short ./...

bin/golangci-lint:
	wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.15.0

.PHONY: lint
lint: bin/golangci-lint
	./bin/golangci-lint run
