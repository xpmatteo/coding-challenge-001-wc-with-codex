package cli

import (
	"io"
)

// App executes the wc command with provided dependencies.
type App struct {
}

// Run parses arguments, processes inputs, and writes results to stdout. It returns the process exit code.
func (r App) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return 0
}
