package cli

import (
	"fmt"
	"io"

	core "wc/internal/wc"
)

// App executes the wc command with provided dependencies.
type App struct {
}

// Run parses arguments, processes inputs, and writes results to stdout. It returns the process exit code.
func (r App) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, err := core.ParseArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err) //nolint:errcheck // best-effort diagnostic
		return 1
	}

	stats, err := core.AnalyzeFiles(cfg)
	if err != nil {
		fmt.Fprintln(stderr, err) //nolint:errcheck // best-effort diagnostic
		return 1
	}

	stats, err = core.AddTotal(cfg, stats)
	if err != nil {
		fmt.Fprintln(stderr, err) //nolint:errcheck // best-effort diagnostic
		return 1
	}

	lines, err := core.Format(cfg, stats)
	if err != nil {
		fmt.Fprintln(stderr, err) //nolint:errcheck // best-effort diagnostic
		return 1
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(stdout, line); err != nil {
			fmt.Fprintln(stderr, err) //nolint:errcheck // best-effort diagnostic
			return 1
		}
	}

	return 0
}
