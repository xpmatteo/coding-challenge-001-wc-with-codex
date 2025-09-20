# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of the Unix `wc` (word count) command. The project follows a clean pipeline architecture where CLI arguments are transformed through stages: ParseArgs -> AnalyzeFiles -> AddTotal -> Format -> output.

## Development Commands

### Testing
- `make test` - Run all tests with custom GOCACHE location
- `go test ./...` - Standard test runner (without custom cache)
- `go test -run TestSpecificTest` - Run a specific test
- `go test -v ./internal/accepttest` - Run acceptance tests with verbose output

### Linting
- `make lint` - Run golangci-lint with custom cache location

### Building
- `go build ./cmd/wc` - Build the wc binary
- `go run ./cmd/wc [args]` - Run without building

## Architecture

### Core Pipeline (internal/wc/pipeline.go)
The main processing pipeline consists of these stages:
1. **ParseArgs**: Converts CLI arguments into a Config struct
2. **AnalyzeFiles**: Processes each file and returns Stats
3. **AddTotal**: Adds total line when multiple files are processed
4. **Format**: Converts Stats into formatted output strings

### Key Types
- **Config**: Captures CLI flags (CountBytes, CountLines, CountWords) and file list
- **Stats**: Results for a single file (Name, Lines, Words, Bytes, Chars)
- **App** (internal/cli/runner.go): Main CLI runner that orchestrates the pipeline

### Entry Point
- `cmd/wc/main.go`: Simple main function that delegates to cli.App
- `internal/cli/runner.go`: Contains the Run method that executes the full pipeline

## Testing Strategy

### Acceptance Tests
The project uses txtar-based acceptance tests in `internal/accepttest/`. Test cases are stored as `.txtar` files in `testdata/` with YAML directives in comments for configuration:

```yaml
# args: ["-l", "file.txt"]
# stdin: input.txt
# env: ["VAR=value"]
```

Each test case can specify expected stdout, stderr, and exit code.

### Unit Tests
- Individual pipeline stages have dedicated unit tests
- Uses stretchr/testify for assertions (follows project conventions)
- Tabular test format preferred

## File Organization

- `cmd/wc/` - Main entry point and CLI tests
- `internal/wc/` - Core pipeline logic and unit tests
- `internal/cli/` - CLI application runner
- `internal/accepttest/` - End-to-end acceptance tests
- `testdata/` - Test fixtures for acceptance tests
- `docs/` - Design documentation and development notes

## Special Notes

- The project uses custom GOCACHE and GOLANGCI_LINT_CACHE locations (`.gocache/` and `.golangci/`)
- All pipeline functions are stateless and testable in isolation
- Error handling follows the pattern of continuing processing while accumulating errors
- The codebase follows the design outlined in `docs/DESIGN.md`