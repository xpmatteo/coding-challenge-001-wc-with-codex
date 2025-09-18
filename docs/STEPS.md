# Implementation Prompts

Each section below translates the Coding Challenges “Build Your Own wc” steps into an incremental, test-first workflow that fits our pipeline design (`ParseArgs → AnalyzeFile/AnalyzeFiles → AddTotal → Format`). At every stage we grow the walking skeleton just enough to satisfy the new requirement, keep other behaviour stubbed, and finish with a green test suite plus a manual smoke check.

## Step Zero – Walking Skeleton

Prompt:
1. Keep every acceptance fixture skipped and add a new one, e.g. `internal/accepttest/testdata/step0-walking-skeleton.txtar`, with `# args: sample.txt` and an expected `stdout.txt` of `0 0 0 sample.txt`; leave the embedded file content tiny (two or three words). Set `# skipped: false` for this fixture only.
2. Run `GOCACHE=$PWD/.gocache go test ./internal/accepttest -run Step0` (or just the default suite) and confirm the failure is because the CLI prints nothing yet rather than a crash or compile error; capture that message so you recognise it later.
3. Introduce unit tests for every pipeline stage that assert the stub behaviour: `ParseArgs` returns a config with the provided file name and no counters selected, `AnalyzeFile`/`AnalyzeFiles` return a `Stats` struct populated with the file name and zero counts, `AddTotal` is a no-op for a single entry, and `Format` renders `0 0 0 sample.txt`. Run `go test ./...` and make sure those tests fail for the logical reason (missing functions or wrong values), not because the package fails to compile.
4. Implement the bare-bones pipeline: add the `wc` package with the four functions returning zeros, wire `cli.App.Run` to call them, and make sure `Run` writes the formatted lines to `stdout` while still leaving real counting unimplemented.
5. Re-run `go test ./...` (include the acceptance test) and verify everything passes. `go build ./cmd/wc` followed by `go run ./cmd/wc sample.txt` should now print `0 0 0 sample.txt`.
6. Stop here and wait for user feedback before progressing.

## Step One – Bytes (-c)

Prompt:
1. Edit `internal/accepttest/testdata/bytes-chars.txtar` so it only exercises byte counting for now (`# args: -c sample.txt`), update `stdout.txt` to the expected byte total (match it with `wc -c` on the embedded file), and flip `# skipped: false`. Keep the Step Zero fixture in place so the default output remains stubbed.
2. Run the acceptance suite and confirm this fixture fails because the output still reads `0` bytes. Note the diff so you know the failure reason is the missing byte counter.
3. Extend unit tests for the pipeline: `ParseArgs` should now recognise `-c`/`--bytes` and mark the bytes counter; `AnalyzeFile` should read from an `io.Reader` and report the byte length; `AnalyzeFiles` should surface filesystem errors cleanly; `AddTotal` remains a no-op; `Format` should render just the byte column (e.g. `      12 sample.txt`). Run the unit tests to see them fail on the logical assertions.
4. Implement the minimum logic to satisfy those tests: real flag parsing for bytes, a streaming byte counter, and formatting that honours the active counters without touching the yet-to-be-implemented ones. Leave default (no flag) behaviour returning zeros for now.
5. Re-run unit and acceptance tests; everything should be green. Build and smoke test manually with `go run ./cmd/wc -c sample.txt` to confirm the byte total matches `wc -c`.
6. Wait for user feedback before taking on the next requirement.

## Step Two – Lines (-l)

Prompt:
1. Modify `internal/accepttest/testdata/flags-lines-words.txtar` to cover only the `-l` flag for this step, update `stdout.txt` to the expected line count for the embedded file, and set `# skipped: false`. Keep the word column commented out for now.
2. Run the acceptance suite; the failure should show that the line count is still zero, confirming the missing implementation rather than a parsing error.
3. Update unit tests so: `ParseArgs` recognises `-l`/`--lines`, `AnalyzeFile` increments the line counter on `
`, `AnalyzeFiles` continues to pass through stats, `AddTotal` still no-op, and `Format` renders only the line column when that is the sole counter. Run the tests and ensure they fail because counts stay at zero.
4. Implement the smallest change set to make those tests pass: add line-counting during analysis and allow formatting to output lines-only while preserving byte support from Step One.
5. Re-run unit + acceptance tests, then manually verify with `go run ./cmd/wc -l sample.txt`. Everything should pass while the default (no flag) path remains the stub from Step Zero.
6. Pause for user feedback.

## Step Three – Words (-w)

