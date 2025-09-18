# Repository Guidelines

## Project Structure & Module Organization
- CLI entry point lives in `cmd/wc`, with `main.go` invoking the internal application.
- Shared pipeline logic sits under `internal/wc` (`ParseArgs`, `AnalyzeFiles`, `AddTotal`, `Format`).
- Integration harness lives at `internal/cli` and acceptance fixtures at `internal/accepttest/testdata` (txtar format).
- Documentation resides in `docs/`, including the challenge brief and step-by-step prompts.

## Build, Test, and Development Commands
- `make test` – runs the full Go test suite with the repo-scoped build cache.
- `make lint` – executes `golangci-lint run` using the local lint cache.
- `GOCACHE=$PWD/.gocache go run ./cmd/wc <file>` – exercise the walking skeleton manually when needed.
- `rm -rf .gocache .golangci` – clean the temporary build/lint caches after runs.

## Coding Style & Naming Conventions
- Standard Go formatting via `gofmt`; run `gofmt -w <files>` after edits.
- Prefer short, intention-revealing identifiers; exported symbols only when shared across packages.
- Keep comments sparse and explanatory—only for non-obvious behaviour.

## Testing Guidelines
- Unit tests live alongside package code (`internal/wc/*_test.go`) with one concern per file.
- Acceptance tests rely on txtar fixtures (`internal/accepttest/testdata`) and `TestAcceptanceSuite`.
- Ensure new features include both unit coverage and an acceptance fixture (unskip or add) that fails before implementation.

## Commit & Pull Request Guidelines
- Follow concise imperative commit messages (e.g., “Implement byte counting for AnalyzeFile”).
- PRs should summarize intent, link relevant issues/challenge steps, and note testing (`go test`, manual checks). Include snapshots of CLI output when behaviour changes.

## Agent-Specific Notes
- Respect existing skips in acceptance fixtures; only unskip when implementing that capability.
- Use the staged prompts in `docs/STEPS.md` to advance one requirement at a time, waiting for review between steps.
