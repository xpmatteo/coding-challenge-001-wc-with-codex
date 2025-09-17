package main

import (
	"os"

	"wc/internal/cli"
	"wc/internal/pipeline/aggregator"
)

func main() {
	runner := cli.Runner{Engine: aggregator.DefaultEngine()}
	exitCode := runner.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
