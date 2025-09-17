# Design Overview

## Goals
- Implement a streaming clone of GNU `wc` in Go with predictable performance on large files.
- Maintain a modular pipeline so alternative sources/tokenizers/counters can be swapped in later.
- Provide clear separation between counting logic, formatting, and CLI wiring to ease testing and future extensions.

## High-Level Architecture
```
Sources -> Tokenizer -> Counters -> Aggregator -> Formatter -> CLI
```

The pipeline converts raw inputs into structured results while remaining composable:
- **Sources** supply `io.Reader` streams backed by files or standard input.
- **Tokenizer** reads buffered bytes and emits semantic tokens (runes, newlines, EOF, decode errors).
- **Counters** consume tokens to track lines, words, bytes, runes, and longest line length.
- **Aggregator** runs the pipeline per file, collates errors, and accumulates totals.
- **Formatter** produces GNU-style aligned output respecting the requested columns.
- **CLI** parses flags, orchestrates the pipeline for each input, and handles exit codes.

## Package Layout
```
cmd/
  wc/           # main package wiring CLI to the engine
internal/
  pipeline/
    source/     # Source implementations
    tokenizer/  # Token emission from readers
    counters/   # Counter state machines
    aggregator/ # Engine coordinating stages
  format/       # Output table rendering
  cli/          # Flag parsing helpers, column selection, totals logic
docs/
  acceptance-tests.md # txtar acceptance suite guide
```

## Pipeline Contracts
- `source.Source` exposes `Open(context.Context) (io.ReadCloser, FileMeta, error)` where `FileMeta` retains the display name and size hints.
- `tokenizer.Tokenizer` delivers `Token{Kind, Rune, Bytes, IsWhitespace, Err}` via `Next() (Token, error)`. It hides buffering details and tracks byte/rune widths for multibyte characters.
- `counters.Counters` owns fields `Bytes`, `Runes`, `Lines`, `Words`, `MaxLineLen` plus internal state (`inWord`, `lineRuneLen`). Method `Consume(Token)` updates counts; Rune errors increment bytes and treat the replacement rune as a single character to match GNU semantics.
- `aggregator.Engine` binds a `Source`, `TokenizerFactory`, and `CounterFactory`, producing per-file `Result{File string, Counts counters.Snapshot, Err error}` along with a total row when more than one file is processed.

## Flag Handling & Column Selection
- Supported flags: `-c` bytes, `-m` runes, `-w` words, `-l` lines, `-L` longest line.
- When no column flag is given, default order mirrors GNU (`lines words bytes`).
- `cli.Columns` captures boolean selections and can render an ordered slice for formatting.
- The CLI reports errors to `stderr`, continues processing remaining files, and exits non-zero if any input fails.

## Testing Strategy
1. **Unit tests** for tokenizer and counters using fixtures that cover:
   - Empty input, multiple whitespace patterns, newline handling, very long lines.
   - Valid and invalid UTF-8 sequences, including multibyte graphemes.
   - Windows-style line endings and binary blobs.
2. **Integration tests** for the CLI using in-memory sources (single file, multiple files, stdin).
3. **Acceptance tests** via txtar archives to assert end-to-end behaviour and formatting consistency.

## Acceptance Tests with txtar
- Each archive resides under `testdata/acceptance/`. The comment header can include directives like `args: -w file.txt` and `stdin: input.txt`.
- Files inside the archive represent inputs (`file.txt`), expected `stdout`, and expected `stderr`.
- The acceptance harness materializes archives into a temporary directory, runs the built binary with the specified args, and compares outputs/exit codes to golden expectations (allowing placeholder substitution for temp paths).

## Future Extensions
- Swap tokenizer for memory-mapped reader on supported platforms.
- Add JSON output by bolting an alternate formatter.
- Parallelize per-file processing while preserving consistent totals.
- Extend acceptance harness to compare against system `wc` when available.