Prompt:
1. Extend `internal/accepttest/testdata/flags-lines-words.txtar` to request both `-l` and `-w`, update `stdout.txt` to the expected two-column output, and keep it unskipped.
2. Run the acceptance suite and make sure the failure is that the word column remains zero—line counting should already succeed.
3. Add/adjust unit tests so that: `ParseArgs` handles `-w`/`--words` and preserves flag order, `AnalyzeFile` counts word boundaries according to whitespace transitions, `AnalyzeFiles` aggregates both counts, `AddTotal` still only mirrors the input for a single file, and `Format` prints the selected columns in traditional `wc` order (lines, words, bytes) for this subset. Run the tests to watch them fail for the missing word totals.
4. Implement word counting in the analyzer and update formatting to include the new column while keeping previous features untouched.
5. Re-run the full test suite and manually check `go run ./cmd/wc -w sample.txt` as well as `-lw` to confirm correctness.
6. Wait for user feedback.

## Step Four – Characters (-m)

Prompt:
1. Restore `internal/accepttest/testdata/bytes-chars.txtar` so it now exercises both `--bytes` and `--chars` (or `-c -m`). Ensure the embedded file contains multi-byte characters, update the expected output to match `wc`, and keep it unskipped. Also unskip `internal/accepttest/testdata/invalid-utf8.txtar` so we verify how the tool behaves on malformed sequences.
2. Run the acceptance tests and confirm the failures come from the character column being wrong (or from panics in the invalid UTF-8 case), not from parsing issues.
3. Expand unit tests: `ParseArgs` should recognise `-m`/`--chars` aliases, `AnalyzeFile` must decode runes (while still reporting raw bytes), include coverage for invalid UTF-8 handling, `AnalyzeFiles` should surface both counters, `AddTotal` remains untouched, and `Format` needs to render whichever subset of columns is requested while maintaining alignment. Run the tests and make sure they fail because the char counts are still zero or mismatched.
4. Implement rune counting alongside byte counting, handle invalid UTF-8 consistently with `wc`, and adjust formatting to keep columns aligned regardless of the active counters.
5. Re-run all tests (unit and acceptance) and manually compare `go run ./cmd/wc -m sample.txt` with the system `wc -m sample.txt` to ensure parity.
6. Await user feedback.

## Step Five – Default Output (no flags)

Prompt:
1. Unskip `internal/accepttest/testdata/basic.txtar` so it checks the default invocation (no flags) and embeds a sample file; ensure the expected output contains lines, words, and bytes in that order. Also unskip `internal/accepttest/testdata/multiple-files.txtar` so totals are exercised.
2. Run the acceptance suite and verify the failures stem from the default path still emitting `0 0 0` and from totals missing, not from unrelated regressions.
3. Extend unit tests: `ParseArgs` should enable the lines/words/bytes trio when no specific counter flags are provided, `AnalyzeFiles` must iterate over multiple targets, `AddTotal` should now add a `total` row when more than one file (including stdin plus files) is processed, and `Format` must align multi-column output and append totals correctly. Add unit coverage for column width calculation and total aggregation. Run the tests to ensure they fail because totals/default behaviour isn’t implemented yet.
4. Implement the default counter selection, multi-file processing, total aggregation, and width-aware formatting (still without anticipating stdin logic beyond treating `cfg.Files`). Ensure existing single-flag paths remain untouched.
5. Re-run all tests and then manually spot check `go run ./cmd/wc file1.txt file2.txt` against the system `wc` to validate totals and alignment.
6. Pause for user feedback before tackling stdin.

## The Final Step – Standard Input

Prompt:
1. Unskip the stdin-related fixtures: `internal/accepttest/testdata/stdin.txtar`, `stdin-dash.txtar`, and `stdin-dup-dash.txtar`. Keep other advanced fixtures (e.g. `maxline`, `missing-file`, `invalid-flag`) skipped unless you are ready to address them now.
2. Run the acceptance suite and confirm the failures arise from stdin not being consumed (or from accepting `-` twice) rather than from previously solved features.
3. Update unit tests so `ParseArgs` recognises the lone `-` sentinel, prevents double stdin consumption, and treats an empty file list as stdin. Add coverage to `AnalyzeFiles` to handle an injected `io.Reader` for stdin, ensure `AddTotal` still works when stdin participates, and verify `Format` keeps totals/columns aligned. Run the unit tests to see them fail for the missing stdin handling.
4. Implement stdin routing: wire `cli.App.Run` to pass the provided `stdin` reader when appropriate, update `AnalyzeFiles` to substitute it for `-` or an empty file list, and enforce the “stdin only once” rule surfaced by the tests.
5. Re-run the full test suite (unit + acceptance) and manually check `cat sample.txt | go run ./cmd/wc -l` and mixed cases like `go run ./cmd/wc file.txt -` to ensure behaviour matches the fixtures.
6. Wait for user feedback before moving on to any optional extras (e.g. `-L`, improved error surfaces).

---

Follow this script step by step; at each stage you should finish with a passing test suite, a binary you can exercise manually, and a clear picture of what changed before soliciting feedback for the next increment.
