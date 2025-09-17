package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"wc/internal/format"
	"wc/internal/pipeline/aggregator"
	"wc/internal/pipeline/counters"
)

// Runner executes the wc command with provided dependencies.
type Runner struct {
	Engine aggregator.Engine
}

// Run parses arguments, processes inputs, and writes results to stdout. It returns the process exit code.
func (r Runner) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, files, err := parseArgs(args, stderr)
	if err != nil {
		return 2
	}

	engine := r.Engine
	if engine.TokenizerFactory == nil {
		engine = aggregator.DefaultEngine()
	}

	// Determine columns to display.
	columns := selectColumns(cfg)

	// Determine inputs.
	if len(files) == 0 {
		snap, err := engine.Run(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "wc: stdin: %v\n", err)
			return 1
		}
		rows := []format.Row{{Counts: snap}}
		if err := format.Render(stdout, columns, rows); err != nil {
			fmt.Fprintf(stderr, "wc: %v\n", err)
			return 1
		}
		return 0
	}

	exitCode := 0
	rows := make([]format.Row, 0, len(files)+1)
	total := counters.Snapshot{}
	var success int
	stdinConsumed := false

	for _, name := range files {
		if name == "-" {
			if stdinConsumed {
				fmt.Fprintf(stderr, "wc: stdin already consumed\n")
				exitCode = 1
				continue
			}
			snap, err := engine.Run(stdin)
			if err != nil {
				fmt.Fprintf(stderr, "wc: -: %v\n", err)
				exitCode = 1
				stdinConsumed = true
				continue
			}
			rows = append(rows, format.Row{Label: "-", Counts: snap})
			addSnapshot(&total, snap)
			success++
			stdinConsumed = true
			continue
		}

		snap, err := processFile(engine, name)
		if err != nil {
			fmt.Fprintf(stderr, "wc: %s: %v\n", name, err)
			exitCode = 1
			continue
		}
		rows = append(rows, format.Row{Label: name, Counts: snap})
		addSnapshot(&total, snap)
		success++
	}

	if success > 1 {
		rows = append(rows, format.Row{Label: "total", Counts: total})
	}

	if len(rows) > 0 {
		if err := format.Render(stdout, columns, rows); err != nil {
			fmt.Fprintf(stderr, "wc: %v\n", err)
			return 1
		}
	}

	if exitCode != 0 {
		return exitCode
	}
	if success == 0 {
		return 1
	}
	return 0
}

type config struct {
	lines bool
	words bool
	bytes bool
	runes bool
	max   bool
}

func parseArgs(args []string, errw io.Writer) (config, []string, error) {
	fs := flag.NewFlagSet("wc", flag.ContinueOnError)
	fs.SetOutput(errw)

	cfg := config{}

	fs.BoolVar(&cfg.bytes, "c", false, "print byte counts")
	fs.BoolVar(&cfg.bytes, "bytes", false, "print byte counts")
	fs.BoolVar(&cfg.lines, "l", false, "print line counts")
	fs.BoolVar(&cfg.lines, "lines", false, "print line counts")
	fs.BoolVar(&cfg.words, "w", false, "print word counts")
	fs.BoolVar(&cfg.words, "words", false, "print word counts")
	fs.BoolVar(&cfg.runes, "m", false, "print character counts")
	fs.BoolVar(&cfg.runes, "chars", false, "print character counts")
	fs.BoolVar(&cfg.max, "L", false, "print maximum line length")

	if err := fs.Parse(args); err != nil {
		return cfg, nil, err
	}

	return cfg, fs.Args(), nil
}

func selectColumns(cfg config) []format.Column {
	requested := []format.Column{}
	if cfg.lines {
		requested = append(requested, format.ColumnLines)
	}
	if cfg.words {
		requested = append(requested, format.ColumnWords)
	}
	if cfg.bytes {
		requested = append(requested, format.ColumnBytes)
	}
	if cfg.runes {
		requested = append(requested, format.ColumnRunes)
	}
	if cfg.max {
		requested = append(requested, format.ColumnMaxLine)
	}

	if len(requested) > 0 {
		return requested
	}

	return []format.Column{format.ColumnLines, format.ColumnWords, format.ColumnBytes}
}

func processFile(engine aggregator.Engine, name string) (counters.Snapshot, error) {
	file, err := os.Open(name)
	if err != nil {
		return counters.Snapshot{}, err
	}
	defer file.Close()

	snap, err := engine.Run(file)
	if err != nil {
		return counters.Snapshot{}, err
	}
	return snap, nil
}

func addSnapshot(dst *counters.Snapshot, src counters.Snapshot) {
	dst.Bytes += src.Bytes
	dst.Runes += src.Runes
	dst.Lines += src.Lines
	dst.Words += src.Words
	if src.MaxLineLen > dst.MaxLineLen {
		dst.MaxLineLen = src.MaxLineLen
	}
}
