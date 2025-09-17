# Acceptance Test Suite

The `testdata/acceptance` directory stores end-to-end scenarios using [`txtar`](https://pkg.go.dev/golang.org/x/tools/txtar). Each archive captures command-line arguments, input files, and expected outputs for a single invocation of `wc`.

## Archive Structure
- **Comment header**: include directives in comment lines at the top of the archive.
  - `args: -w -l foo.txt` — space-separated CLI arguments.
  - `stdin: input.txt` — optional filename whose content will be piped to stdin.
  - `env: KEY=VALUE` — optional environment variables applied to the command.
- **Files**: plain files within the archive. Suggested names:
  - `foo.txt`, `bar.bin` — inputs referenced by `args` or `stdin`.
  - `stdout.txt` — expected standard output.
  - `stderr.txt` — expected standard error (omit if empty).
  - `exitcode` — optional single-line file with the expected numeric exit status (defaults to `0`).

## Test Harness Behaviour
1. Parse the archive with `txtar.ParseFile`.
2. Materialize each file in a temporary working directory.
3. Build (or reuse) the `wc` binary under test.
4. Execute with provided arguments and environment, supplying stdin when requested.
5. Compare stdout, stderr, and exit code against the golden files.
6. Support placeholder substitution (e.g., `%TMPDIR%`) in expected outputs to account for run-specific paths.

## Conventions
- Keep inputs small and focused; large binary blobs can be base64- or hex-encoded within the archive if necessary.
- Prefer UTF-8 text for readability, noting when invalid sequences are intentional.
- Name archives descriptively (`basic.txtar`, `unicode.txtar`, `binary.txtar`).
- Document unusual scenarios in the comment header for quick reference.

## Running the Suite
The acceptance tests will be integrated into `go test ./...` by importing the harness under an `//go:build acceptance` tag or a dedicated test file. They can also be run explicitly via a helper script once the binary exists.

