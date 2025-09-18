# wc Clone Design

This document outlines the planned architecture for the wc implementation. The goal is to provide a small, testable pipeline that mirrors the behaviour of the Unix `wc` command while remaining easy to extend.

## Overview

The CLI entry point (`internal/cli.App`) delegates to a pipeline that transforms command‑line arguments into formatted output lines. Each stage accepts simple inputs, returns data plus errors, and avoids hidden state so that we can unit test stages in isolation.

```
args -> ParseArgs -> AnalyzeFiles -> AddTotal -> Format -> []string
```

`AnalyzeFiles` calls `AnalyzeFile` for each input so we can reuse the counting logic across files and stdin.

## Data Model

- `Config`: captures the requested counters (lines, words, bytes, optional future counters) and the ordered list of file paths. It also records whether stdin should be read (no files or explicit single `-`). Flags default to counting lines, words, and bytes—matching GNU/BSD `wc` when no switches are provided. `ParseArgs` rejects multiple requests to read stdin so we never consume it more than once.
- `Stats`: one record per analyzed input. Fields include `Name string`, `Lines int`, `Words int`, `Bytes int`, and placeholders for future counters such as `Chars int` once we decide how to handle rune counting. A lightweight mask or enum indicates which counters were requested so that formatting can stay generic.

## Pipeline Stages

### ParseArgs(args []string) (Config, error)

- Validates and normalizes CLI flags.
- Distinguishes between file names, `-` (stdin), and no files (stdin implied) while ensuring stdin is requested at most once.
- Determines the active counters: if no counter flags are supplied we enable the default trio; otherwise we honour only the requested ones, keeping the traditional `wc` column order.
- Leaves I/O to later stages; errors cover malformed flag combinations or missing operands.

### AnalyzeFiles(cfg Config) ([]Stats, error)

- Iterates over the resolved file targets and delegates each to `AnalyzeFile`.
- Handles stdin by passing the provided reader when `Name == "-"` or when no paths were supplied.
- Mirrors `wc` behaviour by continuing after individual file errors while recording the failure details so the caller can emit diagnostics and set a non‑zero exit code.

### AnalyzeFile(source io.Reader, name string, counters CounterMask) (Stats, error)

- Streams the content once, incrementing only the requested counters.
- Uses buffered reads to keep memory constant, treating `\n` as the line delimiter.
- Word counting follows the POSIX definition (transitions between whitespace and non‑whitespace).
- Supports byte counts today. Rune counting (e.g. UTF‑8 decoded character totals) will be added once we settle on the exact flag semantics.

### AddTotal(cfg Config, stats []Stats) ([]Stats, error)

- If more than one file (or stdin plus files) contributed, appends a `Stats` entry named `"total"` whose counters are the per‑field sums.
- Leaves the slice unchanged when only a single input exists.

### Format(cfg Config, stats []Stats) ([]string, error)

- Converts `Stats` into text lines matching `wc` spacing.
- Computes column widths by inspecting only the selected counters to ensure alignment across rows.
- Emits entries in the original order followed by `total` when present, retaining `wc`’s default column ordering (lines, words, bytes, …).

## Error Handling

- User input errors (unknown flags, missing filenames, repeated stdin requests) surface from `ParseArgs` and cause `Run` to emit messages on stderr with a non‑zero exit code.
- I/O errors are tagged with the offending file name. We continue processing subsequent files, accumulate whether any error occurred, and reflect that in the exit status—matching the stock `wc` behaviour.

## Testing Strategy

- Unit tests for each stage: flags parsing, counter logic, formatting alignment.
- Property‑style fuzz tests for `AnalyzeFile` can guard against panics for arbitrary byte streams.
- Acceptance tests (already scaffolded under `internal/accepttest`) exercise the end‑to‑end CLI, comparing stdout, stderr, and exit codes against fixtures.

## Extensibility

- New counters only require updating `Config`, `Stats`, `AnalyzeFile`, and `Format` to include the extra column.
- Alternative front ends (e.g. library API) can reuse `AnalyzeFile` directly for ad‑hoc counting.
- Parallel file processing can be added later by introducing a worker pool around `AnalyzeFile` without touching downstream stages.
