# Repository Guidelines

## Project Structure & Module Organization
- `cmd/wc` holds the CLI entrypoint; keep binaries thin and delegate logic.
- `internal/cli` wires CLI parsing to the pipeline engine; extend flags here.
- `internal/pipeline` contains `source`, `tokenizer`, `counters`, and `aggregator` stages; add new stages behind focused packages.
- `docs/DESIGN.md` captures the architecture and pipeline flow; update diagrams and invariants when modifying stages.
- `internal/accepttest` hosts txtar-driven end-to-end tests, with fixtures under `internal/accepttest/testdata`.
- `testdata` stores shared sample inputs; document unusual encodings in README stubs within that folder.

## Build, Test, and Development Commands
- `go build ./cmd/wc` builds the CLI binary; prefer this over ad-hoc `go install` during development.
- `go run ./cmd/wc --help` is the quickest way to verify argument handling and default output.
- `go test ./...` runs unit and integration tests; add `-run` filters when iterating on a single package.
- `go test ./internal/accepttest` executes the txtar acceptance suite; append `-v` for detailed diff output.

## Coding Style & Naming Conventions
- Format Go code with `gofmt` (automatic via `go fmt ./...`); no custom style exceptions.
- Use Go module-relative import paths (`wc/internal/...`).
- Exported identifiers should read like phrases (`WordCounter`), while unexported ones stay concise (`countTokens`).
- Directory names favor nouns for packages (`tokenizer`, `aggregator`); avoid camel case on disk.

## Testing Guidelines
- Unit tests live alongside code with `_test.go` suffixes and table-driven cases where meaningful.
- Name tests `Test<Thing>` or `Test_<Scenario>`; acceptance archives mirror scenario names (`basic.txtar`).
- Aim to cover new branches with either Go tests or txtar cases; validate edge cases such as empty input and large files.
- Snapshot outputs belong in txtar fixtures; keep them minimal and comment directives at the top.

## Commit & Pull Request Guidelines
- Follow the existing imperative, lowercase history (e.g., `Implement initial wc pipeline`, `Now the AT can find the test cases`).
- Reference linked issues in the body when applicable and note behavioural changes.
- Include repro steps or sample commands for CLI changes; attach before/after output snippets when useful.
- Verify `go test ./...` before opening a PR and mention acceptance coverage in the description.
