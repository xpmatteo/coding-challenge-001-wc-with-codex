
.PHONY: test
test:
	GOCACHE=$$(pwd)/.gocache go test ./...

.PHONY: lint
lint: 
	GOLANGCI_LINT_CACHE=$$(pwd)/.golangci golangci-lint run
